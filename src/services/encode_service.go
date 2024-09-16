package services

import (
	"fmt"
	"sync"
	"transfigurr/interfaces"
	"transfigurr/models"
	"transfigurr/utils"
)

type EncodeService struct {
	encodeQueue  []models.Item
	encodeSet    map[string]struct{}
	mu           sync.Mutex
	cond         *sync.Cond
	eventService interfaces.EventServiceInterface
	seriesRepo   interfaces.SeriesRepositoryInterface
	seasonRepo   interfaces.SeasonRepositoryInterface
	episodeRepo  interfaces.EpisodeRepositoryInterface
	movieRepo    interfaces.MovieRepositoryInterface
	settingRepo  interfaces.SettingRepositoryInterface
	systemRepo   interfaces.SystemRepositoryInterface
	profileRepo  interfaces.ProfileRepositoryInterface
	authRepo     interfaces.AuthRepositoryInterface
	userRepo     interfaces.UserRepositoryInterface
	historyRepo  interfaces.HistoryRepositoryInterface
	eventRepo    interfaces.EventRepositoryInterface
	codecRepo    interfaces.CodecRepositoryInterface
}

func NewEncodeService(eventService interfaces.EventServiceInterface, seriesRepo interfaces.SeriesRepositoryInterface, seasonRepo interfaces.SeasonRepositoryInterface, episodeRepo interfaces.EpisodeRepositoryInterface, movieRepo interfaces.MovieRepositoryInterface, settingRepo interfaces.SettingRepositoryInterface, systemRepo interfaces.SystemRepositoryInterface, profileRepo interfaces.ProfileRepositoryInterface, authRepo interfaces.AuthRepositoryInterface, userRepo interfaces.UserRepositoryInterface, historyRepo interfaces.HistoryRepositoryInterface, eventRepo interfaces.EventRepositoryInterface, codecRepo interfaces.CodecRepositoryInterface) interfaces.EncodeServiceInterface {
	service := &EncodeService{
		encodeQueue:  make([]models.Item, 0),
		encodeSet:    make(map[string]struct{}),
		eventService: eventService,
		seriesRepo:   seriesRepo,
		seasonRepo:   seasonRepo,
		episodeRepo:  episodeRepo,
		movieRepo:    movieRepo,
		settingRepo:  settingRepo,
		systemRepo:   systemRepo,
		profileRepo:  profileRepo,
		authRepo:     authRepo,
		userRepo:     userRepo,
		historyRepo:  historyRepo,
		eventRepo:    eventRepo,
		codecRepo:    codecRepo,
	}
	service.cond = sync.NewCond(&service.mu)
	return service
}

func (s *EncodeService) Enqueue(item models.Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	itemID := fmt.Sprintf("%s_%s", item.Type, item.Id)
	if _, ok := s.encodeSet[itemID]; !ok {
		s.encodeSet[itemID] = struct{}{}
		s.encodeQueue = append(s.encodeQueue, item)
		s.cond.Signal()
	}
}

func (s *EncodeService) process() {
	for {
		s.mu.Lock()
		for len(s.encodeQueue) == 0 {
			s.cond.Wait()
		}
		item := s.encodeQueue[0]
		s.encodeQueue = s.encodeQueue[1:]
		s.mu.Unlock()
		s.processItem(item)
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

	queue := make([]models.Item, len(s.encodeQueue))
	copy(queue, s.encodeQueue)
	return queue
}

func (s *EncodeService) Startup() {
	go s.process()
}
