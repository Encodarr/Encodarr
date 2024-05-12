package services

import (
	"fmt"
	"log"
	"sync"
	"time"
	"transfigurr/interfaces"
	"transfigurr/models"
	"transfigurr/tasks"
)

type ItemScanService struct {
	scanQueue   chan models.Item
	scanSet     map[string]struct{}
	mu          sync.Mutex
	seriesRepo  interfaces.SeriesRepositoryInterface
	seasonRepo  interfaces.SeasonRepositoryInterface
	episodeRepo interfaces.EpisodeRepositoryInterface
	movieRepo   interfaces.MovieRepositoryInterface
	settingRepo interfaces.SettingRepositoryInterface
	systemRepo  interfaces.SystemRepositoryInterface
	profileRepo interfaces.ProfileRepositoryInterface
	authRepo    interfaces.AuthRepositoryInterface
	userRepo    interfaces.UserRepositoryInterface
	historyRepo interfaces.HistoryRepositoryInterface
	eventRepo   interfaces.EventRepositoryInterface
	codecRepo   interfaces.CodecRepositoryInterface
}

func NewItemScanService(seriesRepo interfaces.SeriesRepositoryInterface, seasonRepo interfaces.SeasonRepositoryInterface, episodeRepo interfaces.EpisodeRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface, authRepo interfaces.AuthRepositoryInterface, userRepo interfaces.UserRepositoryInterface, historyRepo interfaces.HistoryRepositoryInterface, eventRepo interfaces.EventRepositoryInterface, codecRepo interfaces.CodecRepositoryInterface) interfaces.ScanServiceInterface {
	return &ItemScanService{
		scanQueue:   make(chan models.Item),
		scanSet:     make(map[string]struct{}),
		seriesRepo:  seriesRepo,
		seasonRepo:  seasonRepo,
		episodeRepo: episodeRepo,
		movieRepo:   movieRepo,
		settingRepo: settingRepo,
		systemRepo:  systemRepo,
		profileRepo: profileRepo,
		authRepo:    authRepo,
		userRepo:    userRepo,
		historyRepo: historyRepo,
		eventRepo:   eventRepo,
		codecRepo:   codecRepo,
	}
}

func (s *ItemScanService) Enqueue(item models.Item) {
	s.mu.Lock()
	defer s.mu.Unlock()

	itemID := fmt.Sprintf("%s_%s", item.Type, item.Id)
	if _, ok := s.scanSet[itemID]; !ok {
		log.Print("Enqueuing item")
		s.scanSet[itemID] = struct{}{}
		s.scanQueue <- item
	}
}

func (s *ItemScanService) process() {
	for {
		log.Print("Processing")
		select {
		case item, ok := <-s.scanQueue:
			if ok {
				log.Print("Processing item")
				go s.processItem(item)
			}
		case <-time.After(1 * time.Second):
			continue
		}
	}
}

func (s *ItemScanService) processItem(item models.Item) {
	fmt.Printf("Processing item: %s of type: %s\n", item.Id, item.Type)
	if item.Type == "movie" {
		log.Print("Processing movie")
		tasks.ScanMovie(item.Id, s.movieRepo, s.settingRepo, s.profileRepo)
	}

	s.mu.Lock()
	delete(s.scanSet, fmt.Sprintf("%s_%s", item.Type, item.Id))
	s.mu.Unlock()
}

func (s *ItemScanService) Startup() {
	go s.process()
}
