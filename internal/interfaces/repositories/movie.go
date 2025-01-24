package repositories

import "transfigurr/internal/models"

type MovieRepositoryInterface interface {
	GetMovies() ([]models.Movie, error)
	UpsertMovie(id string, inputMovie models.Movie) (models.Movie, error)
	GetMovieById(id string) (models.Movie, error)
	DeleteMovieById(id string) error
}
