package router

import (
	"net/http"
	"transfigurr/internal/api/handlers"
	"transfigurr/internal/router/middleware"
	"transfigurr/internal/types"
)

func SetupRouter(mux *http.ServeMux, services *types.Services, repositories *types.Repositories) {
	// API Routes prefix handler
	// apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//     // CORS middleware could go here
	//     w.Header().Set("Access-Control-Allow-Origin", "*")``
	//     w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	//     w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	//     if r.Method == "OPTIONS" {
	//         w.WriteHeader(http.StatusOK)
	//         return
	//     }
	// })

	user, err := repositories.UserRepo.GetUser()
	if err != nil {
		panic(err)
	}
	jwtSecret := []byte(user.Secret)

	// Public Routes
	mux.Handle("/api/auth/", http.HandlerFunc(handlers.HandleAuth(repositories.AuthRepo)))
	assetsHandler, rootHandler := handlers.HandleStatic("./frontend/dist")
	mux.Handle("/assets/", http.HandlerFunc(assetsHandler))
	mux.Handle("/", http.HandlerFunc(rootHandler))

	// Protected Routes
	mux.Handle("/api/series/", middleware.Protected(http.HandlerFunc(handlers.HandleSeries(repositories.SeriesRepo, repositories.SeasonRepo, repositories.EpisodeRepo, services.ScanService)), jwtSecret))
	mux.Handle("/api/movies/", middleware.Protected(http.HandlerFunc(handlers.HandleMovies(services.ScanService, repositories.MovieRepo)), jwtSecret))
	mux.Handle("/api/settings/", middleware.Protected(http.HandlerFunc(handlers.HandleSettings(repositories.SettingRepo)), jwtSecret))
	mux.Handle("/api/system/", middleware.Protected(http.HandlerFunc(handlers.HandleSystem(repositories.SystemRepo)), jwtSecret))
	mux.Handle("/api/profiles/", middleware.Protected(http.HandlerFunc(handlers.HandleProfiles(services.ScanService, repositories.ProfileRepo, repositories.MovieRepo, repositories.SeriesRepo)), jwtSecret))
	mux.Handle("/api/history/", middleware.Protected(http.HandlerFunc(handlers.HandleHistory(repositories.HistoryRepo)), jwtSecret))
	mux.Handle("/api/events/", middleware.Protected(http.HandlerFunc(handlers.HandleEvents(repositories.EventRepo)), jwtSecret))
	mux.Handle("/api/codecs/", middleware.Protected(http.HandlerFunc(handlers.HandleCodecs(repositories.CodecRepo)), jwtSecret))
	mux.Handle("/api/actions/", middleware.Protected(http.HandlerFunc(handlers.HandleActions(services.ScanService, services.MetadataService)), jwtSecret))
	mux.Handle("/api/artwork/", middleware.Protected(http.HandlerFunc(handlers.HandleArtwork()), jwtSecret))
	mux.Handle("/api/events/stream", middleware.Protected(handlers.HandleEventStream(services.EncodeService, repositories), jwtSecret))

}
