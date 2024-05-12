package interfaces

import (
	"transfigurr/models"
)

type EpisodeRepositoryInterface interface {
	GetEpisodes(seriesId string, seasonNumber string) ([]models.Episode, error)
	UpsertEpisode(seriesId string, seasonNumber string, episodeNumber string, inputEpisode models.Episode) (models.Episode, error)
	GetEpisodeById(seriesId string, seasonNumber string, episodeNumber string) (models.Episode, error)
	DeleteEpisodeById(seriesId string, seasonNumber string, episodeNumber string) error
}
