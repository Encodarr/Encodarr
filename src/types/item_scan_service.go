package models

import (
	"sync"
	"transfigurr/interfaces"
	"transfigurr/models"
)

type ItemScanService struct {
	ScanQueue   chan models.Item
	ScanSet     map[string]struct{}
	Mu          sync.Mutex
	SeriesRepo  *interfaces.SeriesRepositoryInterface
	SeasonRepo  *interfaces.SeasonRepositoryInterface
	EpisodeRepo *interfaces.EpisodeRepositoryInterface
	MovieRepo   *interfaces.MovieRepositoryInterface
	SettingRepo *interfaces.SettingRepositoryInterface
	SystemRepo  *interfaces.SystemRepositoryInterface
	ProfileRepo *interfaces.ProfileRepositoryInterface
	AuthRepo    *interfaces.AuthRepositoryInterface
	UserRepo    *interfaces.UserRepositoryInterface
	HistoryRepo *interfaces.HistoryRepositoryInterface
	EventRepo   *interfaces.EventRepositoryInterface
	CodecRepo   *interfaces.CodecRepositoryInterface
}
