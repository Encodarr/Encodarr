package services

import (
	"log"
	"os"
	"time"
	"transfigurr/interfaces"
	"transfigurr/models"
)

type EventService struct {
	repo       interfaces.EventRepositoryInterface
	formatter  *log.Logger
	eventQueue chan models.Event
}

func NewEventService(eventRepo interfaces.EventRepositoryInterface, queueSize int) interfaces.EventServiceInterface {
	return &EventService{
		repo:       eventRepo,
		formatter:  log.New(os.Stdout, "", log.LstdFlags),
		eventQueue: make(chan models.Event, queueSize),
	}
}

func (s *EventService) Log(level, service, message string) {
	timestamp := time.Now().Format("2006-01-02T15:04:05.000")
	entry := models.Event{
		Timestamp: timestamp,
		Level:     level,
		Service:   service,
		Message:   message,
	}
	s.eventQueue <- entry
}

func (s *EventService) processLogQueue() {
	for {
		select {
		case entry, ok := <-s.eventQueue:
			if ok {
				err := s.repo.Log(entry.Level, entry.Service, entry.Message)
				if err != nil {
					s.formatter.Printf("Failed to log to event repository: %v", err)
				}
				s.formatter.Printf("%s [%s] %s: %s", entry.Timestamp, entry.Level, entry.Service, entry.Message)
			}
		case <-time.After(1 * time.Second):
			continue
		}
	}
}

func (s *EventService) Startup(logLevel string) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.SetFlags(0)
	if logLevel == "debug" {
		logger.SetFlags(log.LstdFlags | log.Lshortfile)
	}
	go s.processLogQueue()
}
