package services

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
	"transfigurr/interfaces"
	"transfigurr/models"

	"github.com/fsnotify/fsnotify"
)

type WatchdogHandler struct {
	mediaType    string
	watcher      *fsnotify.Watcher
	scanService  interfaces.ScanServiceInterface
	pollInterval time.Duration
	dirCache     map[string]time.Time
	mutex        sync.RWMutex
}

func NewWatchdogService(scanService interfaces.ScanServiceInterface) *WatchdogHandler {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil
	}
	return &WatchdogHandler{
		watcher:      watcher,
		scanService:  scanService,
		pollInterval: 30 * time.Second,
		dirCache:     make(map[string]time.Time),
	}
}

func (w *WatchdogHandler) OnCreated(path string) {
	if !isDirectory(path) {
		w.WaitUntilDone(path)
	}
	w.HandleChange(path)
}

func (w *WatchdogHandler) OnDeleted(path string) {
	if !isDirectory(path) {
		w.WaitUntilDone(path)
	} else {
	}
	w.HandleChange(path)
}

func (w *WatchdogHandler) OnModified(path string) {
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
		} else {
			oldFileSize = newFileSize
			time.Sleep(5 * time.Second)
		}
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

func (w *WatchdogHandler) watchDirectory(directory string) {
	err := w.watcher.Add(directory)
	if err != nil {
		return
	}
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return
	}
	for _, file := range files {
		path := directory + "/" + file.Name()
		if file.IsDir() {
			w.watchDirectory(path)
		}
	}
}

func (w *WatchdogHandler) process(directory, contentType string) {
	log.Printf("Starting watchdog for %s directory: %s", contentType, directory)
	w.mediaType = contentType
	w.watchDirectory(directory)

	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					log.Println("[ERROR] Event channel closed")
					return
				}
				log.Printf("Received event: %s for %s", event.Op, event.Name)

				switch event.Op {
				case fsnotify.Create:
					w.OnCreated(event.Name)
				case fsnotify.Remove:
					w.OnDeleted(event.Name)
				case fsnotify.Write:
					w.OnModified(event.Name)
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					log.Println("[ERROR] Error channel closed")
					return
				}
				log.Printf("[ERROR] Watcher error: %v", err)
			}
		}
	}()

	select {}
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

func (w *WatchdogHandler) startPolling(directory string) {
	ticker := time.NewTicker(w.pollInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.pollDirectory(directory)
			}
		}
	}()
}

func (w *WatchdogHandler) pollDirectory(directory string) {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		w.mutex.Lock()
		lastMod, exists := w.dirCache[path]
		currentMod := info.ModTime()
		w.dirCache[path] = currentMod
		w.mutex.Unlock()

		if !exists {
			w.OnCreated(path)
		} else if currentMod.After(lastMod) {
			w.OnModified(path)
		}

		return nil
	})

	if err != nil {
		log.Printf("[ERROR] Polling error: %v", err)
	}
}

func (w *WatchdogHandler) Startup(directory, contentType string) {
	w.startPolling(directory)            // Start polling
	go w.process(directory, contentType) // Keep existing fsnotify
}
