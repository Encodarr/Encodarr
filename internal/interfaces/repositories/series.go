package repositories

import "transfigurr/internal/models"

type SeriesRepositoryInterface interface {
	GetSeries() ([]models.Series, error)
	UpsertSeries(id string, series models.Series) (models.Series, error)
	GetSeriesByID(id string) (models.Series, error)
	DeleteSeriesByID(id string) error
}
