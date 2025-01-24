package repositories

import "transfigurr/internal/models"

type SeasonRepositoryInterface interface {
	GetSeasons(seriesId string) ([]models.Season, error)
	UpsertSeason(seriesId string, seasonNumber int, inputSeason models.Season) (models.Season, error)
	GetSeasonById(seriesId string, seasonNumber int) (models.Season, error)
	DeleteSeasonById(seriesId string, seasonNumber int) error
}
