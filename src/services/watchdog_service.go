package services

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
	"transfigurr/interfaces"
	"transfigurr/models"
)

type WatchdogHandler struct {
	mediaType    string
	scanService  interfaces.ScanServiceInterface
	pollInterval time.Duration
	fileStates   map[string]FileState
	mutex        sync.RWMutex
	stopChan     chan struct{}
}

type FileState struct {
	ModTime time.Time
	Size    int64
	IsDir   bool
}

func NewWatchdogService(scanService interfaces.ScanServiceInterface) *WatchdogHandler {
	return &WatchdogHandler{
		scanService:  scanService,
		pollInterval: 30 * time.Second,
		fileStates:   make(map[string]FileState),
		stopChan:     make(chan struct{}),
	}
}

func (w *WatchdogHandler) Startup(directory, contentType string) {
	w.mediaType = contentType
	go w.watchDirectory(directory)
}

func (w *WatchdogHandler) watchDirectory(directory string) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	// Initial scan
	w.scanDirectory(directory)

	for {
		select {
		case <-ticker.C:
			w.scanDirectory(directory)
		case <-w.stopChan:
			return
		}
	}
}

func (w *WatchdogHandler) scanDirectory(directory string) {
	currentStates := make(map[string]FileState)

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("[ERROR] Error accessing path %s: %v", path, err)
			return filepath.SkipDir
		}

		currentStates[path] = FileState{
			ModTime: info.ModTime(),
			Size:    info.Size(),
			IsDir:   info.IsDir(),
		}

		w.mutex.RLock()
		oldState, exists := w.fileStates[path]
		w.mutex.RUnlock()

		if !exists {
			w.OnCreated(path)
		} else if oldState.ModTime != info.ModTime() || oldState.Size != info.Size() {
			w.OnModified(path)
		}

		return nil
	})

	if err != nil {
		log.Printf("[ERROR] Directory scan error: %v", err)
		return
	}

	// Check for deletions
	w.mutex.RLock()
	for path := range w.fileStates {
		if _, exists := currentStates[path]; !exists {
			w.OnDeleted(path)
		}
	}
	w.mutex.RUnlock()

	// Update states
	w.mutex.Lock()
	w.fileStates = currentStates
	w.mutex.Unlock()
}

func (w *WatchdogHandler) OnCreated(path string) {
	log.Printf("File created: %s", path)
	if !isDirectory(path) {
		w.WaitUntilDone(path)
	}
	w.HandleChange(path)
}

func (w *WatchdogHandler) OnDeleted(path string) {
	log.Printf("File deleted: %s", path)
	w.HandleChange(path)
}

func (w *WatchdogHandler) OnModified(path string) {
	log.Printf("File modified: %s", path)
	if !isDirectory(path) {
		w.WaitUntilDone(path)
	}
	w.HandleChange(path)
}

func (w *WatchdogHandler) WaitUntilDone(path string) {
	oldFileSize := int64(-1)
	for {
		newFileSize, err := getFileSize(path)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		if newFileSize == oldFileSize {
			break
		}
		oldFileSize = newFileSize
		time.Sleep(5 * time.Second)
	}
}

func (w *WatchdogHandler) HandleChange(path string) {
	var media string
	if w.mediaType == "series" {
		media = GetSeriesName(path)
		if media == "" {
			w.scanService.EnqueueAllSeries()
		} else {
			w.scanService.Enqueue(models.Item{Id: media, Type: "series"})
		}
	} else {
		media = GetMovieName(path)
		if media == "" {
			w.scanService.EnqueueAllMovies()
		} else {
			w.scanService.Enqueue(models.Item{Id: media, Type: "movie"})
		}
	}
}

func GetSeriesName(path string) string {
	re := regexp.MustCompile(`/series/([^/]*)`)
	match := re.FindStringSubmatch(path)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func GetMovieName(path string) string {
	re := regexp.MustCompile(`/movies/([^/]*)`)
	match := re.FindStringSubmatch(path)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func getFileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func (w *WatchdogHandler) Stop() {
	close(w.stopChan)
}
