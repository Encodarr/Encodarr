package routes

import (
	"net/http"
	"strings"
	"transfigurr/api/controllers"
	"transfigurr/interfaces"
)

func HandleUsers(userRepo interfaces.UserRepositoryInterface) http.HandlerFunc {
	controller := controllers.NewUserController(userRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := strings.TrimPrefix(r.URL.Path, "/api/users/")
		segments := strings.Split(strings.Trim(path, "/"), "/")

		switch {
		case r.Method == http.MethodGet && len(segments) == 1 && segments[0] == "":
			controller.GetUsers(w, r)
		case r.Method == http.MethodPost && len(segments) == 1 && segments[0] == "":
			//controller.UpdateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
