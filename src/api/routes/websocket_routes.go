package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"transfigurr/interfaces"
	"transfigurr/types"
)

func HandleEventStream(encodeService interfaces.EncodeServiceInterface, repositories *types.Repositories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set headers for SSE with retry
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-Accel-Buffering", "no") // Disable proxy buffering
		w.Header().Set("retry", "3000")           // Tell client to retry after 3s

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}

		// Create context with cancellation
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		// Handle client disconnect
		go func() {
			<-ctx.Done()
			cancel()
		}()

		// Create fetchers map with error handling
		fetchers := map[string]func() (interface{}, error){
			"settings":   func() (interface{}, error) { return repositories.SettingRepo.GetAllSettings() },
			"system":     func() (interface{}, error) { return repositories.SystemRepo.GetSystems() },
			"profiles":   func() (interface{}, error) { return repositories.ProfileRepo.GetAllProfiles() },
			"containers": func() (interface{}, error) { return repositories.CodecRepo.GetContainers(), nil },
			"codecs":     func() (interface{}, error) { return repositories.CodecRepo.GetCodecs(), nil },
			"encoders":   func() (interface{}, error) { return repositories.CodecRepo.GetEncoders(), nil },
			"queue":      func() (interface{}, error) { return encodeService.GetQueue(), nil },
			"series":     func() (interface{}, error) { return repositories.SeriesRepo.GetSeries() },
			"movies":     func() (interface{}, error) { return repositories.MovieRepo.GetMovies() },
			"history":    func() (interface{}, error) { return repositories.HistoryRepo.GetHistories() },
			"logs":       func() (interface{}, error) { return repositories.EventRepo.GetEvents() },
		}

		// Start event loop with heartbeat
		ticker := time.NewTicker(250 * time.Millisecond)
		heartbeat := time.NewTicker(15 * time.Second)
		defer ticker.Stop()
		defer heartbeat.Stop()

		// Use mutex for thread-safe writes
		var mu sync.Mutex

		for {
			select {
			case <-ctx.Done():
				return
			case <-heartbeat.C:
				mu.Lock()
				fmt.Fprintf(w, ": heartbeat\n\n")
				flusher.Flush()
				mu.Unlock()
			case <-ticker.C:
				events := make(map[string]interface{})

				// Batch fetch all data
				for eventType, fetcher := range fetchers {
					if data, err := fetcher(); err == nil {
						events[eventType] = data
					}
				}

				// Send batch update
				if len(events) > 0 {
					mu.Lock()
					for eventType, data := range events {
						if jsonData, err := json.Marshal(data); err == nil {
							fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, jsonData)
						}
					}
					flusher.Flush()
					mu.Unlock()
				}
			}
		}
	}
}
