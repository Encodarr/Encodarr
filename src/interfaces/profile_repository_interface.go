package interfaces

import "transfigurr/models"

type ProfileRepositoryInterface interface {
	GetAllProfiles() ([]models.Profile, error)
	GetProfileById(id string) (models.Profile, error)
	UpsertProfile(profileId string, inputProfile models.Profile) (models.Profile, error)
	DeleteProfileById(profileId string) error
}
