package repositories

import (
	"errors"
	"transfigurr/internal/models"

	"gorm.io/gorm"
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

func (repo *MovieRepository) UpsertMovie(id string, movie models.Movie) (models.Movie, error) {
	var existingMovie models.Movie
	if err := repo.DB.Where("id = ?", id).First(&existingMovie).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If the movie does not exist, create a new one
			if err := repo.DB.Create(&movie).Error; err != nil {
				return models.Movie{}, err
			}
			return movie, nil
		}
		return models.Movie{}, err
	}

	// If the movie exists, update it
	if err := repo.DB.Model(&existingMovie).Select("*").Updates(movie).Error; err != nil {
		return models.Movie{}, err
	}
	return existingMovie, nil
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
