package repositories

import "transfigurr/internal/models"

type ProfileRepositoryInterface interface {
	GetAllProfiles() ([]models.Profile, error)
	GetProfileById(id int) (models.Profile, error)
	UpsertProfile(profileId int, inputProfile models.Profile) (models.Profile, error)
	DeleteProfileById(profileId int) error
}
