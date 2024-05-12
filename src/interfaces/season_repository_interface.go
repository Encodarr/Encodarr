package interfaces

import "transfigurr/models"

type SeasonRepositoryInterface interface {
	GetSeasons(seriesId string) ([]models.Season, error)
	UpsertSeason(seriesId string, seasonNumber string, inputSeason models.Season) (models.Season, error)
	GetSeasonById(seriesId string, seasonNumber string) (models.Season, error)
	DeleteSeasonById(seriesId string, seasonNumber string) error
}
