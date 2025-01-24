package handlers

import (
	"net/http"
	"strings"
	"transfigurr/internal/api/controllers"
	"transfigurr/internal/interfaces/repositories"
	"transfigurr/internal/interfaces/services"
)

func HandleSeries(seriesRepo repositories.SeriesRepositoryInterface, seasonRepo repositories.SeasonRepositoryInterface, episodeRepo repositories.EpisodeRepositoryInterface, scanService services.ScanServiceInterface) http.HandlerFunc {
	seriesController := controllers.NewSeriesController(seriesRepo, scanService)
	seasonController := controllers.NewSeasonController(seasonRepo, scanService)
	episodeController := controllers.NewEpisodeController(episodeRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Parse URL path
		path := strings.TrimPrefix(r.URL.Path, "/api/series/")
		segments := strings.Split(strings.Trim(path, "/"), "/")

		switch {
		// Series endpoints
		case r.Method == http.MethodGet && len(segments) == 1 && segments[0] == "":
			seriesController.GetSeries(w, r)
		case r.Method == http.MethodGet && len(segments) == 1:
			seriesController.GetSeriesByID(w, r, segments[0])
		case r.Method == http.MethodPut && len(segments) == 1:
			seriesController.UpsertSeries(w, r, segments[0])
		case r.Method == http.MethodDelete && len(segments) == 1:
			seriesController.DeleteSeriesByID(w, r, segments[0])

		// Season endpoints
		case r.Method == http.MethodGet && len(segments) == 2 && segments[1] == "seasons":
			seasonController.GetSeasons(w, r, segments[0])
		case r.Method == http.MethodGet && len(segments) == 3 && segments[1] == "seasons":
			seasonController.GetSeasonByID(w, r, segments[0], segments[2])
		case r.Method == http.MethodPut && len(segments) == 3 && segments[1] == "seasons":
			seasonController.UpsertSeason(w, r, segments[0], segments[2])
		case r.Method == http.MethodDelete && len(segments) == 3 && segments[1] == "seasons":
			seasonController.DeleteSeasonByID(w, r, segments[0], segments[2])

		// Episode endpoints
		case r.Method == http.MethodGet && len(segments) == 4 && segments[1] == "seasons" && segments[3] == "episodes":
			episodeController.GetEpisodes(w, r, segments[0], segments[2])
		case r.Method == http.MethodGet && len(segments) == 5 && segments[1] == "seasons" && segments[3] == "episodes":
			episodeController.GetEpisodeBySeriesSeasonEpisode(w, r, segments[0], segments[2], segments[4])
		case r.Method == http.MethodPut && len(segments) == 5 && segments[1] == "seasons" && segments[3] == "episodes":
			episodeController.UpsertEpisode(w, r, segments[0], segments[2], segments[4])
		case r.Method == http.MethodDelete && len(segments) == 5 && segments[1] == "seasons" && segments[3] == "episodes":
			episodeController.DeleteEpisodeById(w, r, segments[0], segments[2], segments[4])
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
