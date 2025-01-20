package routes

import (
    "net/http"
    "strings"
    "transfigurr/api/controllers"
    "transfigurr/interfaces"
)

func HandleActions(scanService interfaces.ScanServiceInterface, metadataService interfaces.MetadataServiceInterface) http.HandlerFunc {
    controller := controllers.NewActionController(scanService, metadataService)

    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        path := strings.TrimPrefix(r.URL.Path, "/api/actions/")
        segments := strings.Split(strings.Trim(path, "/"), "/")

        switch {
        case len(segments) == 1 && segments[0] == "restart":
            controller.Restart(w, r)
        case len(segments) == 1 && segments[0] == "shutdown":
            controller.Shutdown(w, r)
        case len(segments) == 2 && segments[0] == "refresh" && segments[1] == "metadata":
            controller.RefreshMetadata(w, r)
        case len(segments) == 1 && segments[0] == "scan":
            controller.Scan(w, r)
        case len(segments) == 4 && segments[0] == "refresh" && segments[1] == "metadata" && segments[2] == "series":
            controller.RefreshSeriesMetadata(w, r, segments[3])
        case len(segments) == 3 && segments[0] == "scan" && segments[1] == "series":
            controller.ScanSeries(w, r, segments[2])
        case len(segments) == 4 && segments[0] == "refresh" && segments[1] == "metadata" && segments[2] == "movies":
            controller.RefreshMovieMetadata(w, r, segments[3])
        case len(segments) == 3 && segments[0] == "scan" && segments[1] == "movies":
            controller.ScanMovie(w, r, segments[2])
        default:
            http.Error(w, "Not found", http.StatusNotFound)
        }
    }
}