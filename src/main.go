package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"transfigurr/startup"
)

func main() {
	sqlDB, services, repositories := startup.Startup()
	defer sqlDB.Close()

	// Create server mux
	mux := http.NewServeMux()

	// Setup routes
	SetupRouter(mux, services, repositories)

	// Configure server
	server := &http.Server{
		Addr:    "0.0.0.0:7889",
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
