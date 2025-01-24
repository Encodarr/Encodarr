package handlers

import (
	"net/http"
	"strings"
	"transfigurr/internal/api/controllers"
	"transfigurr/internal/interfaces/repositories"
)

func HandleSystem(systemRepo repositories.SystemRepositoryInterface) http.HandlerFunc {
	controller := controllers.NewSystemController(systemRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := strings.TrimPrefix(r.URL.Path, "/api/system/")
		segments := strings.Split(strings.Trim(path, "/"), "/")

		switch {
		case r.Method == http.MethodGet && len(segments) == 1 && segments[0] == "":
			controller.GetSystems(w, r)
		case r.Method == http.MethodGet && len(segments) == 1:
			controller.GetSystemById(w, r, segments[0])
		case r.Method == http.MethodPut && len(segments) == 1:
			controller.UpsertSystem(w, r, segments[0])
		case r.Method == http.MethodDelete && len(segments) == 1:
			controller.DeleteSystemById(w, r, segments[0])
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
