package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"transfigurr/internal/interfaces/repositories"
	"transfigurr/internal/interfaces/services"
	"transfigurr/internal/models"

	"gorm.io/gorm"
)

type MovieController struct {
	Repo        repositories.MovieRepositoryInterface
	ScanService services.ScanServiceInterface
}

func NewMovieController(repo repositories.MovieRepositoryInterface, scanService services.ScanServiceInterface) *MovieController {
	return &MovieController{
		Repo:        repo,
		ScanService: scanService,
	}
}

func (ctrl *MovieController) GetMovies(w http.ResponseWriter, r *http.Request) {
	movieList, err := ctrl.Repo.GetMovies()
	if err != nil {
		http.Error(w, "Error retrieving movies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movieList)
}

func (ctrl *MovieController) UpsertMovie(w http.ResponseWriter, r *http.Request, movieId string) {
	var inputMovie models.Movie

	if err := json.NewDecoder(r.Body).Decode(&inputMovie); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	movie, err := ctrl.Repo.UpsertMovie(movieId, inputMovie)
	if err != nil {
		http.Error(w, "Error upserting movie", http.StatusInternalServerError)
		return
	}
	ctrl.ScanService.Enqueue(models.Item{Id: movie.Id, Type: "movie"})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

func (ctrl *MovieController) GetMovieByID(w http.ResponseWriter, r *http.Request, movieId string) {
	movie, err := ctrl.Repo.GetMovieById(movieId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Movie not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving movie", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

func (ctrl *MovieController) DeleteMovieByID(w http.ResponseWriter, r *http.Request, movieId string) {
	err := ctrl.Repo.DeleteMovieById(movieId)
	if err != nil {
		http.Error(w, "Error deleting movie", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Movie deleted successfully"})
}
