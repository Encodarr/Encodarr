package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"transfigurr/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WebsocketRoutes(rg *gin.RouterGroup, encodeService interfaces.EncodeServiceInterface, seriesRepo interfaces.SeriesRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface, historyRepo interfaces.HistoryRepositoryInterface, eventRepo interfaces.EventRepositoryInterface, codecRepo interfaces.CodecRepositoryInterface) {
	rg.GET("/ws", func(c *gin.Context) {
		WebsocketHandler(c.Writer, c.Request, encodeService, seriesRepo, movieRepo, profileRepo, settingRepo, systemRepo, historyRepo, eventRepo, codecRepo)
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request, encodeService interfaces.EncodeServiceInterface, seriesRepo interfaces.SeriesRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface, historyRepo interfaces.HistoryRepositoryInterface, eventRepo interfaces.EventRepositoryInterface, codecRepo interfaces.CodecRepositoryInterface) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	// Set up ping/pong handlers
	conn.SetPingHandler(func(appData string) error {
		log.Println("Received ping")
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})

	conn.SetPongHandler(func(appData string) error {
		log.Println("Received pong")
		return nil
	})

	// Create a channel to signal when to stop the goroutines
	stopChan := make(chan struct{})
	defer close(stopChan)

	// Create a mutex for synchronized writes
	var writeMutex sync.Mutex

	// Start goroutines to fetch and send data concurrently
	go startDataFetchers(conn, seriesRepo, movieRepo, profileRepo, settingRepo, systemRepo, historyRepo, eventRepo, codecRepo, encodeService, stopChan, &writeMutex)

	// Read messages to keep the connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
	}
	log.Println("WebSocket connection closed")
}

func startDataFetchers(conn *websocket.Conn, seriesRepo interfaces.SeriesRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface, historyRepo interfaces.HistoryRepositoryInterface, eventRepo interfaces.EventRepositoryInterface, codecRepo interfaces.CodecRepositoryInterface, encodeService interfaces.EncodeServiceInterface, stopChan chan struct{}, writeMutex *sync.Mutex) {
	type fetcher struct {
		dataType string
		getter   func() (interface{}, error)
	}

	fetchers := []fetcher{
		{"settings", func() (interface{}, error) { return settingRepo.GetAllSettings() }},
		{"system", func() (interface{}, error) { return systemRepo.GetSystems() }},
		{"profiles", func() (interface{}, error) { return profileRepo.GetAllProfiles() }},
		{"containers", func() (interface{}, error) { return codecRepo.GetContainers(), nil }},
		{"codecs", func() (interface{}, error) { return codecRepo.GetCodecs(), nil }},
		{"encoders", func() (interface{}, error) { return codecRepo.GetEncoders(), nil }},
		{"queue", func() (interface{}, error) { return encodeService.GetQueue(), nil }},
		{"series", func() (interface{}, error) { return seriesRepo.GetSeries() }},
		{"movies", func() (interface{}, error) { return movieRepo.GetMovies() }},
		{"history", func() (interface{}, error) { return historyRepo.GetHistories() }},
		{"logs", func() (interface{}, error) { return eventRepo.GetEvents() }},
	}

	for _, f := range fetchers {
		go func(f fetcher) {
			ticker := time.NewTicker(250 * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					data, err := f.getter()
					if err != nil {
						log.Printf("Error fetching %s: %v", f.dataType, err)
						continue
					}

					message := map[string]interface{}{
						f.dataType: data,
					}
					jsonData, err := json.Marshal(message)
					if err != nil {
						log.Printf("Error marshaling %s data: %v", f.dataType, err)
						continue
					}

					writeMutex.Lock()
					if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
						log.Printf("Error writing %s message: %v", f.dataType, err)
						writeMutex.Unlock()
						return
					}
					writeMutex.Unlock()
				case <-stopChan:
					return
				}
			}
		}(f)
	}
}
