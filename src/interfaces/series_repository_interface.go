package interfaces

import "transfigurr/models"

type SeriesRepositoryInterface interface {
	GetSeries() ([]models.Series, error)
	UpsertSeries(id string, series models.Series) (models.Series, error)
	GetSeriesByID(id string) (models.Series, error)
	DeleteSeriesByID(id string) error
}
