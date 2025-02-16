package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"transfigurr/internal/config"
	"transfigurr/internal/router"
	"transfigurr/internal/scripts/startup"
)

func main() {
	db, services, repositories := startup.Startup()
	defer db.Close()

	// Create server mux
	mux := http.NewServeMux()

	// Setup handlers
	router.SetupRouter(mux, services, repositories)

	// Configure server
	server := &http.Server{
		Addr:    config.ServerAddress,
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
	}()

	// Start server
	log.Printf("Server starting on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server error:", err)
	}
}
