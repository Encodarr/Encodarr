package utils

import (
	"os"
	"path/filepath"
	"transfigurr/internal/config"
	"transfigurr/internal/interfaces/repositories"
)

func ValidateSeries(seriesID string, seriesRepo repositories.SeriesRepositoryInterface, seasonRepo repositories.SeasonRepositoryInterface, episodeRepo repositories.EpisodeRepositoryInterface) error {

	series, err := seriesRepo.GetSeriesByID(seriesID)

	seriesPath := filepath.Join(config.SeriesPath, seriesID)
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
						//log.Println("Error removing episode:", episode.Id, err)
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
					//log.Println("Error removing season:", season.Id, err)
				}
			}
		}
	}
	return nil
}

func ValidateMovie(movieID string, movieRepo repositories.MovieRepositoryInterface) error {

	_, err := movieRepo.GetMovieById(movieID)
	if err != nil {
		return err
	}

	moviePath := filepath.Join(config.MoviesPath, movieID)
	if _, err := os.Stat(moviePath); os.IsNotExist(err) {
		err := movieRepo.DeleteMovieById(movieID)
		if err != nil {
		}
		return err
	}

	return nil
}
