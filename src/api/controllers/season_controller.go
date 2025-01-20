package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"transfigurr/interfaces"
	"transfigurr/models"
)

type SeasonController struct {
	seasonRepo  interfaces.SeasonRepositoryInterface
	scanService interfaces.ScanServiceInterface
}

func NewSeasonController(seasonRepo interfaces.SeasonRepositoryInterface, scanService interfaces.ScanServiceInterface) *SeasonController {
	return &SeasonController{
		seasonRepo:  seasonRepo,
		scanService: scanService,
	}
}

func (c *SeasonController) GetSeasons(w http.ResponseWriter, r *http.Request, seriesId string) {
	seasons, err := c.seasonRepo.GetSeasons(seriesId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(seasons)
}

func (c *SeasonController) GetSeasonByID(w http.ResponseWriter, r *http.Request, seriesId string, seasonNumberString string) {
	seasonNumber, err := strconv.Atoi(seasonNumberString)
	if err != nil {
		http.Error(w, "Invalid season number", http.StatusBadRequest)
		return
	}
	season, err := c.seasonRepo.GetSeasonById(seriesId, seasonNumber)
	if err != nil {
		http.Error(w, "Season not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(season)
}

func (c *SeasonController) UpsertSeason(w http.ResponseWriter, r *http.Request, seriesId string, seasonNumberString string) {
	var season models.Season
	seasonNum, err := strconv.Atoi(seasonNumberString)
	if err != nil {
		http.Error(w, "Invalid season number", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&season); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := c.seasonRepo.UpsertSeason(string(seriesId), seasonNum, season); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(season)
}

func (c *SeasonController) DeleteSeasonByID(w http.ResponseWriter, r *http.Request, seriesId string, seasonNumberString string) {
	seasonNumber, err := strconv.Atoi(seasonNumberString)
	if err != nil {
		http.Error(w, "Invalid season number", http.StatusBadRequest)
		return
	}
	if err := c.seasonRepo.DeleteSeasonById(seriesId, seasonNumber); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
