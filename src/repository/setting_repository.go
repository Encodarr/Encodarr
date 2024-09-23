package repository

import (
	"log"
	"transfigurr/models"

	"gorm.io/gorm"
)

type SettingRepository struct {
	DB *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{
		DB: db,
	}
}

func (repo *SettingRepository) GetAllSettings() (map[string]models.Setting, error) {
	var settings []models.Setting
	err := repo.DB.Find(&settings).Error
	if err != nil {
		return nil, err
	}

	settingMap := make(map[string]models.Setting)
	for _, setting := range settings {
		settingMap[setting.Id] = setting
	}

	return settingMap, nil
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
	log.Print(setting.Id, setting.Value)
	return repo.DB.Save(&setting).Select("*").Error
}

func (repo *SettingRepository) DeleteSetting(setting models.Setting) error {
	return repo.DB.Delete(&setting).Error
}
