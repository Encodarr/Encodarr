package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"transfigurr/internal/interfaces/repositories"
	"transfigurr/internal/models"

	"gorm.io/gorm"
)

type SettingController struct {
	Repo repositories.SettingRepositoryInterface
}

func NewSettingController(repo repositories.SettingRepositoryInterface) *SettingController {
	return &SettingController{
		Repo: repo,
	}
}

func (ctrl *SettingController) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := ctrl.Repo.GetAllSettings()
	if err != nil {
		http.Error(w, "Error retrieving settings", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(settings)
}

func (ctrl *SettingController) GetSettingById(w http.ResponseWriter, r *http.Request, settingId string) {
	setting, err := ctrl.Repo.GetSettingById(settingId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Setting not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving setting", http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(setting)
}

func (ctrl *SettingController) UpsertSetting(w http.ResponseWriter, r *http.Request, settingId string) {
	var inputSetting models.Setting
	if err := json.NewDecoder(r.Body).Decode(&inputSetting); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := ctrl.Repo.UpdateSetting(inputSetting)
	if err != nil {
		http.Error(w, "Error upserting setting", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(inputSetting)
}

func (ctrl *SettingController) DeleteSettingById(w http.ResponseWriter, r *http.Request, setting string) {
	var inputSetting models.Setting
	if err := json.NewDecoder(r.Body).Decode(&inputSetting); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := ctrl.Repo.DeleteSetting(inputSetting)
	if err != nil {
		http.Error(w, "Error deleting setting", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Setting deleted successfully"})
}
