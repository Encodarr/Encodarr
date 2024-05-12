package repository

import (
	"transfigurr/models"

	"github.com/jinzhu/gorm"
)

type MovieRepository struct {
	DB *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{
		DB: db,
	}
}

func (repo *MovieRepository) GetMovies() ([]models.Movie, error) {
	var movieList []models.Movie
	if err := repo.DB.Find(&movieList).Error; err != nil {
		return nil, err
	}
	return movieList, nil
}

func (repo *MovieRepository) UpsertMovie(id string, inputMovie models.Movie) (models.Movie, error) {
	var movie models.Movie
	inputMovie.Id = id
	result := repo.DB.Where("id = ?", id).First(&movie)

	if result.RecordNotFound() {
		movie = inputMovie
		if err := repo.DB.Create(&movie).Error; err != nil {
			return models.Movie{}, err
		}
	} else {
		repo.DB.Model(&movie).Updates(inputMovie)
		if err := repo.DB.Save(&movie).Error; err != nil {
			return models.Movie{}, err
		}
	}
	return movie, nil
}

func (repo *MovieRepository) GetMovieById(id string) (models.Movie, error) {
	var movie models.Movie
	if err := repo.DB.Where("id = ?", id).First(&movie).Error; err != nil {
		return models.Movie{}, err
	}
	return movie, nil
}

func (repo *MovieRepository) DeleteMovieById(id string) error {
	var movie models.Movie
	if err := repo.DB.Where("id = ?", id).First(&movie).Error; err != nil {
		return err
	}

	db := repo.DB.Delete(&movie)
	if db.Error != nil {
		return db.Error
	}

	return nil
}
