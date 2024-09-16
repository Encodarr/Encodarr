package repository

import (
	"strconv"
	"transfigurr/models"

	"github.com/jinzhu/gorm"
)

type SeasonRepository struct {
	DB *gorm.DB
}

func NewSeasonRepository(db *gorm.DB) *SeasonRepository {
	return &SeasonRepository{
		DB: db,
	}
}

func (repo *SeasonRepository) GetSeasons(seriesId string) ([]models.Season, error) {
	var seasons []models.Season
	if err := repo.DB.Where("series_id = ?", seriesId).Preload("Episodes").Find(&seasons).Error; err != nil {
		return nil, err
	}
	return seasons, nil
}

func (repo *SeasonRepository) UpsertSeason(seriesId string, seasonNumber int, inputSeason models.Season) (models.Season, error) {
	var season models.Season
	inputSeason.Id = seriesId + strconv.Itoa(seasonNumber)
	inputSeason.SeriesId = seriesId
	result := repo.DB.Where("series_id = ? AND season_number = ?", seriesId, seasonNumber).First(&season)

	if result.RecordNotFound() {
		season = inputSeason
		if err := repo.DB.Create(&season).Error; err != nil {
			return models.Season{}, err
		}
	} else {
		repo.DB.Model(&season).Updates(inputSeason)
		if err := repo.DB.Save(&season).Error; err != nil {
			return models.Season{}, err
		}
	}
	return season, nil
}

func (repo *SeasonRepository) GetSeasonById(seriesId string, seasonNumber int) (models.Season, error) {
	var season models.Season
	if err := repo.DB.Where("series_id = ? AND season_number = ?", seriesId, seasonNumber).First(&season).Error; err != nil {
		return models.Season{}, err
	}
	return season, nil
}

func (repo *SeasonRepository) DeleteSeasonById(seriesId string, seasonNumber int) error {
	var season models.Season
	if err := repo.DB.Where("series_id = ? AND season_number = ?", seriesId, seasonNumber).First(&season).Error; err != nil {
		return err
	}

	db := repo.DB.Delete(&season)
	if db.Error != nil {
		return db.Error
	}

	return nil
}
