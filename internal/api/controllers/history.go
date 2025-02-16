package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"transfigurr/internal/interfaces/repositories"
	"transfigurr/internal/models"
)

type HistoryController struct {
	Repo repositories.HistoryRepositoryInterface
}

func NewHistoryController(repo repositories.HistoryRepositoryInterface) *HistoryController {
	return &HistoryController{
		Repo: repo,
	}
}

func (ctrl HistoryController) GetHistories(w http.ResponseWriter, r *http.Request) {
	histories, err := ctrl.Repo.GetHistories()
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Histories not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving histories", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(histories)
}

func (ctrl HistoryController) GetHistoryById(w http.ResponseWriter, r *http.Request, historyId string) {
	history, err := ctrl.Repo.GetHistoryById(historyId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "History not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving history", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func (ctrl HistoryController) UpsertHistory(w http.ResponseWriter, r *http.Request, historyId string) {
	var inputHistory models.History
	if err := json.NewDecoder(r.Body).Decode(&inputHistory); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := ctrl.Repo.UpsertHistoryById(&inputHistory)
	if err != nil {
		http.Error(w, "Error upserting history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inputHistory)
}

func (ctrl HistoryController) DeleteHistoryById(w http.ResponseWriter, r *http.Request, historyId string) {
	var inputHistory models.History
	if err := json.NewDecoder(r.Body).Decode(&inputHistory); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := ctrl.Repo.DeleteHistoryById(&inputHistory)
	if err != nil {
		http.Error(w, "Error deleting history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "History deleted successfully"})
}
