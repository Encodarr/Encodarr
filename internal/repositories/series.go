package repositories

import (
	"errors"
	"transfigurr/internal/models"

	"gorm.io/gorm"
)

type SeriesRepository struct {
	DB *gorm.DB
}

func NewSeriesRepository(db *gorm.DB) *SeriesRepository {
	return &SeriesRepository{
		DB: db,
	}
}

func (repo *SeriesRepository) GetSeries() ([]models.Series, error) {
	var seriesList []models.Series
	if err := repo.DB.Preload("Seasons.Episodes").Find(&seriesList).Error; err != nil {
		return nil, err
	}
	return seriesList, nil
}

func (repo *SeriesRepository) UpsertSeries(id string, inputSeries models.Series) (models.Series, error) {
	var series models.Series
	inputSeries.Id = id

	// Try to find the series by ID
	result := repo.DB.Where("id = ?", id).First(&series)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Series not found, create a new one
		series = inputSeries
		if err := repo.DB.Create(&series).Error; err != nil {
			return models.Series{}, err
		}
	} else if result.Error != nil {
		// Other errors
		return models.Series{}, result.Error
	} else {
		// Series found, update it
		if err := repo.DB.Model(&series).Select("*").Updates(inputSeries).Error; err != nil {
			return models.Series{}, err
		}
	}

	return series, nil
}

func (repo *SeriesRepository) GetSeriesByID(id string) (models.Series, error) {
	var series models.Series
	if err := repo.DB.Preload("Seasons.Episodes").Where("id = ?", id).First(&series).Error; err != nil {
		return models.Series{}, err
	}
	return series, nil
}

func (repo *SeriesRepository) DeleteSeriesByID(id string) error {
	var series models.Series
	if err := repo.DB.Where("id = ?", id).First(&series).Error; err != nil {
		return err
	}

	db := repo.DB.Delete(&series)
	if db.Error != nil {
		return db.Error
	}

	return nil
}
