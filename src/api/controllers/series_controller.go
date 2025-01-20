package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"transfigurr/interfaces"
	"transfigurr/models"
)

type SeriesController struct {
	Repo        interfaces.SeriesRepositoryInterface
	ScanService interfaces.ScanServiceInterface
}

func NewSeriesController(repo interfaces.SeriesRepositoryInterface, scanService interfaces.ScanServiceInterface) *SeriesController {
	return &SeriesController{
		Repo:        repo,
		ScanService: scanService,
	}
}

func (ctrl *SeriesController) GetSeries(w http.ResponseWriter, r *http.Request) {
	seriesList, err := ctrl.Repo.GetSeries()
	if err != nil {
		http.Error(w, "Error retrieving series", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(seriesList)
}

func (ctrl *SeriesController) UpsertSeries(w http.ResponseWriter, r *http.Request, id string) {
	var inputSeries models.Series

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Parse JSON body
	if err := json.Unmarshal(body, &inputSeries); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	series, err := ctrl.Repo.UpsertSeries(id, inputSeries)
	if err != nil {
		http.Error(w, "Error upserting series", http.StatusInternalServerError)
		return
	}

	ctrl.ScanService.Enqueue(models.Item{Id: series.Id, Type: "series"})
	json.NewEncoder(w).Encode(series)
}

func (ctrl *SeriesController) GetSeriesByID(w http.ResponseWriter, r *http.Request, id string) {
	series, err := ctrl.Repo.GetSeriesByID(id)
	if err != nil {
		http.Error(w, "Error retrieving series", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(series)
}

func (ctrl *SeriesController) DeleteSeriesByID(w http.ResponseWriter, r *http.Request, id string) {
	err := ctrl.Repo.DeleteSeriesByID(id)
	if err != nil {
		http.Error(w, "Error deleting series", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Series deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
