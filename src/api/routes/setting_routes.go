package routes

import (
	"net/http"
	"strings"
	"transfigurr/api/controllers"
	"transfigurr/interfaces"
)

func HandleSettings(settingRepo interfaces.SettingRepositoryInterface) http.HandlerFunc {
	controller := controllers.NewSettingController(settingRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := strings.TrimPrefix(r.URL.Path, "/api/settings/")
		segments := strings.Split(strings.Trim(path, "/"), "/")

		switch {
		case r.Method == http.MethodGet && len(segments) == 1 && segments[0] == "":
			controller.GetSettings(w, r)
		case r.Method == http.MethodGet && len(segments) == 1:
			controller.GetSettingById(w, r, segments[0])
		case r.Method == http.MethodPut && len(segments) == 1:
			controller.UpsertSetting(w, r, segments[0])
		case r.Method == http.MethodDelete && len(segments) == 1:
			controller.DeleteSettingById(w, r, segments[0])
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
