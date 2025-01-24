package middleware

import (
	"net/http"
	"strings"
	"transfigurr/pkg/jwt"
)

func AuthMiddleware(next http.HandlerFunc, jwtSecret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		_, err := jwt.Parse(tokenString, jwtSecret)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func Protected(handler http.HandlerFunc, jwtSecret []byte) http.HandlerFunc {
	return AuthMiddleware(handler, jwtSecret)
}
