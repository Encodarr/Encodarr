package repository

import (
	"strconv"
	"transfigurr/models"

	"github.com/jinzhu/gorm"
)

type EpisodeRepository struct {
	DB *gorm.DB
}

func NewEpisodeRepository(db *gorm.DB) *EpisodeRepository {
	return &EpisodeRepository{
		DB: db,
	}
}

func (repo *EpisodeRepository) GetEpisodes(seriesId string, seasonNumber int) ([]models.Episode, error) {
	var episodes []models.Episode
	if err := repo.DB.Where("series_id = ? AND season_number = ?", seriesId, seasonNumber).Find(&episodes).Error; err != nil {
		return nil, err
	}
	return episodes, nil
}

func (repo *EpisodeRepository) UpsertEpisode(seriesId string, seasonNumber int, episodeNumber int, inputEpisode models.Episode) (models.Episode, error) {
	var episode models.Episode
	inputEpisode.Id = seriesId + strconv.Itoa(seasonNumber) + strconv.Itoa(episodeNumber)
	inputEpisode.SeriesId = seriesId
	result := repo.DB.Where("id = ?", seriesId+strconv.Itoa(seasonNumber)+strconv.Itoa(episodeNumber)).First(&episode)

	if result.RecordNotFound() {
		episode = inputEpisode
		if err := repo.DB.Create(&episode).Error; err != nil {
			return models.Episode{}, err
		}
	} else {
		repo.DB.Model(&episode).Updates(inputEpisode)
		if err := repo.DB.Save(&episode).Error; err != nil {
			return models.Episode{}, err
		}
	}
	return episode, nil
}

func (repo *EpisodeRepository) GetEpisodeById(episodeId string) (models.Episode, error) {
	var episode models.Episode
	if err := repo.DB.Where("id = ?", episodeId).First(&episode).Error; err != nil {
		return models.Episode{}, err
	}
	return episode, nil
}

func (repo *EpisodeRepository) GetEpisodeBySeriesSeasonEpisode(seriesId string, seasonNumber int, episodeNumber int) (models.Episode, error) {
	var episode models.Episode
	if err := repo.DB.Where("series_id = ? AND season_number = ? AND episode_number = ?", seriesId, seasonNumber, episodeNumber).First(&episode).Error; err != nil {
		return models.Episode{}, err
	}
	return episode, nil
}

func (repo *EpisodeRepository) DeleteEpisodeById(seriesId string, seasonNumber int, episodeNumber int) error {
	var episode models.Episode
	if err := repo.DB.Where("series_id = ? AND season_number = ? AND episode_number = ?", seriesId, seasonNumber, episodeNumber).First(&episode).Error; err != nil {
		return err
	}

	db := repo.DB.Delete(&episode)
	if db.Error != nil {
		return db.Error
	}

	return nil
}
