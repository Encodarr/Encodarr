import asyncio
from collections import defaultdict
from functools import partial
from pathlib import Path
import os

file_locks = defaultdict(asyncio.Lock)


async def verify_folders():
    root_folder_path = await get_root_folder()
    paths = []
    config_folder_path = root_folder_path + '/config'
    paths.append(config_folder_path)
    paths.append(config_folder_path + '/artwork/series')
    paths.append(config_folder_path + '/artwork/movies')
    paths.append(config_folder_path + '/db')
    paths.append(root_folder_path + '/series')
    paths.append(root_folder_path + '/movies')
    paths.append(root_folder_path + '/transcode')

    loop = asyncio.get_event_loop()

    for path in paths:
        await loop.run_in_executor(None, partial(os.makedirs, path, exist_ok=True))
    return


async def get_root_folder():
    current_file = Path(__file__).resolve()
    root_folder = current_file.parents[2]
    return str(root_folder)


def get_root_folder2():
    current_file = Path(__file__).resolve()
    root_folder = current_file.parents[2]
    return str(root_folder)


async def get_config_folder():
    return await get_root_folder() + '/config'


async def get_series_artwork_folder():
    return await get_root_folder() + '/config/artwork/series'


async def get_series_folder():
    return await get_root_folder() + '/series'


async def get_movies_folder():
    return await get_root_folder() + '/movies'


async def get_series_metadata_folder():
    return await get_root_folder() + '/config/metadata/series'


async def get_transcode_folder():
    return await get_root_folder() + '/transcode'
