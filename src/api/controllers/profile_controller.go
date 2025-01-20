package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"transfigurr/interfaces"
	"transfigurr/models"
)

type ProfileController struct {
	ProfileRepo interfaces.ProfileRepositoryInterface
	MovieRepo   interfaces.MovieRepositoryInterface
	SeriesRepo  interfaces.SeriesRepositoryInterface
	ScanService interfaces.ScanServiceInterface
}

func NewProfileController(scanService interfaces.ScanServiceInterface, profileRepo interfaces.ProfileRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, seriesRepo interfaces.SeriesRepositoryInterface) *ProfileController {
	return &ProfileController{
		ProfileRepo: profileRepo,
		MovieRepo:   movieRepo,
		SeriesRepo:  seriesRepo,
		ScanService: scanService,
	}
}

func (ctrl *ProfileController) GetProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := ctrl.ProfileRepo.GetAllProfiles()
	if err != nil {
		http.Error(w, "Error retrieving profiles", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(profiles)
}

func (ctrl *ProfileController) GetProfileById(w http.ResponseWriter, r *http.Request, profileId string) {
	profileIdInt, err := strconv.Atoi(profileId)
	if err != nil {
		http.Error(w, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	profile, err := ctrl.ProfileRepo.GetProfileById(profileIdInt)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func (ctrl *ProfileController) UpsertProfile(w http.ResponseWriter, r *http.Request, profileId string) {
	var inputProfile models.Profile

	profileIdInt, err := strconv.Atoi(profileId)
	if err != nil {
		http.Error(w, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&inputProfile); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	profile, err := ctrl.ProfileRepo.UpsertProfile(profileIdInt, inputProfile)
	if err != nil {
		http.Error(w, "Error updating profile", http.StatusInternalServerError)
		return
	}

	movies, err := ctrl.MovieRepo.GetMovies()
	if err != nil {
		http.Error(w, "Error retrieving movies", http.StatusInternalServerError)
		return
	}

	series, err := ctrl.SeriesRepo.GetSeries()
	if err != nil {
		http.Error(w, "Error retrieving series", http.StatusInternalServerError)
		return
	}

	for _, movie := range movies {
		if movie.ProfileID == profileIdInt {
			ctrl.ScanService.Enqueue(models.Item{Id: movie.Id, Type: "movie"})
		}
	}

	for _, serie := range series {
		if serie.ProfileID == profileIdInt {
			ctrl.ScanService.Enqueue(models.Item{Id: serie.Id, Type: "series"})
		}
	}

	json.NewEncoder(w).Encode(profile)
}

func (ctrl *ProfileController) DeleteProfileById(w http.ResponseWriter, r *http.Request, profileId string) {
	profileIdInt, err := strconv.Atoi(profileId)
	if err != nil {
		http.Error(w, "Invalid profile ID", http.StatusBadRequest)
		return
	}

	err = ctrl.ProfileRepo.DeleteProfileById(profileIdInt)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Profile deleted successfully"})
}
