package services

import (
	"fmt"
	"sync"
	"time"
	"transfigurr/interfaces"
	"transfigurr/models"
	"transfigurr/utils"
)

type EncodeService struct {
	encodeQueue chan models.Item
	encodeSet   map[string]struct{}
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

func NewEncodeService(seriesRepo interfaces.SeriesRepositoryInterface, seasonRepo interfaces.SeasonRepositoryInterface, episodeRepo interfaces.EpisodeRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface, authRepo interfaces.AuthRepositoryInterface, userRepo interfaces.UserRepositoryInterface, historyRepo interfaces.HistoryRepositoryInterface, eventRepo interfaces.EventRepositoryInterface, codecRepo interfaces.CodecRepositoryInterface) interfaces.EncodeServiceInterface {
	return &EncodeService{
		encodeQueue: make(chan models.Item, 100),
		encodeSet:   make(map[string]struct{}),
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

func (s *EncodeService) Enqueue(item models.Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	itemID := fmt.Sprintf("%s_%s", item.Type, item.Id)
	if _, ok := s.encodeSet[itemID]; !ok {
		s.encodeSet[itemID] = struct{}{}
		s.encodeQueue <- item
	}
}

func (s *EncodeService) process() {
	for {
		select {
		case item, ok := <-s.encodeQueue:
			if ok {
				s.processItem(item)
			}
		case <-time.After(1 * time.Second):
			continue
		}
	}
}

func (s *EncodeService) processItem(item models.Item) {
	if item.Type == "movie" {
		utils.EncodeMovie(item, s.movieRepo, s.settingRepo, s.profileRepo)
	}
	if item.Type == "episode" {
		utils.EncodeEpisode(item, s.seriesRepo, s.episodeRepo, s.settingRepo, s.profileRepo)
	}

	s.mu.Lock()
	delete(s.encodeSet, fmt.Sprintf("%s_%s", item.Type, item.Id))
	s.mu.Unlock()
}

func (s *EncodeService) GetQueue() []models.Item {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue := make([]models.Item, 0, len(s.encodeQueue))
	for len(s.encodeQueue) > 0 {
		item := <-s.encodeQueue
		queue = append(queue, item)
		s.encodeQueue <- item
	}
	return queue
}

func (s *EncodeService) Startup() {
	go s.process()
}
