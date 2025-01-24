package repositories

import "transfigurr/internal/models"

type SettingRepositoryInterface interface {
	GetAllSettings() (map[string]models.Setting, error)
	GetSettingById(id string) (models.Setting, error)
	CreateSetting(setting models.Setting) error
	UpdateSetting(setting models.Setting) error
	DeleteSetting(setting models.Setting) error
}
