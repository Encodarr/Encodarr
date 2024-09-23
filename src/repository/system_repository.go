package repository

import (
	"errors"
	"transfigurr/models"

	"gorm.io/gorm"
)

type SystemRepository struct {
	DB *gorm.DB
}

func NewSystemRepository(db *gorm.DB) *SystemRepository {
	return &SystemRepository{
		DB: db,
	}
}

func (repo *SystemRepository) GetSystems() ([]models.System, error) {
	var systems []models.System
	if err := repo.DB.Find(&systems).Error; err != nil {
		return nil, err
	}
	return systems, nil
}

func (repo *SystemRepository) UpsertSystem(systemId string, inputSystem models.System) (models.System, error) {
	var system models.System
	inputSystem.Id = systemId
	result := repo.DB.Where("id = ?", systemId).First(&system)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		system = inputSystem
		if err := repo.DB.Create(&system).Error; err != nil {
			return models.System{}, err
		}
	} else {
		repo.DB.Model(&system).Updates(inputSystem)
		if err := repo.DB.Save(&system).Error; err != nil {
			return models.System{}, err
		}
	}
	return system, nil
}

func (repo *SystemRepository) GetSystemById(systemId string) (models.System, error) {
	var system models.System
	if err := repo.DB.Where("id = ?", systemId).First(&system).Error; err != nil {
		return models.System{}, err
	}
	return system, nil
}

func (repo *SystemRepository) DeleteSystemById(systemId string) error {
	var system models.System
	if err := repo.DB.Where("id = ?", systemId).First(&system).Error; err != nil {
		return err
	}

	db := repo.DB.Delete(&system)
	if db.Error != nil {
		return db.Error
	}

	return nil
}
