package repositories

import (
	"database/sql"
	"transfigurr/internal/models"
)

type SettingRepository struct {
	DB *sql.DB
}

func NewSettingRepository(db *sql.DB) *SettingRepository {
	return &SettingRepository{
		DB: db,
	}
}

func (repo *SettingRepository) GetAllSettings() (map[string]models.Setting, error) {
	rows, err := repo.DB.Query("SELECT id, value FROM settings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settingMap := make(map[string]models.Setting)
	for rows.Next() {
		var setting models.Setting
		err := rows.Scan(&setting.Id, &setting.Value)
		if err != nil {
			return nil, err
		}
		settingMap[setting.Id] = setting
	}

	return settingMap, nil
}

func (repo *SettingRepository) GetSettingById(id string) (models.Setting, error) {
	var setting models.Setting
	err := repo.DB.QueryRow("SELECT id, value FROM settings WHERE id = ?", id).
		Scan(&setting.Id, &setting.Value)
	return setting, err
}

func (repo *SettingRepository) CreateSetting(setting models.Setting) error {
	_, err := repo.DB.Exec(
		"INSERT INTO settings (id, value) VALUES (?, ?)",
		setting.Id, setting.Value,
	)
	return err
}

func (repo *SettingRepository) UpdateSetting(setting models.Setting) error {
	_, err := repo.DB.Exec(
		"UPDATE settings SET value = ? WHERE id = ?",
		setting.Value, setting.Id,
	)
	return err
}

func (repo *SettingRepository) DeleteSetting(setting models.Setting) error {
	_, err := repo.DB.Exec("DELETE FROM settings WHERE id = ?", setting.Id)
	return err
}
