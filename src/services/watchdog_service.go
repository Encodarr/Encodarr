package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/fsnotify/fsnotify"
)

type WatchdogHandler struct {
	mediaType string
	watcher   *fsnotify.Watcher
}

func NewWatchdogService(queueSize int) *WatchdogHandler {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	return &WatchdogHandler{watcher: watcher}
}

func (w *WatchdogHandler) OnCreated(path string) {
	log.Println("Watchdog detected a file creation")
	w.WaitUntilDone(path)
	w.HandleChange(path)
}

func (w *WatchdogHandler) OnDeleted(path string) {
	log.Println("Watchdog detected a file deletion")
	w.HandleChange(path)
}

func (w *WatchdogHandler) OnModified(path string) {
	log.Println("Watchdog detected a file modification")
	w.WaitUntilDone(path)
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
	var media string
	if w.mediaType == "series" {
		media = GetSeriesName(path)
		if media != "" {
			fmt.Printf("Enqueue series: %s\n", media)
		} else {
			fmt.Println("Enqueue all series")
		}
	} else {
		media = GetMovieName(path)
		if media != "" {
			fmt.Printf("Enqueue movie: %s\n", media)
		} else {
			fmt.Println("Enqueue all movies")
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

func (w *WatchdogHandler) Startup(directory, contentType string) {
	go w.process(directory, contentType)
}
