package repository

import (
	"transfigurr/models"

	"github.com/jinzhu/gorm"
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
	result := repo.DB.Where("id = ?", id).First(&series)

	if result.RecordNotFound() {
		series = inputSeries
		if err := repo.DB.Create(&series).Error; err != nil {
			return models.Series{}, err
		}
	} else {
		repo.DB.Model(&series).Updates(inputSeries)
		if err := repo.DB.Save(&series).Error; err != nil {
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
