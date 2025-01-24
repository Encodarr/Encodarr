package handlers

import (
	"net/http"
	"transfigurr/internal/api/controllers"
	"transfigurr/internal/interfaces/services"
	"transfigurr/internal/types"
)

func HandleEventStream(encodeService services.EncodeServiceInterface, repositories *types.Repositories) http.HandlerFunc {
	controller := controllers.NewSSEController(encodeService, repositories)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		controller.HandleEventStream(w, r)
	}
}
