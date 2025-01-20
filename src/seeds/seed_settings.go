package seeds

import (
	"transfigurr/models"

	"gorm.io/gorm"
)

func SeedSettings(db *gorm.DB) {
	if !db.Migrator().HasTable(&models.Seed{}) {
		db.Migrator().CreateTable(&models.Seed{})
		db.Migrator().CreateIndex(&models.Seed{}, "idx_name")
	}

	var seed models.Seed
	db.Where("name = ?", "SeedSettings").First(&seed)
	if seed.Name == "SeedSettings" {
		return
	}
	defaultSettings := []models.Setting{
		{
			Id:    "theme",
			Value: "auto",
		},
		{
			Id:    "defaultProfile",
			Value: "1",
		},
		{
			Id:    "queueStatus",
			Value: "active",
		},
		{
			Id:    "queueStartupState",
			Value: "previous",
		},
		{
			Id:    "logLevel",
			Value: "info",
		},
		{
			Id:    "mediaView",
			Value: "posters",
		}, {
			Id:    "mediaSort",
			Value: "title",
		},
		{
			Id:    "massEditorSort",
			Value: "title",
		},
		{
			Id:    "massEditorSortDirection",
			Value: "ascending",
		},
		{
			Id:    "massEditorFilter",
			Value: "all",
		},
		{
			Id:    "mediaSortDirection",
			Value: "ascending",
		},
		{
			Id:    "mediaFilter",
			Value: "all",
		},
		{
			Id:    "TMDB",
			Value: "ZXlKaGJHY2lPaUpJVXpJMU5pSjkuZXlKaGRXUWlPaUprT1RCalpqQmhaREEyT0dJd01XVXpNVFkxTWpjNVltWXpPRE0xWmpRNU9TSXNJbk4xWWlJNklqWTFOR0UxWVRReE5qZGlOakV6TURFeFpqUXdaV0ZpWVNJc0luTmpiM0JsY3lJNld5SmhjR2xmY21WaFpDSmRMQ0oyWlhKemFXOXVJam94ZlEuNU1LVjViaXV0RmZvQkRuMk14aFMxQU1wbV9DTmE4QTh4WE5XTkFKUVNnTQ==",
		},
		{
			Id:    "mediaPosterPosterSize",
			Value: "medium",
		},
		{
			Id:    "mediaPosterDetailedProgressBar",
			Value: "false",
		},
		{
			Id:    "mediaPosterShowTitle",
			Value: "true",
		},
		{
			Id:    "mediaPosterShowMonitored",
			Value: "true",
		},
		{
			Id:    "mediaPosterShowProfile",
			Value: "true",
		},
		{
			Id:    "mediaTableShowNetwork",
			Value: "false",
		},
		{
			Id:    "mediaTableShowProfile",
			Value: "true",
		},
		{
			Id:    "mediaTableShowSeasons",
			Value: "true",
		},
		{
			Id:    "mediaTableShowEpisodes",
			Value: "true",
		},
		{
			Id:    "mediaTableShowEpisodeCount",
			Value: "false",
		},
		{
			Id:    "mediaTableShowYear",
			Value: "true",
		},
		{
			Id:    "mediaTableShowType",
			Value: "true",
		},
		{
			Id:    "mediaTableShowSizeOnDisk",
			Value: "true",
		},
		{
			Id:    "mediaTableShowSizeSaved",
			Value: "true",
		},
		{
			Id:    "mediaTableShowGenre",
			Value: "false",
		},
		{
			Id:    "mediaOverviewPosterSize",
			Value: "medium",
		},
		{
			Id:    "mediaOverviewDetailedProgressBar",
			Value: "false",
		},
		{
			Id:    "mediaOverviewShowMonitored",
			Value: "true",
		},
		{
			Id:    "mediaOverviewShowNetwork",
			Value: "true",
		},
		{
			Id:    "mediaOverviewShowProfile",
			Value: "true",
		},
		{
			Id:    "mediaOverviewShowSeasonCount",
			Value: "true",
		},
		{
			Id:    "mediaOverviewShowPath",
			Value: "false",
		},
		{
			Id:    "mediaOverviewShowSizeOnDisk",
			Value: "true",
		},
		{
			Id:    "queueFilter",
			Value: "all",
		},
		{
			Id:    "queuePageSize",
			Value: "12",
		},
		{
			Id:    "historyFilter",
			Value: "all",
		},
		{
			Id:    "historyPageSize",
			Value: "15",
		},
		{
			Id:    "eventsFilter",
			Value: "all",
		},
		{
			Id:    "eventsPageSize",
			Value: "15",
		},
		{
			Id:    "port",
			Value: "9889",
		},
	}
	for _, defaultSetting := range defaultSettings {
		db.Create(&defaultSetting)
	}
	db.Create(&models.Seed{Name: "SeedSettings"})

}
