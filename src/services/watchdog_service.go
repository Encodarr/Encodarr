package services

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"
	"transfigurr/interfaces"
	"transfigurr/models"

	"github.com/fsnotify/fsnotify"
)

type WatchdogHandler struct {
	mediaType   string
	watcher     *fsnotify.Watcher
	scanService interfaces.ScanServiceInterface
}

func NewWatchdogService(scanService interfaces.ScanServiceInterface) *WatchdogHandler {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	return &WatchdogHandler{watcher: watcher, scanService: scanService}
}

func (w *WatchdogHandler) OnCreated(path string) {
	log.Println("Watchdog detected a file creation")
	if !isDirectory(path) {
		w.WaitUntilDone(path)
	}
	w.HandleChange(path)
}

func (w *WatchdogHandler) OnDeleted(path string) {
	log.Println("Watchdog detected a file deletion", path)
	if !isDirectory(path) {
		log.Println("Path is not a directory, waiting until done:", path)
		w.WaitUntilDone(path)
	} else {
		log.Println("Path is a directory, skipping WaitUntilDone:", path)
	}
	log.Println("Handling change for path:", path)
	w.HandleChange(path)
}

func (w *WatchdogHandler) OnModified(path string) {
	log.Println("Watchdog detected a file modification")
	if !isDirectory(path) {
		w.WaitUntilDone(path)
	}
	w.HandleChange(path)
}

func (w *WatchdogHandler) WaitUntilDone(path string) {
	oldFileSize := int64(-1)
	for {
		newFileSize, err := getFileSize(path)
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
	log.Print("handle change for", path)
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
		log.Print("media", media)
		if media == "" {
			log.Print("Enqueuing all movies")
			w.scanService.EnqueueAllMovies()
		} else {
			log.Print("Enqueuing movie", media)
			w.scanService.Enqueue(models.Item{Id: media, Type: "movie"})
		}
	}
}

func (w *WatchdogHandler) watchDirectory(directory string) {
	err := w.watcher.Add(directory)
	if err != nil {
		log.Println("Error adding directory to watcher:", err)
		return
	}
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Println("Error reading directory:", err)
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
	w.mediaType = contentType
	w.watchDirectory(directory)

	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				log.Printf("Event detected: %s on %s", event.Op, event.Name)
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
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	log.Print("starting watchdog for", directory)
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

func (w *WatchdogHandler) Startup(directory, contentType string) {
	go w.process(directory, contentType)
}
