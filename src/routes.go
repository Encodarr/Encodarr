package main

import (
	"net/http"
	"transfigurr/api/middleware"
	"transfigurr/api/routes"
	"transfigurr/types"
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
	mux.Handle("/api/auth/", http.HandlerFunc(routes.HandleAuth(repositories.AuthRepo)))
	assetsHandler, rootHandler := routes.HandleStatic("./frontend/dist")
	mux.Handle("/assets/", http.HandlerFunc(assetsHandler))
	mux.Handle("/", http.HandlerFunc(rootHandler))

	// Protected Routes
	mux.Handle("/api/series/", middleware.Protected(http.HandlerFunc(routes.HandleSeries(repositories.SeriesRepo, repositories.SeasonRepo, repositories.EpisodeRepo, services.ScanService)), jwtSecret))
	mux.Handle("/api/movies/", middleware.Protected(http.HandlerFunc(routes.HandleMovies(services.ScanService, repositories.MovieRepo)), jwtSecret))
	mux.Handle("/api/settings/", middleware.Protected(http.HandlerFunc(routes.HandleSettings(repositories.SettingRepo)), jwtSecret))
	mux.Handle("/api/system/", middleware.Protected(http.HandlerFunc(routes.HandleSystem(repositories.SystemRepo)), jwtSecret))
	mux.Handle("/api/profiles/", middleware.Protected(http.HandlerFunc(routes.HandleProfiles(services.ScanService, repositories.ProfileRepo, repositories.MovieRepo, repositories.SeriesRepo)), jwtSecret))
	mux.Handle("/api/users/", middleware.Protected(http.HandlerFunc(routes.HandleUsers(repositories.UserRepo)), jwtSecret))
	mux.Handle("/api/history/", middleware.Protected(http.HandlerFunc(routes.HandleHistory(repositories.HistoryRepo)), jwtSecret))
	mux.Handle("/api/events/", middleware.Protected(http.HandlerFunc(routes.HandleEvents(repositories.EventRepo)), jwtSecret))
	mux.Handle("/api/codecs/", middleware.Protected(http.HandlerFunc(routes.HandleCodecs(repositories.CodecRepo)), jwtSecret))
	mux.Handle("/api/actions/", middleware.Protected(http.HandlerFunc(routes.HandleActions(services.ScanService, services.MetadataService)), jwtSecret))
	mux.Handle("/api/artwork/", middleware.Protected(http.HandlerFunc(routes.HandleArtwork()), jwtSecret))
	mux.Handle("/api/events/stream", middleware.Protected(routes.HandleEventStream(services.EncodeService, repositories), jwtSecret))

}
