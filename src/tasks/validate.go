package tasks

import (
	"log"
	"os"
	"path/filepath"
	"transfigurr/constants"
	"transfigurr/interfaces"
)

func ValidateSeries(seriesID string, seriesRepo interfaces.SeriesRepositoryInterface, seasonRepo interfaces.SeasonRepositoryInterface, episodeRepo interfaces.EpisodeRepositoryInterface) error {
	log.Println("Validating series:", seriesID)

	series, err := seriesRepo.GetSeriesByID(seriesID)

	seriesPath := filepath.Join(constants.SeriesPath, seriesID)
	if _, err := os.Stat(seriesPath); os.IsNotExist(err) {
		return seriesRepo.DeleteSeriesByID(seriesID)
	}

	for _, season := range series.Seasons {
		for _, episode := range season.Episodes {
			episodePathInSeries := filepath.Join(seriesPath, episode.Filename)
			episodePathInSeason := filepath.Join(seriesPath, episode.SeasonName, episode.Filename)
			if _, err := os.Stat(episodePathInSeries); os.IsNotExist(err) {
				if _, err := os.Stat(episodePathInSeason); os.IsNotExist(err) {
					if err := episodeRepo.DeleteEpisodeById(series.Id, season.SeasonNumber, episode.EpisodeNumber); err != nil {
						log.Println("Error removing episode:", episode.Id, err)
					}
				}
			}
		}

		if err != nil {
			return err
		}

		for _, season := range series.Seasons {
			if len(season.Episodes) == 0 {
				if err := seasonRepo.DeleteSeasonById(series.Id, season.SeasonNumber); err != nil {
					log.Println("Error removing season:", season.Id, err)
				}
			}
		}
	}
	return nil
}

func ValidateMovie(movieID string, movieRepo interfaces.MovieRepositoryInterface) error {
	log.Println("Validating movie:", movieID)

	movie, err := movieRepo.GetMovieById(movieID)
	if err != nil {
		log.Printf("Error fetching movie with ID %s: %v", movieID, err)
		return err
	}

	log.Printf("Fetched movie: %+v", movie)

	moviePath := filepath.Join(constants.MoviesPath, movieID)
	if _, err := os.Stat(moviePath); os.IsNotExist(err) {
		log.Printf("Movie directory not found: %s", moviePath)
		err := movieRepo.DeleteMovieById(movieID)
		if err != nil {
			log.Printf("Error deleting movie with ID %s: %v", movieID, err)
		}
		log.Print("Movie deleted successfully")
		return err
	}

	log.Printf("Movie validation completed successfully for ID: %s", movieID)
	return nil
}
