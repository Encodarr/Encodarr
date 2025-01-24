package handlers

import (
	"net/http"
	"strings"
	"transfigurr/internal/api/controllers"
)

func HandleArtwork() http.HandlerFunc {
	controller := controllers.NewArtworkController()

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/api/artwork/")
		segments := strings.Split(strings.Trim(path, "/"), "/")

		if len(segments) != 3 {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		switch {
		case segments[0] == "series" && segments[2] == "backdrop":
			controller.GetSeriesBackdrop(w, r, segments[1])
		case segments[0] == "series" && segments[2] == "poster":
			controller.GetSeriesPoster(w, r, segments[1])
		case segments[0] == "movies" && segments[2] == "backdrop":
			controller.GetMovieBackdrop(w, r, segments[1])
		case segments[0] == "movies" && segments[2] == "poster":
			controller.GetMoviePoster(w, r, segments[1])
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}
}
