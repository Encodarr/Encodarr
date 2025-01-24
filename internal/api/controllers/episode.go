package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"transfigurr/internal/interfaces/repositories"
	"transfigurr/internal/models"
)

type EpisodeController struct {
	Repo repositories.EpisodeRepositoryInterface
}

func NewEpisodeController(repo repositories.EpisodeRepositoryInterface) *EpisodeController {
	return &EpisodeController{
		Repo: repo,
	}
}

func (ctrl *EpisodeController) GetEpisodes(w http.ResponseWriter, r *http.Request, seriesId string, seasonNumber string) {
	seasonNum, err := strconv.Atoi(seasonNumber)
	if err != nil {
		http.Error(w, "Invalid season number", http.StatusBadRequest)
		return
	}

	episodes, err := ctrl.Repo.GetEpisodes(seriesId, seasonNum)
	if err != nil {
		http.Error(w, "Error retrieving episodes", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(episodes)
}

func (ctrl *EpisodeController) UpsertEpisode(w http.ResponseWriter, r *http.Request, seriesId string, seasonNumber string, episodeNumber string) {
	var inputEpisode models.Episode

	seasonNum, err := strconv.Atoi(seasonNumber)
	if err != nil {
		http.Error(w, "Invalid season number", http.StatusBadRequest)
		return
	}

	episodeNum, err := strconv.Atoi(episodeNumber)
	if err != nil {
		http.Error(w, "Invalid episode number", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&inputEpisode); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	episode, err := ctrl.Repo.UpsertEpisode(seriesId, seasonNum, episodeNum, inputEpisode)
	if err != nil {
		http.Error(w, "Error upserting episode", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(episode)
}

func (ctrl *EpisodeController) GetEpisodeBySeriesSeasonEpisode(w http.ResponseWriter, r *http.Request, seriesId string, seasonNumber string, episodeNumber string) {
	seasonNum, err := strconv.Atoi(seasonNumber)
	if err != nil {
		http.Error(w, "Invalid season number", http.StatusBadRequest)
		return
	}

	episodeNum, err := strconv.Atoi(episodeNumber)
	if err != nil {
		http.Error(w, "Invalid episode number", http.StatusBadRequest)
		return
	}

	episode, err := ctrl.Repo.GetEpisodeBySeriesSeasonEpisode(seriesId, seasonNum, episodeNum)
	if err != nil {
		http.Error(w, "Error retrieving episode", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(episode)
}

func (ctrl *EpisodeController) DeleteEpisodeById(w http.ResponseWriter, r *http.Request, seriesId string, seasonNumber string, episodeNumber string) {
	seasonNum, err := strconv.Atoi(seasonNumber)
	if err != nil {
		http.Error(w, "Invalid season number", http.StatusBadRequest)
		return
	}

	episodeNum, err := strconv.Atoi(episodeNumber)
	if err != nil {
		http.Error(w, "Invalid episode number", http.StatusBadRequest)
		return
	}

	err = ctrl.Repo.DeleteEpisodeById(seriesId, seasonNum, episodeNum)
	if err != nil {
		http.Error(w, "Error deleting episode", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Episode deleted successfully"})
}
