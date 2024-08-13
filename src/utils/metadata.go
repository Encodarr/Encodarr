package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"transfigurr/constants"
	"transfigurr/models"
)

func GetMovieMetadata(movie models.Movie) (models.Movie, error) {
	client := &http.Client{}

	// Create the search parameters
	searchParams := url.Values{}
	searchParams.Add("query", movie.Id)

	// Create the request
	req, err := http.NewRequest("GET", constants.MOVIE_URL, nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return movie, err
	}

	// Set the headers and parameters
	decoded, err := base64.StdEncoding.DecodeString(constants.TEST)
	if err != nil {
		log.Printf("Error decoding base64: %v\n", err)
		return movie, err
	}
	strtest := string(decoded)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", strtest))
	req.URL.RawQuery = searchParams.Encode()

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return movie, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return movie, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return movie, err
	}
	var movieSearchArray map[string]interface{}
	err = json.Unmarshal(body, &movieSearchArray)
	if err != nil {
		log.Printf("Error unmarshalling response body: %v\n", err)
		return movie, err
	}

	if movieSearchArray["results"] == nil {
		return movie, nil
	}

	movieBestMatch := movieSearchArray["results"].([]interface{})[0].(map[string]interface{})
	movieUrl := fmt.Sprintf("https://api.themoviedb.org/3/movie/%v", movieBestMatch["id"])

	// Create the request for the movie
	req, err = http.NewRequest("GET", movieUrl, nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return movie, err
	}

	// Create a new header and set it to the request
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", strtest))
	req.Header = header

	// Send the request
	resp, err = client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
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
		log.Printf("Error unmarshalling response body: %v\n", err)
		return movie, err
	}

	// Set the movie data
	movie.Name = movieData.Title
	movie.Overview = movieData.Overview
	movie.ReleaseDate = movieData.ReleaseDate
	movie.Runtime = movieData.Runtime
	movie.Genre = movieData.Genres[0].Name
	movie.Studio = movieData.ProductionCompanies[0].Name
	movie.Status = movieData.Status
	return movie, nil
}

func GetSeriesMetadata(series models.Series) (models.Series, error) {
	client := &http.Client{}

	// Create the search parameters
	searchParams := url.Values{}
	searchParams.Add("query", series.Id)

	// Create the request
	req, err := http.NewRequest("GET", constants.SERIES_URL, nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return series, err
	}

	// Set the headers and parameters
	decoded, err := base64.StdEncoding.DecodeString(constants.TEST)
	if err != nil {
		log.Printf("Error decoding base64: %v\n", err)
		return series, err
	}
	strtest := string(decoded)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", strtest))
	req.URL.RawQuery = searchParams.Encode()

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return series, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return series, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return series, err
	}
	var seriesSearchArray map[string]interface{}
	err = json.Unmarshal(body, &seriesSearchArray)
	if err != nil {
		log.Printf("Error unmarshalling response body: %v\n", err)
		return series, err
	}

	if seriesSearchArray["results"] == nil {
		return series, nil
	}

	seriesBestMatch := seriesSearchArray["results"].([]interface{})[0].(map[string]interface{})
	seriesUrl := fmt.Sprintf("https://api.themoviedb.org/3/tv/%v", seriesBestMatch["id"])

	// Create the request for the series
	req, err = http.NewRequest("GET", seriesUrl, nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return series, err
	}

	// Create a new header and set it to the request
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", strtest))
	req.Header = header

	// Send the request
	resp, err = client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return series, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return series, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	// Parse the response
	var seriesData models.TMDBSeries
	err = json.NewDecoder(resp.Body).Decode(&seriesData)
	if err != nil {
		log.Printf("Error unmarshalling response body: %v\n", err)
		return series, err
	}

	// Set the series data
	series.Name = seriesData.Name
	series.Overview = seriesData.Overview
	series.Runtime = seriesData.EpisodeRunTime[0]
	series.Genre = seriesData.Genres[0].Name
	series.Status = seriesData.Status
	series.ReleaseDate = seriesData.FirstAirDate
	series.LastAirDate = seriesData.LastAirDate
	series.Networks = seriesData.Networks[0].Name
	return series, nil
}
