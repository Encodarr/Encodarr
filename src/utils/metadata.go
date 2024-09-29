package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"transfigurr/constants"
	"transfigurr/models"

	"github.com/rainycape/unidecode"

	"github.com/mitchellh/mapstructure"
)

const (
	SERIES_URL  = "https://api.themoviedb.org/3/search/tv"
	MOVIES_URL  = "https://api.themoviedb.org/3/search/movie"
	ARTWORK_URL = "https://image.tmdb.org/t/p/original"
)

var header = func() map[string]string {
	decoded, err := base64.StdEncoding.DecodeString(constants.TEST)
	if err != nil {
		return nil
	}
	return map[string]string{
		"Authorization": "Bearer " + string(decoded),
	}
}()

var re = regexp.MustCompile(`\s*\(.*?\)\s*`)

func parseMovie(movieID string) (models.TMDBMovie, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", MOVIES_URL, nil)
	if err != nil {
		return models.TMDBMovie{}, err
	}

	q := req.URL.Query()
	q.Add("query", unidecode.Unidecode(movieID))
	req.URL.RawQuery = q.Encode()

	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.TMDBMovie{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.TMDBMovie{}, fmt.Errorf("failed to fetch movie data")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.TMDBMovie{}, err
	}

	if len(result["results"].([]interface{})) == 0 {
		return models.TMDBMovie{}, fmt.Errorf("no results found")
	}

	tmdbMovie := result["results"].([]interface{})[0].(map[string]interface{})

	var movie models.TMDBMovie
	if err := mapstructure.Decode(tmdbMovie, &movie); err != nil {
		return models.TMDBMovie{}, err
	}

	return movie, nil
}

func parseSeries(seriesID string) (models.TMDBSeries, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", SERIES_URL, nil)
	if err != nil {
		return models.TMDBSeries{}, err
	}

	q := req.URL.Query()
	q.Add("query", unidecode.Unidecode(seriesID))
	req.URL.RawQuery = q.Encode()

	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.TMDBSeries{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.TMDBSeries{}, fmt.Errorf("failed to fetch series data, status code: %d", resp.StatusCode)
	}

	var seriesSearchResponse struct {
		Results []models.TMDBSeries `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&seriesSearchResponse); err != nil {
		return models.TMDBSeries{}, err
	}

	if len(seriesSearchResponse.Results) == 0 {
		return models.TMDBSeries{}, fmt.Errorf("no results found")
	}

	seriesBestMatch := seriesSearchResponse.Results[0]

	seriesURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%v", seriesBestMatch.ID)
	req, err = http.NewRequest("GET", seriesURL, nil)
	if err != nil {
		return models.TMDBSeries{}, err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	seriesResponse, err := client.Do(req)
	if err != nil {
		return models.TMDBSeries{}, err
	}
	defer seriesResponse.Body.Close()

	if seriesResponse.StatusCode != http.StatusOK {
		return models.TMDBSeries{}, fmt.Errorf("failed to fetch detailed series data, status code: %d", seriesResponse.StatusCode)
	}

	var seriesData models.TMDBSeries
	if err := json.NewDecoder(seriesResponse.Body).Decode(&seriesData); err != nil {
		return models.TMDBSeries{}, err
	}

	return seriesData, nil
}

func downloadMediaArtwork(mediaData map[string]interface{}, mediaID, folder string) error {
	mediaFolder := filepath.Join(folder, mediaID)
	if err := os.MkdirAll(mediaFolder, os.ModePerm); err != nil {
		return err
	}

	client := &http.Client{}

	// Download poster
	if posterPath, ok := mediaData["poster_path"].(string); ok {
		posterURL := ARTWORK_URL + posterPath
		posterFilePath := filepath.Join(mediaFolder, "poster.webp")

		// Remove existing file if it exists
		if _, err := os.Stat(posterFilePath); err == nil {
			if err := os.Remove(posterFilePath); err != nil {
				return err
			}
		}

		resp, err := client.Get(posterURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			return err
		}

		out, err := os.Create(posterFilePath)
		if err != nil {
			return err
		}
		defer out.Close()

		if err := jpeg.Encode(out, img, &jpeg.Options{Quality: constants.PosterQuality}); err != nil {
			return err
		}
	}

	// Download backdrop
	if backdropPath, ok := mediaData["backdrop_path"].(string); ok {
		backdropURL := ARTWORK_URL + backdropPath
		backdropFilePath := filepath.Join(mediaFolder, "backdrop.webp")

		// Remove existing file if it exists
		if _, err := os.Stat(backdropFilePath); err == nil {
			if err := os.Remove(backdropFilePath); err != nil {
				return err
			}
		}

		resp, err := client.Get(backdropURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			return err
		}

		out, err := os.Create(backdropFilePath)
		if err != nil {
			return err
		}
		defer out.Close()

		if err := jpeg.Encode(out, img, &jpeg.Options{Quality: constants.BackdropQuality}); err != nil {
			return err
		}
	}

	return nil
}

func parseEpisode(series models.Series, season models.Season, seriesData models.TMDBSeries, seasonNumber int, episodeNumber int) (*models.Episode, error) {
	client := &http.Client{}
	episodeURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/season/%d/episode/%d", seriesData.ID, seasonNumber, episodeNumber)

	req, err := http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch episode data")
	}

	var episodeData models.TMDBEpisode
	if err := json.NewDecoder(resp.Body).Decode(&episodeData); err != nil {
		return nil, err
	}

	episode := &models.Episode{
		Id:            fmt.Sprintf("%s%d%d", series.Id, seasonNumber, episodeNumber),
		SeriesId:      series.Id,
		SeasonName:    season.Name,
		SeasonNumber:  seasonNumber,
		EpisodeName:   episodeData.Name,
		EpisodeNumber: int(episodeData.EpisodeNumber),
		AirDate:       episodeData.AirDate,
	}

	return episode, nil
}

func GetSeriesMetadata(series models.Series) (models.Series, error) {

	seriesData, err := parseSeries(series.Id)
	if err != nil {
		return series, err
	}
	series.Name = seriesData.Name
	series.Overview = seriesData.Overview
	series.ReleaseDate = seriesData.FirstAirDate
	series.LastAirDate = seriesData.LastAirDate
	if len(seriesData.Genres) > 0 {
		series.Genre = seriesData.Genres[0].Name
	}
	if len(seriesData.Networks) > 0 {
		series.Networks = seriesData.Networks[0].Name
	}
	series.Status = seriesData.Status
	series.Runtime = seriesData.LastEpisodeToAir.Runtime

	// Integrate downloadMediaArtwork
	if err := downloadMediaArtwork(map[string]interface{}{
		"poster_path":   seriesData.PosterPath,
		"backdrop_path": seriesData.BackdropPath,
	}, series.Id, filepath.Join(constants.ConfigPath, "artwork", "series")); err != nil {
		return series, err
	}

	if err != nil {
		return series, err
	}

	for seasonNumber, season := range series.Seasons {
		if season.Episodes == nil {
			continue
		}

		for episodeNumber := range season.Episodes {
			tmdbEpisode, err := parseEpisode(series, season, seriesData, series.Seasons[seasonNumber].SeasonNumber, series.Seasons[seasonNumber].Episodes[episodeNumber].EpisodeNumber)
			if err != nil {
				continue
			}
			series.Seasons[seasonNumber].Episodes[episodeNumber].EpisodeName = tmdbEpisode.EpisodeName
			series.Seasons[seasonNumber].Episodes[episodeNumber].AirDate = tmdbEpisode.AirDate
		}
	}

	return series, nil
}
func GetMovieMetadata(movie models.Movie) (models.Movie, error) {
	client := &http.Client{}

	// Create the search parameters
	searchParams := url.Values{}

	transliteratedID := unidecode.Unidecode(movie.Id)
	cleanedID := re.ReplaceAllString(transliteratedID, "")

	searchParams.Add("query", cleanedID)

	// Create the request
	req, err := http.NewRequest("GET", MOVIES_URL, nil)
	if err != nil {
		return movie, err
	}

	// Set the headers and parameters
	decoded, err := base64.StdEncoding.DecodeString(constants.TEST)
	if err != nil {
		return movie, err
	}
	strtest := string(decoded)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", strtest))
	req.URL.RawQuery = searchParams.Encode()
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return movie, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return movie, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return movie, err
	}
	var movieSearchArray map[string]interface{}
	err = json.Unmarshal(body, &movieSearchArray)
	if err != nil {
		return movie, err
	}

	results, ok := movieSearchArray["results"].([]interface{})
	if !ok || len(results) == 0 {
		return movie, fmt.Errorf("no results found for movie ID: %s", movie.Id)
	}

	movieBestMatch := results[0].(map[string]interface{})
	movieUrl := fmt.Sprintf("https://api.themoviedb.org/3/movie/%v", movieBestMatch["id"])

	// Create the request for the movie
	req, err = http.NewRequest("GET", movieUrl, nil)
	if err != nil {
		return movie, err
	}

	// Create a new header and set it to the request
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", strtest))
	req.Header = header

	// Send the request
	resp, err = client.Do(req)
	if err != nil {
		return movie, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return movie, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response
	var movieData models.TMDBMovie
	err = json.NewDecoder(resp.Body).Decode(&movieData)
	if err != nil {
		return movie, err
	}

	// Set the movie data
	movie.Name = movieData.Title
	movie.Overview = movieData.Overview
	movie.ReleaseDate = movieData.ReleaseDate
	movie.Runtime = movieData.Runtime
	if len(movieData.Genres) > 0 {
		movie.Genre = movieData.Genres[0].Name
	}
	if len(movieData.ProductionCompanies) > 0 {
		movie.Studio = movieData.ProductionCompanies[0].Name
	}
	movie.Status = movieData.Status

	// Integrate downloadMediaArtwork
	if err := downloadMediaArtwork(map[string]interface{}{
		"poster_path":   movieData.PosterPath,
		"backdrop_path": movieData.BackdropPath,
	}, movie.Id, filepath.Join(constants.ConfigPath, "artwork", "movies")); err != nil {
		return movie, err
	}

	return movie, nil
}
