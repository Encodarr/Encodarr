package handlers

import (
	"net/http"
	"strings"
	"transfigurr/internal/api/controllers"
	"transfigurr/internal/interfaces/repositories"
)

func HandleHistory(historyRepo repositories.HistoryRepositoryInterface) http.HandlerFunc {
	controller := controllers.NewHistoryController(historyRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := strings.TrimPrefix(r.URL.Path, "/api/history/")
		segments := strings.Split(strings.Trim(path, "/"), "/")

		switch {
		case r.Method == http.MethodGet && len(segments) == 1 && segments[0] == "":
			controller.GetHistories(w, r)
		case r.Method == http.MethodGet && len(segments) == 1:
			controller.GetHistoryById(w, r, segments[0])
		case r.Method == http.MethodPut && len(segments) == 1:
			controller.UpsertHistory(w, r, segments[0])
		case r.Method == http.MethodDelete && len(segments) == 1:
			controller.DeleteHistoryById(w, r, segments[0])
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
