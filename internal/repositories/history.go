package repositories

import (
	"errors"
	"transfigurr/internal/models"

	"gorm.io/gorm"
)

type HistoryRepository struct {
	DB *gorm.DB
}

func NewHistoryRepository(db *gorm.DB) *HistoryRepository {
	return &HistoryRepository{
		DB: db,
	}
}

func (repo *HistoryRepository) GetHistories() ([]models.History, error) {
	var histories []models.History
	if err := repo.DB.Find(&histories).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return histories, nil
}

func (repo *HistoryRepository) GetHistoryById(id string) (models.History, error) {
	var history models.History
	if err := repo.DB.Where("id = ?", id).First(&history).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.History{}, ErrRecordNotFound
		}
		return history, err
	}
	return history, nil
}

func (repo *HistoryRepository) UpsertHistoryById(history *models.History) error {
	if err := repo.DB.Save(history).Error; err != nil {
		return err
	}
	return nil
}

func (repo *HistoryRepository) DeleteHistoryById(history *models.History) error {
	if err := repo.DB.Delete(history).Error; err != nil {
		return err
	}
	return nil
}
