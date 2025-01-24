package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"transfigurr/internal/interfaces/services"
	"transfigurr/internal/models"
)

type ActionController struct {
	scanService     services.ScanServiceInterface
	metadataService services.MetadataServiceInterface
}

func NewActionController(scanService services.ScanServiceInterface, metadataService services.MetadataServiceInterface) *ActionController {
	return &ActionController{
		scanService:     scanService,
		metadataService: metadataService,
	}
}

func (ctrl ActionController) Restart(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	if err := os.Chtimes(filepath.Join("/config", "restart.txt"), now, now); err != nil {
		http.Error(w, "Failed to trigger restart", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Application is restarting"})
}

func (ctrl ActionController) Shutdown(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	if err := os.Chtimes(filepath.Join("/config", "shutdown.txt"), now, now); err != nil {
		http.Error(w, "Failed to trigger shutdown", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Application is shutting down"})
}

func (ctrl ActionController) RefreshMetadata(w http.ResponseWriter, r *http.Request) {
	ctrl.metadataService.EnqueueAll()
	json.NewEncoder(w).Encode(map[string]string{"message": "Metadata refresh enqueued"})
}

func (ctrl ActionController) Scan(w http.ResponseWriter, r *http.Request) {
	ctrl.scanService.EnqueueAll()
	json.NewEncoder(w).Encode(map[string]string{"message": "Scan enqueued"})
}

func (ctrl ActionController) RefreshSeriesMetadata(w http.ResponseWriter, r *http.Request, seriesId string) {
	ctrl.metadataService.Enqueue(models.Item{Id: seriesId, Type: "series"})
	json.NewEncoder(w).Encode(map[string]string{"message": "Refresh enqueued"})
}

func (ctrl ActionController) ScanSeries(w http.ResponseWriter, r *http.Request, seriesId string) {
	ctrl.scanService.Enqueue(models.Item{Id: seriesId, Type: "series"})
	json.NewEncoder(w).Encode(map[string]string{"message": "Scan enqueued"})
}

func (ctrl ActionController) RefreshMovieMetadata(w http.ResponseWriter, r *http.Request, movieId string) {
	ctrl.metadataService.Enqueue(models.Item{Id: movieId, Type: "movie"})
	json.NewEncoder(w).Encode(map[string]string{"message": "Refresh enqueued"})
}

func (ctrl ActionController) ScanMovie(w http.ResponseWriter, r *http.Request, movieId string) {
	ctrl.scanService.Enqueue(models.Item{Id: movieId, Type: "movie"})
	json.NewEncoder(w).Encode(map[string]string{"message": "Scan enqueued"})
}
