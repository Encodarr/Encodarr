package handlers

import (
	"net/http"
	"strings"
	"transfigurr/internal/api/controllers"
	"transfigurr/internal/interfaces/repositories"
	"transfigurr/internal/interfaces/services"
)

func HandleMovies(scanService services.ScanServiceInterface, movieRepo repositories.MovieRepositoryInterface) http.HandlerFunc {
	controller := controllers.NewMovieController(movieRepo, scanService)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Parse URL path
		path := strings.TrimPrefix(r.URL.Path, "/api/movies/")
		segments := strings.Split(strings.Trim(path, "/"), "/")

		switch {
		case r.Method == http.MethodGet && len(segments) == 1 && segments[0] == "":
			controller.GetMovies(w, r)
		case r.Method == http.MethodGet && len(segments) == 1:
			controller.GetMovieByID(w, r, segments[0])
		case r.Method == http.MethodPut && len(segments) == 1:
			controller.UpsertMovie(w, r, segments[0])
		case r.Method == http.MethodDelete && len(segments) == 1:
			controller.DeleteMovieByID(w, r, segments[0])
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
