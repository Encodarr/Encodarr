package interfaces

import "transfigurr/models"

type SystemRepositoryInterface interface {
	GetSystems() ([]models.System, error)
	UpsertSystem(systemId string, inputSystem models.System) (models.System, error)
	GetSystemById(systemID string) (models.System, error)
	DeleteSystemById(systemID string) error
}
