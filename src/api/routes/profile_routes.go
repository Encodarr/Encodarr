package routes

import (
	"net/http"
	"strings"
	"transfigurr/api/controllers"
	"transfigurr/interfaces"
)

func HandleProfiles(scanService interfaces.ScanServiceInterface, profileRepo interfaces.ProfileRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, seriesRepo interfaces.SeriesRepositoryInterface) http.HandlerFunc {
	controller := controllers.NewProfileController(scanService, profileRepo, movieRepo, seriesRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := strings.TrimPrefix(r.URL.Path, "/api/profiles/")
		segments := strings.Split(strings.Trim(path, "/"), "/")

		switch {
		case r.Method == http.MethodGet && len(segments) == 1 && segments[0] == "":
			controller.GetProfiles(w, r)
		case r.Method == http.MethodGet && len(segments) == 1:
			controller.GetProfileById(w, r, segments[0])
		case r.Method == http.MethodPut && len(segments) == 1:
			controller.UpsertProfile(w, r, segments[0])
		case r.Method == http.MethodDelete && len(segments) == 1:
			controller.DeleteProfileById(w, r, segments[0])
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
