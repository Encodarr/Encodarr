package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"transfigurr/internal/interfaces/services"
	"transfigurr/internal/types"
)

type SSEController struct {
	EncodeService services.EncodeServiceInterface
	Repositories  *types.Repositories
}

func NewSSEController(encodeService services.EncodeServiceInterface, repositories *types.Repositories) *SSEController {
	return &SSEController{
		EncodeService: encodeService,
		Repositories:  repositories,
	}
}

func (ctrl *SSEController) HandleEventStream(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE with retry
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("retry", "3000")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go func() {
		<-ctx.Done()
		cancel()
	}()

	fetchers := map[string]func() (interface{}, error){
		"settings":   func() (interface{}, error) { return ctrl.Repositories.SettingRepo.GetAllSettings() },
		"system":     func() (interface{}, error) { return ctrl.Repositories.SystemRepo.GetSystems() },
		"profiles":   func() (interface{}, error) { return ctrl.Repositories.ProfileRepo.GetAllProfiles() },
		"containers": func() (interface{}, error) { return ctrl.Repositories.CodecRepo.GetContainers(), nil },
		"codecs":     func() (interface{}, error) { return ctrl.Repositories.CodecRepo.GetCodecs(), nil },
		"encoders":   func() (interface{}, error) { return ctrl.Repositories.CodecRepo.GetEncoders(), nil },
		"queue":      func() (interface{}, error) { return ctrl.EncodeService.GetQueue(), nil },
		"series":     func() (interface{}, error) { return ctrl.Repositories.SeriesRepo.GetSeries() },
		"movies":     func() (interface{}, error) { return ctrl.Repositories.MovieRepo.GetMovies() },
		"history":    func() (interface{}, error) { return ctrl.Repositories.HistoryRepo.GetHistories() },
		"logs":       func() (interface{}, error) { return ctrl.Repositories.EventRepo.GetEvents() },
	}

	ticker := time.NewTicker(250 * time.Millisecond)
	heartbeat := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	defer heartbeat.Stop()

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
			for eventType, fetcher := range fetchers {
				if data, err := fetcher(); err == nil {
					events[eventType] = data
				}
			}

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
