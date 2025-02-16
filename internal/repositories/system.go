package repositories

import (
	"database/sql"
	"transfigurr/internal/models"
)

type SystemRepository struct {
	DB *sql.DB
}

func NewSystemRepository(db *sql.DB) *SystemRepository {
	return &SystemRepository{
		DB: db,
	}
}

func (repo *SystemRepository) GetSystems() ([]models.System, error) {
	rows, err := repo.DB.Query("SELECT id, value FROM systems")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var systems []models.System
	for rows.Next() {
		var system models.System
		err := rows.Scan(&system.Id, &system.Value)
		if err != nil {
			return nil, err
		}
		systems = append(systems, system)
	}
	return systems, nil
}

func (repo *SystemRepository) UpsertSystem(systemId string, inputSystem models.System) (models.System, error) {
	var exists bool
	err := repo.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM systems WHERE id = ?)", systemId).Scan(&exists)
	if err != nil {
		return models.System{}, err
	}

	if exists {
		_, err = repo.DB.Exec(
			"UPDATE systems SET value = ? WHERE id = ?",
			inputSystem.Value, systemId,
		)
	} else {
		_, err = repo.DB.Exec(
			"INSERT INTO systems (id, value) VALUES (?, ?)",
			systemId, inputSystem.Value,
		)
	}
	if err != nil {
		return models.System{}, err
	}

	return repo.GetSystemById(systemId)
}

func (repo *SystemRepository) GetSystemById(systemId string) (models.System, error) {
	var system models.System
	err := repo.DB.QueryRow("SELECT id, value FROM systems WHERE id = ?", systemId).
		Scan(&system.Id, &system.Value)
	if err != nil {
		return models.System{}, err
	}
	return system, nil
}

func (repo *SystemRepository) DeleteSystemById(systemId string) error {
	_, err := repo.DB.Exec("DELETE FROM systems WHERE id = ?", systemId)
	return err
}
