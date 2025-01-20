package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"transfigurr/interfaces"
	"transfigurr/models"

	"gorm.io/gorm"
)

type SystemController struct {
	Repo interfaces.SystemRepositoryInterface
}

func NewSystemController(repo interfaces.SystemRepositoryInterface) *SystemController {
	return &SystemController{
		Repo: repo,
	}
}

func (ctrl *SystemController) GetSystems(w http.ResponseWriter, r *http.Request) {
	systems, err := ctrl.Repo.GetSystems()
	if err != nil {
		http.Error(w, "Error retrieving systems", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(systems)
}

func (ctrl *SystemController) GetSystemById(w http.ResponseWriter, r *http.Request, systemId string) {
	system, err := ctrl.Repo.GetSystemById(systemId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "System not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving system", http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(system)
}

func (ctrl *SystemController) UpsertSystem(w http.ResponseWriter, r *http.Request, systemId string) {
	var inputSystem models.System
	if err := json.NewDecoder(r.Body).Decode(&inputSystem); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	system, err := ctrl.Repo.UpsertSystem(systemId, inputSystem)
	if err != nil {
		http.Error(w, "Error upserting system", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(system)
}

func (ctrl *SystemController) DeleteSystemById(w http.ResponseWriter, r *http.Request, systemId string) {
	err := ctrl.Repo.DeleteSystemById(systemId)
	if err != nil {
		http.Error(w, "Error deleting system", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "System deleted successfully"})
}
