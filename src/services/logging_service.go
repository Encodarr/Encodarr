package services

import (
	"log"
	"os"
	"time"
	"transfigurr/repository"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type LogEntry struct {
	Timestamp string
	Level     string
	Service   string
	Message   string
}

type LogService struct {
	repo      *repository.LoggingRepository
	formatter *log.Logger
	logQueue  chan LogEntry
}

func NewLogService(logRepo *repository.LoggingRepository, queueSize int) *LogService {
	return &LogService{
		repo:      logRepo,
		formatter: log.New(os.Stdout, "", log.LstdFlags),
		logQueue:  make(chan LogEntry, queueSize),
	}
}

func (s *LogService) Log(level, service, message string) {
	timestamp := time.Now().Format("2006-01-02T15:04:05.000")
	entry := LogEntry{
		Timestamp: timestamp,
		Level:     level,
		Service:   service,
		Message:   message,
	}
	s.logQueue <- entry
}

func (s *LogService) processLogQueue() {
	for entry := range s.logQueue {
		err := s.repo.Log(entry.Level, entry.Service, entry.Message)
		if err != nil {
			s.formatter.Printf("Failed to log to repository: %v", err)
		}
		s.formatter.Printf("%s [%s] %s: %s", entry.Timestamp, entry.Level, entry.Service, entry.Message)
	}
}

func (s *LogService) Startup(logLevel string) {
	logService := NewLogService(s.repo, 100)

	go logService.processLogQueue()

	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.SetFlags(0)
	if logLevel == "debug" {
		logger.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	// Example usage
	logService.Log("INFO", "example_service", "This is an info message")
	logService.Log("ERROR", "example_service", "This is an error message")

	// Allow some time for logs to be processed
	time.Sleep(2 * time.Second)
	close(logService.logQueue)
}
