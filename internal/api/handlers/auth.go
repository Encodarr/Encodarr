package handlers

import (
	"net/http"
	"strings"
	"transfigurr/internal/api/controllers"
	"transfigurr/internal/interfaces/repositories"
)

func HandleAuth(authRepo repositories.AuthRepositoryInterface) http.HandlerFunc {
	controller := controllers.NewAuthController(authRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := strings.TrimPrefix(r.URL.Path, "/api/auth/")
		segments := strings.Split(strings.Trim(path, "/"), "/")

		switch {
		case r.Method == http.MethodGet && segments[0] == "activated":
			controller.GetActivated(w, r)
		case r.Method == http.MethodPost && segments[0] == "register":
			controller.Register(w, r)
		case r.Method == http.MethodPost && segments[0] == "login":
			controller.Login(w, r)
		case r.Method == http.MethodPost && segments[0] == "logintoken":
			controller.LoginToken(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
