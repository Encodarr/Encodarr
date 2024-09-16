package routes

import (
	"encoding/json"
	"log"
	"net/http"
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
		return
	}
	defer conn.Close()

	// Set up ping/pong handlers
	conn.SetPingHandler(func(appData string) error {
		log.Println("Received ping")
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})

	conn.SetPongHandler(func(appData string) error {
		log.Println("Received pong")
		return nil
	})

	// Create a ticker that triggers every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				sendData(seriesRepo, movieRepo, profileRepo, settingRepo, systemRepo, historyRepo, eventRepo, codecRepo, encodeService, conn)
			}
		}
	}()

	// Read messages to keep the connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
	}
}

func sendData(seriesRepo interfaces.SeriesRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface, historyRepo interfaces.HistoryRepositoryInterface, eventRepo interfaces.EventRepositoryInterface, codecRepo interfaces.CodecRepositoryInterface, encodeService interfaces.EncodeServiceInterface, conn *websocket.Conn) {
	series, err := seriesRepo.GetSeries()
	if err != nil {
		log.Println("Error fetching series:", err)
	}

	settings, err := settingRepo.GetAllSettings()
	if err != nil {
		log.Println("Error fetching settings:", err)
	}

	movies, err := movieRepo.GetMovies()
	if err != nil {
		log.Println("Error fetching movies:", err)
	}

	profiles, err := profileRepo.GetAllProfiles()
	if err != nil {
		log.Println("Error fetching profiles:", err)
	}

	system, err := systemRepo.GetSystems()
	if err != nil {
		log.Println("Error fetching system:", err)
	}

	history, err := historyRepo.GetHistories()
	if err != nil {
		log.Println("Error fetching history:", err)
	}

	containers := codecRepo.GetContainers()
	if err != nil {
		log.Println("Error fetching containers:", err)
	}

	codecs := codecRepo.GetCodecs()

	encoders := codecRepo.GetEncoders()

	logs, err := eventRepo.GetEvents()
	if err != nil {
		log.Println("Error fetching logs:", err)
	}

	data := map[string]interface{}{
		"series":     series,
		"movies":     movies,
		"profiles":   profiles,
		"settings":   settings,
		"system":     system,
		"history":    history,
		"containers": containers,
		"codecs":     codecs,
		"encoders":   encoders,
		"logs":       logs,
		"queue":      encodeService.GetQueue(),
	}

	// Send the fetched data over the WebSocket connection
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling data:", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.Println("Error writing message:", err)
		return
	}
}
