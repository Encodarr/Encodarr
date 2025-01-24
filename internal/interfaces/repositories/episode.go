package repositories

import (
	"transfigurr/internal/models"
)

type EpisodeRepositoryInterface interface {
	GetEpisodes(seriesId string, seasonNumber int) ([]models.Episode, error)
	UpsertEpisode(seriesId string, seasonNumber int, episodeNumber int, inputEpisode models.Episode) (models.Episode, error)
	GetEpisodeBySeriesSeasonEpisode(seriesId string, seasonNumber int, episodeNumber int) (models.Episode, error)
	GetEpisodeById(episodeId string) (models.Episode, error)
	DeleteEpisodeById(seriesId string, seasonNumber int, episodeNumber int) error
}
