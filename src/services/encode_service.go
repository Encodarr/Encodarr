package services

import (
	"fmt"
	"sync"
	"transfigurr/interfaces"
	"transfigurr/models"
	"transfigurr/types"
	"transfigurr/utils"
)

type EncodeService struct {
	encodeQueue  []models.Item
	encodeSet    map[string]struct{}
	progress     *float64
	stage        *string
	eta          *int
	queueStatus  *string
	current      *models.Item
	mu           sync.Mutex
	cond         *sync.Cond
	eventService interfaces.EventServiceInterface
	repositories *types.Repositories
}

func NewEncodeService(eventService interfaces.EventServiceInterface, repositories *types.Repositories) interfaces.EncodeServiceInterface {
	service := &EncodeService{
		encodeQueue:  make([]models.Item, 0),
		encodeSet:    make(map[string]struct{}),
		progress:     new(float64),
		stage:        new(string),
		eta:          new(int),
		queueStatus:  new(string),
		current:      new(models.Item),
		eventService: eventService,
		repositories: repositories,
	}
	service.cond = sync.NewCond(&service.mu)

	queueStartupState, err := repositories.SettingRepo.GetSettingById("queueStartupState")
	if err != nil {
		queueStartupState = models.Setting{Value: "inactive"}
	}

	queueStatus, err := repositories.SettingRepo.GetSettingById("queueStatus")
	if err != nil {
		queueStatus = models.Setting{Value: "inactive"}
	}

	if queueStartupState.Value == "previous" {
		service.queueStatus = &queueStatus.Value
	} else {
		service.queueStatus = &queueStartupState.Value
	}

	idleStage := "Idle"
	service.stage = &idleStage

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
		queueStatus, err := s.repositories.SettingRepo.GetSettingById("queueStatus")
		if err != nil {
			queueStatus = models.Setting{Value: "inactive"}
		}
		if queueStatus.Value == "inactive" {
			continue
		}
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
		utils.EncodeMovie(item, s.repositories.MovieRepo, s.repositories.HistoryRepo, s.repositories.SettingRepo, s.repositories.ProfileRepo, s.stage, s.progress, s.eta, s.queueStatus, s.current)
	}
	if item.Type == "episode" {
		utils.EncodeEpisode(item, s.repositories.SeriesRepo, s.repositories.HistoryRepo, s.repositories.EpisodeRepo, s.repositories.SettingRepo, s.repositories.ProfileRepo, s.stage, s.progress, s.eta, s.queueStatus, s.current)
	}

	s.mu.Lock()
	delete(s.encodeSet, fmt.Sprintf("%s_%s", item.Type, item.Id))
	s.mu.Unlock()
}

func (s *EncodeService) GetQueue() models.QueueStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	queue := make([]models.Item, len(s.encodeQueue))
	copy(queue, s.encodeQueue)
	return models.QueueStatus{
		Queue:       queue,
		Progress:    *s.progress,
		Stage:       *s.stage,
		ETA:         *s.eta,
		Current:     *s.current,
		QueueStatus: *s.queueStatus,
	}

}

func (s *EncodeService) Startup() {
	go s.process()
}
