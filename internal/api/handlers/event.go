package handlers

import (
	"net/http"
	"strings"
	"transfigurr/internal/api/controllers"
	"transfigurr/internal/interfaces/repositories"
)

func HandleEvents(eventRepo repositories.EventRepositoryInterface) http.HandlerFunc {
	controller := controllers.NewEventController(eventRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := strings.TrimPrefix(r.URL.Path, "/api/events/")
		segments := strings.Split(strings.Trim(path, "/"), "/")

		switch {
		case r.Method == http.MethodGet && len(segments) == 1 && segments[0] == "":
			controller.GetEvents(w, r)
		case r.Method == http.MethodGet && len(segments) == 1:
			controller.GetEventById(w, r, segments[0])
		case r.Method == http.MethodPut && len(segments) == 1:
			controller.UpsertEvent(w, r, segments[0])
		case r.Method == http.MethodDelete && len(segments) == 1:
			controller.DeleteEventById(w, r, segments[0])
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
