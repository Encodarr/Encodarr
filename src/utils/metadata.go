package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"transfigurr/constants"
	"transfigurr/models"

	"github.com/chai2010/webp"
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

func parseSeries(series models.Series) (models.Series, error) {
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

func parseEpisode(series models.Series, season models.Season, seasonNumber int, episodeNumber int) (models.Episode, error) {
	seriesID := series.Id
	episodeURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%v/season/%d/episode/%d", seriesID, seasonNumber, episodeNumber)

	client := &http.Client{}
	req, err := http.NewRequest("GET", episodeURL, nil)
	if err != nil {
		log.Printf("An error occurred while creating the request: %v", err)
		return models.Episode{}, err
	}

	// Add headers to the request
	decoded, err := base64.StdEncoding.DecodeString(constants.TEST)
	if err != nil {
		log.Printf("Error decoding base64: %v\n", err)
		return models.Episode{}, err
	}
	strtest := string(decoded)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", strtest))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("An error occurred while making the request: %v", err)
		return models.Episode{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-200 response: %d", resp.StatusCode)
		return models.Episode{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var episodeData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&episodeData); err != nil {
		log.Printf("An error occurred while decoding the response: %v", err)
		return models.Episode{}, err
	}

	episode := models.Episode{
		Id:            seriesID + fmt.Sprintf("%d%d", seasonNumber, episodeNumber),
		SeriesId:      seriesID,
		SeasonName:    season.Name,
		SeasonNumber:  seasonNumber,
		EpisodeName:   episodeData["name"].(string),
		EpisodeNumber: int(episodeData["episode_number"].(float64)),
		AirDate:       episodeData["air_date"].(string),
	}

	return episode, nil
}

func GetSeriesMetadata(series models.Series) (models.Series, []models.Episode, error) {
	_, err := parseSeries(series)
	if err != nil {
		log.Printf("An error occurred while parsing series: %v", err)
		return series, nil, err
	}

	//configFolder := constants.ConfigPath
	if err != nil {
		log.Printf("An error occurred while getting config folder: %v", err)
		return series, nil, err
	}

	//err = downloadMediaArtwork(seriesData, series.Id, configFolder+"/artwork/series")
	if err != nil {
		log.Printf("An error occurred while downloading media artwork: %v", err)
		return series, nil, err
	}

	var episodes []models.Episode
	for seasonNumber, season := range series.Seasons {
		if season.Episodes == nil {
			continue
		}
		for episodeNumber := range season.Episodes {
			episode, err := parseEpisode(series, season, seasonNumber, episodeNumber)
			if err != nil {
				log.Printf("An error occurred while parsing episode: %v", err)
				continue
			}
			episodes = append(episodes, episode)
		}
	}
	return series, episodes, nil
}

func downloadMediaArtwork(mediaData models.TMDBSeries, mediaID string, folder string) error {
	mediaFolder := filepath.Join(folder, mediaID)
	err := os.MkdirAll(mediaFolder, os.ModePerm)
	if err != nil {
		log.Printf("An error occurred while creating the folder: %v", err)
		return err
	}

	client := &http.Client{}

	// Download poster
	posterPath := mediaData.PosterPath
	if posterPath == "" {
		log.Println("No poster path provided.")
	} else {
		posterURL := fmt.Sprintf("https://image.tmdb.org/t/p/original%s", posterPath)
		posterFilePath := filepath.Join(mediaFolder, "poster.webp")
		if _, err := os.Stat(posterFilePath); os.IsNotExist(err) {
			response, err := client.Get(posterURL)
			if err != nil || response.StatusCode != http.StatusOK {
				log.Println("Failed to download poster.")
			} else {
				defer response.Body.Close()
				img, _, err := image.Decode(response.Body)
				if err != nil {
					log.Printf("An error occurred while decoding the poster: %v", err)
				} else {
					file, err := os.Create(posterFilePath)
					if err != nil {
						log.Printf("An error occurred while creating the poster file: %v", err)
					} else {
						defer file.Close()
						err = webp.Encode(file, img, &webp.Options{Quality: 5})
						if err != nil {
							log.Printf("An error occurred while saving the poster: %v", err)
						}
					}
				}
			}
		}
	}

	// Download backdrop
	backdropPath := mediaData.BackdropPath
	if backdropPath == "" {
		log.Println("No backdrop path provided.")
	} else {
		backdropURL := fmt.Sprintf("https://image.tmdb.org/t/p/original%s", backdropPath)
		backdropFilePath := filepath.Join(mediaFolder, "backdrop.webp")
		if _, err := os.Stat(backdropFilePath); os.IsNotExist(err) {
			response, err := client.Get(backdropURL)
			if err != nil || response.StatusCode != http.StatusOK {
				log.Println("Failed to download backdrop.")
			} else {
				defer response.Body.Close()
				img, _, err := image.Decode(response.Body)
				if err != nil {
					log.Printf("An error occurred while decoding the backdrop: %v", err)
				} else {
					file, err := os.Create(backdropFilePath)
					if err != nil {
						log.Printf("An error occurred while creating the backdrop file: %v", err)
					} else {
						defer file.Close()
						err = webp.Encode(file, img, &webp.Options{Quality: 5})
						if err != nil {
							log.Printf("An error occurred while saving the backdrop: %v", err)
						}
					}
				}
			}
		}
	}
	return nil
}
