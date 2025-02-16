package services

import (
	"fmt"
	"log"
	"sync"
	"transfigurr/internal/interfaces/repositories"
	"transfigurr/internal/interfaces/services"
	"transfigurr/internal/models"
	"transfigurr/internal/types"
	"transfigurr/internal/utils"
)

type MetadataService struct {
	metadataQueue []models.Item
	metadataSet   map[string]struct{}
	mu            sync.Mutex
	cond          *sync.Cond
	eventService  services.EventServiceInterface
	seriesRepo    repositories.SeriesRepositoryInterface
	seasonRepo    repositories.SeasonRepositoryInterface
	episodeRepo   repositories.EpisodeRepositoryInterface
	movieRepo     repositories.MovieRepositoryInterface
	settingRepo   repositories.SettingRepositoryInterface
	systemRepo    repositories.SystemRepositoryInterface
	profileRepo   repositories.ProfileRepositoryInterface
	authRepo      repositories.AuthRepositoryInterface
	userRepo      repositories.UserRepositoryInterface
	historyRepo   repositories.HistoryRepositoryInterface
	eventRepo     repositories.EventRepositoryInterface
	codecRepo     repositories.CodecRepositoryInterface
}

func NewMetadataService(eventService services.EventServiceInterface, repositories *types.Repositories) services.MetadataServiceInterface {
	service := &MetadataService{
		metadataQueue: make([]models.Item, 0),
		metadataSet:   make(map[string]struct{}),
		eventService:  eventService,
		seriesRepo:    repositories.SeriesRepo,
		seasonRepo:    repositories.SeasonRepo,
		episodeRepo:   repositories.EpisodeRepo,
		movieRepo:     repositories.MovieRepo,
		settingRepo:   repositories.SettingRepo,
		systemRepo:    repositories.SystemRepo,
		profileRepo:   repositories.ProfileRepo,
		authRepo:      repositories.AuthRepo,
		userRepo:      repositories.UserRepo,
		historyRepo:   repositories.HistoryRepo,
		eventRepo:     repositories.EventRepo,
		codecRepo:     repositories.CodecRepo,
	}
	service.cond = sync.NewCond(&service.mu)
	return service
}

func (s *MetadataService) EnqueueAll() {
	series, err := s.seriesRepo.GetSeries()
	if err != nil {
	}
	for _, series := range series {
		s.Enqueue(models.Item{Type: "series", Id: series.Id})
	}
	movies, err := s.movieRepo.GetMovies()
	if err != nil {
	}
	for _, m := range movies {
		s.Enqueue(models.Item{Type: "movie", Id: m.Id})
	}
}

func (s *MetadataService) Enqueue(item models.Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	itemID := fmt.Sprintf("%s_%s", item.Type, item.Id)
	if _, ok := s.metadataSet[itemID]; !ok {
		s.metadataSet[itemID] = struct{}{}
		s.metadataQueue = append(s.metadataQueue, item)
		s.cond.Signal()
	}
}

func (s *MetadataService) process() {
	for {
		s.mu.Lock()
		for len(s.metadataQueue) == 0 {
			s.cond.Wait()
		}
		item := s.metadataQueue[0]
		s.metadataQueue = s.metadataQueue[1:]
		s.mu.Unlock()
		s.processItem(item)
	}
}

func (s *MetadataService) processItem(item models.Item) {
	if item.Type == "movie" {
		movie, err := s.movieRepo.GetMovieById(item.Id)
		if err != nil {
		}
		movie, err = utils.GetMovieMetadata(movie)
		if err != nil {
		}
		s.movieRepo.UpsertMovie(movie.Id, movie)
	} else if item.Type == "series" {
		series, err := s.seriesRepo.GetSeriesByID(item.Id)
		if err != nil {
		}
		seriesMeta, err := utils.GetSeriesMetadata(series)
		if err != nil {
			log.Print(err)
		}

		series.Name = seriesMeta.Name
		series.Overview = seriesMeta.Overview
		series.ReleaseDate = seriesMeta.ReleaseDate
		series.Genre = seriesMeta.Genre
		series.Status = seriesMeta.Status
		series.Networks = seriesMeta.Networks
		series.Runtime = seriesMeta.Runtime
		s.seriesRepo.UpsertSeries(series.Id, series)
		for _, season := range series.Seasons {
			for index, episode := range season.Episodes {
				episode.EpisodeName = seriesMeta.Seasons[season.SeasonNumber-1].Episodes[index].EpisodeName
				episode.AirDate = seriesMeta.Seasons[season.SeasonNumber-1].Episodes[index].AirDate
				s.episodeRepo.UpsertEpisode(series.Id, season.SeasonNumber, episode.EpisodeNumber, episode)
			}
		}

	}

	s.mu.Lock()
	delete(s.metadataSet, fmt.Sprintf("%s_%s", item.Type, item.Id))
	s.mu.Unlock()
}

func (s *MetadataService) Startup() {
	go s.process()
}
