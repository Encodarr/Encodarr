package repository

import (
	"transfigurr/models"

	"github.com/jinzhu/gorm"
)

type ProfileRepository struct {
	DB *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{
		DB: db,
	}
}

func (repo *ProfileRepository) GetAllProfiles() ([]models.Profile, error) {
	var profiles []models.Profile
	if err := repo.DB.Find(&profiles).Error; err != nil {
		return nil, err
	}
	return profiles, nil
}

func (repo *ProfileRepository) GetProfileById(profileId string) (models.Profile, error) {
	var profile models.Profile
	if err := repo.DB.Where("id = ?", profileId).First(&profile).Error; err != nil {
		return models.Profile{}, err
	}
	return profile, nil
}

func (repo *ProfileRepository) UpsertProfile(profileId string, inputProfile models.Profile) (models.Profile, error) {
	var profile models.Profile
	result := repo.DB.Where("id = ?", profileId).First(&profile)
	if result.RecordNotFound() {
		profile = inputProfile
		if err := repo.DB.Create(&profile).Error; err != nil {
			return models.Profile{}, err
		}
	} else {
		repo.DB.Model(&profile).Updates(inputProfile)
		if err := repo.DB.Save(&profile).Error; err != nil {
			return models.Profile{}, err
		}
	}
	return profile, nil
}

func (repo *ProfileRepository) DeleteProfileById(profileId string) error {
	var profile models.Profile
	if err := repo.DB.Where("id = ?", profileId).First(&profile).Error; err != nil {
		return err
	}

	db := repo.DB.Delete(&profile)
	if db.Error != nil {
		return db.Error
	}

	return nil
}
