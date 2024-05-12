package repository

import (
	"transfigurr/models"

	"github.com/jinzhu/gorm"
)

type SettingRepository struct {
	DB *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{
		DB: db,
	}
}

func (repo *SettingRepository) GetAllSettings() ([]models.Setting, error) {
	var settings []models.Setting
	err := repo.DB.Find(&settings).Error
	return settings, err
}

func (repo *SettingRepository) GetSettingById(id string) (models.Setting, error) {
	var setting models.Setting
	err := repo.DB.Where("id = ?", id).First(&setting).Error
	return setting, err
}

func (repo *SettingRepository) CreateSetting(setting models.Setting) error {
	return repo.DB.Create(&setting).Error
}

func (repo *SettingRepository) UpdateSetting(setting models.Setting) error {
	return repo.DB.Save(&setting).Error
}

func (repo *SettingRepository) DeleteSetting(setting models.Setting) error {
	return repo.DB.Delete(&setting).Error
}
