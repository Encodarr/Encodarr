package interfaces

import "transfigurr/models"

type HistoryRepositoryInterface interface {
	GetHistories() ([]models.History, error)
	GetHistoryById(id string) (models.History, error)
	UpsertHistoryById(history models.History) error
	DeleteHistoryById(history models.History) error
}
