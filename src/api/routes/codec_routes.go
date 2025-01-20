package routes

import (
	"net/http"
	"strings"
	"transfigurr/api/controllers"
	"transfigurr/interfaces"
)

func HandleCodecs(codecRepo interfaces.CodecRepositoryInterface) http.HandlerFunc {
	controller := controllers.NewCodecController(codecRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := strings.TrimPrefix(r.URL.Path, "/api/codecs/")
		segments := strings.Split(strings.Trim(path, "/"), "/")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		switch segments[0] {
		case "":
			controller.GetCodecs(w, r)
		case "containers":
			controller.GetContainers(w, r)
		case "encoders":
			controller.GetEncoders(w, r)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}
}
