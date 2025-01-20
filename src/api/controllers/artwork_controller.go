package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"transfigurr/constants"
)

type ArtworkController struct{}

func NewArtworkController() *ArtworkController {
	return &ArtworkController{}
}

func (a *ArtworkController) GetSeriesBackdrop(w http.ResponseWriter, r *http.Request, seriesId string) {
	filePath := filepath.Join(constants.ArtworkPath, "series", seriesId, "backdrop.webp")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Backdrop not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filePath)
}

func (a *ArtworkController) GetSeriesPoster(w http.ResponseWriter, r *http.Request, seriesId string) {
	filePath := filepath.Join(constants.ArtworkPath, "series", seriesId, "poster.webp")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Poster not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filePath)
}

func (a *ArtworkController) GetMovieBackdrop(w http.ResponseWriter, r *http.Request, movieId string) {
	filePath := filepath.Join(constants.ArtworkPath, "movies", movieId, "backdrop.webp")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Backdrop not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filePath)
}

func (a *ArtworkController) GetMoviePoster(w http.ResponseWriter, r *http.Request, movieId string) {
	filePath := filepath.Join(constants.ArtworkPath, "movies", movieId, "poster.webp")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Poster not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filePath)
}
