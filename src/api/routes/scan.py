from dataclasses import asdict
import os
import re
from dotenv import dotenv_values
from fastapi import APIRouter
from src.api.routes.profiles import getProfiles

import httpx
import aiofiles

from src.api.utils import analyze_media_file, get_config_folder, get_series_folder, get_series_metadata_folder, open_json, write_json
from src.models.episode_model import episode_model
from src.models.profile_model import profile_model
from src.models.season_model import season_model
from src.models.series_model import series_model
from src.models.settings_model import settings_model
from src.global_state import GlobalState

global_state = GlobalState()

router = APIRouter()
config = dotenv_values("src/.env")
API_KEY = config['API_KEY']



@router.get("/api/scan/series/metadata")
async def get_all_series_metadata():
    series_list = await global_state.get_series()
    for series_name in series_list:
        await get_series_metadata(series_name)
    return

@router.get("/api/scan/series/metadata/{series_name}") 
async def get_series_metadata(series_name):
    print('getting metadata for', series_name)
    async with httpx.AsyncClient() as client:
        series: series_model = series_model()
        series_folder = await get_series_metadata_folder()
        metadata_folder_path = os.path.join(series_folder, series_name)
        tmdb_url = f"https://api.themoviedb.org/3/search/tv"
        params = {"api_key": API_KEY, "query": series_name}
        response = await client.get(tmdb_url, params=params)
        if response.status_code == 200:
            data = response.json()
            if data.get("results"):
                series_data = data["results"][0]
                series.id = series_data.get("id")
                if series.id:
                    season_url = f"https://api.themoviedb.org/3/tv/{series.id}"
                    params = {"api_key": API_KEY}
                    response = await client.get(season_url, params=params)
                    if response.status_code == 200:
                        data = response.json()
                        series.name = data.get("name")
                        series.overview = data.get("overview")
                        series.first_air_date = data.get("first_air_date")
                        series.last_air_date = data.get("last_air_date")
                        series.genre = data.get("genres")[0]['name']
                        series.networks = data.get("networks")[0]['name']
                        series.status = data.get('status')
                        series.seasons = {}
                        series.profile = 0
                # Download poster
                poster_path = series_data.get("poster_path")
                if poster_path:
                    poster_url = f"https://image.tmdb.org/t/p/original{poster_path}"
                    poster_path = os.path.join(metadata_folder_path, "poster.jpg")
                    response = await client.get(poster_url)
                    if response.status_code == 200:
                        async with aiofiles.open(poster_path, "wb") as poster_file:
                            await poster_file.write(response.content)                # Download backdrop
                backdrop_path = series_data.get("backdrop_path")
                if backdrop_path:
                    backdrop_url = f"https://image.tmdb.org/t/p/original{backdrop_path}"
                    backdrop_file_path = os.path.join(metadata_folder_path, "backdrop.jpg")
                    response = await client.get(backdrop_url)
                    if response.status_code == 200:
                        async with aiofiles.open(backdrop_file_path, "wb") as backdrop_file:
                            await backdrop_file.write(response.content)
                if series.id:
                    season_url = f"https://api.themoviedb.org/3/tv/{series.id}"
                    params = {"api_key": API_KEY}
                    season_response = await client.get(season_url, params=params)
                    if season_response.status_code == 200:
                        season_data = season_response.json()
                        for season in season_data.get('seasons'):
                            newSeason: season_model = season_model()
                            newSeason.name = season['name']
                            newSeason.season_number = season['season_number']
                            newSeason.episode_count = season['episode_count']
                            newSeason.overview = season['overview']
                            newSeason.episodes = {}
                            series.seasons[season['season_number']] = newSeason
                            newSeason.profile = 0
                            for episode_number in range(1, newSeason.episode_count + 1):
                                episode_url = f"https://api.themoviedb.org/3/tv/{series.id}/season/{newSeason.season_number}/episode/{episode_number}"
                                episode_response = await client.get(episode_url, params=params)
                                if episode_response.status_code == 200:
                                    episode= episode_model()
                                    episode_data = episode_response.json()
                                    episode.series_name = series.name
                                    episode.season_name = newSeason.name
                                    episode.season_number = newSeason.season_number
                                    episode.episode_name = episode_data.get('name')
                                    episode.episode_number = episode_data.get('episode_number')
                                    episode.profile = 0
                                    newSeason.episodes[episode_number] = episode
                                else:
                                    print(f"Failed to fetch episode {episode_number} information. Status code: {response.status_code}")
        await global_state.set_tvdb(series_name, asdict(series))
        return


@router.get('/api/scan/series')
async def scan_all_series():
    print('scanning all series')
    series_folder = await get_series_folder()
    series_result = await global_state.get_series_list()
    series = set()
    for i in series_result:
        series.add(i)
    for series_name in os.listdir(series_folder):
        if series_name == ".DS_Store":
            continue
        await scan_series(series_name)
        series.add(series_name)
        await global_state.set_series_list(list(series))
    return

        


@router.get("/api/scan/series/{series_name}")
async def scan_series(series_name):
    if series_name == ".DS_Store":
        return
    series_path = os.path.join(await get_series_folder(), series_name)
    if os.path.isdir(series_path) and series_name != '.DS_Store':
        series = {"name": series_name, "series_path": series_path, "seasons": {}}
        for season_name in os.listdir(series_path):
            if season_name == ".DS_Store":
                continue
            season_number = int("".join(re.findall(r'\d+', season_name)))
            season_path = os.path.join(series_path, season_name)
            if os.path.isdir(season_path) and season_name != '.DS_Store':
                files = [f for f in os.listdir(season_path) if os.path.isfile(os.path.join(season_path, f))]
                episodes = {}
                for file in files:
                    pattern = r"(?:S(\d{2})E(\d{2})|E(\d{2}))"
                    match = re.search(pattern, file)
                    if match:
                        if match.group(1):
                            episode_number = int(match.group(2))
                        else:
                            episode_number = int(match.group(3))
                        episode_path = os.path.join(season_path, file)
                        analysis_data = await analyze_media_file(episode_path)
                        episodes[episode_number] = {'episode_name': file, 'season': season_number, 'episode_number': episode_number,'filename': file, 'path': season_path + '/', 'codec': analysis_data}
                    else:
                        print('missing')
                series["seasons"][season_number] = {"season": season_number, "name": season_name, "episodes": episodes}
        await global_state.set_series(series_name, series)
        await global_state.get_series_config(series_name)
        if await global_state.get_tvdb(series_name) == {}:
            print('about to fetch metadata for', series_name)
            await get_series_metadata(series_name)
        return
    
@router.get('/api/scan/queue')
async def scan_queue():
    print('scanning queue')
    await scan_all_series()
    q = []
    series_names = await global_state.get_series_list()
    profiles = await getProfiles()
    for series_name in series_names:
        if series_name == ".DS_Store":
            continue
        settings = await global_state.get_series_config(series_name)
        profile_id = settings['profile_id']
        monitored = settings['monitored']
        if not monitored:
            continue
        profile = profiles[profile_id]
        codec = profile['codec']
        path = os.path.join(await get_series_metadata_folder(),series_name)
        series = await global_state.get_series(series_name)
        for season in series['seasons']:
            for episode in series['seasons'][season]['episodes']:
                if series['seasons'][season]['episodes'][episode]['codec'] != codec:
                    episode = series['seasons'][season]['episodes'][episode]
                    q.append({'profile': profile_id, 'episode': episode})
    await global_state.set_queue(q)
