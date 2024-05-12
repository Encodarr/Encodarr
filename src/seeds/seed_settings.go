package seeds

import (
	"transfigurr/models"

	"github.com/jinzhu/gorm"
)

func SeedSettings(db *gorm.DB) {
	if !db.HasTable(&models.Seed{}) {
		db.CreateTable(&models.Seed{})
		db.Model(&models.Seed{}).AddUniqueIndex("Idx_name", "name")
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
			Id:    "default_profile",
			Value: "1",
		},
		{
			Id:    "queue_status",
			Value: "active",
		},
		{
			Id:    "queue_startup_state",
			Value: "previous",
		},
		{
			Id:    "log_level",
			Value: "info",
		},
		{
			Id:    "media_view",
			Value: "posters",
		}, {
			Id:    "media_sort",
			Value: "title",
		},
		{
			Id:    "massEditor_sort",
			Value: "title",
		},
		{
			Id:    "massEditor_sort_direction",
			Value: "ascending",
		},
		{
			Id:    "massEditor_filter",
			Value: "all",
		},
		{
			Id:    "media_sort_direction",
			Value: "ascending",
		},
		{
			Id:    "media_filter",
			Value: "all",
		},
		{
			Id:    "TMDB",
			Value: "ZXlKaGJHY2lPaUpJVXpJMU5pSjkuZXlKaGRXUWlPaUprT1RCalpqQmhaREEyT0dJd01XVXpNVFkxTWpjNVltWXpPRE0xWmpRNU9TSXNJbk4xWWlJNklqWTFOR0UxWVRReE5qZGlOakV6TURFeFpqUXdaV0ZpWVNJc0luTmpiM0JsY3lJNld5SmhjR2xmY21WaFpDSmRMQ0oyWlhKemFXOXVJam94ZlEuNU1LVjViaXV0RmZvQkRuMk14aFMxQU1wbV9DTmE4QTh4WE5XTkFKUVNnTQ==",
		},
		{
			Id:    "media_poster_posterSize",
			Value: "medium",
		},
		{
			Id:    "media_poster_detailedProgressBar",
			Value: "false",
		},
		{
			Id:    "media_poster_showTitle",
			Value: "true",
		},
		{
			Id:    "media_poster_showMonitored",
			Value: "true",
		},
		{
			Id:    "media_poster_showProfile",
			Value: "true",
		},
		{
			Id:    "media_table_showNetwork",
			Value: "false",
		},
		{
			Id:    "media_table_showProfile",
			Value: "true",
		},
		{
			Id:    "media_table_showSeasons",
			Value: "true",
		},
		{
			Id:    "media_table_showEpisodes",
			Value: "true",
		},
		{
			Id:    "media_table_showEpisodeCount",
			Value: "false",
		},
		{
			Id:    "media_table_showYear",
			Value: "true",
		},
		{
			Id:    "media_table_showType",
			Value: "true",
		},
		{
			Id:    "media_table_showSizeOnDisk",
			Value: "true",
		},
		{
			Id:    "media_table_showSizeSaved",
			Value: "true",
		},
		{
			Id:    "media_table_showGenre",
			Value: "false",
		},
		{
			Id:    "media_overview_posterSize",
			Value: "medium",
		},
		{
			Id:    "media_overview_detailedProgressBar",
			Value: "false",
		},
		{
			Id:    "media_overview_showMonitored",
			Value: "true",
		},
		{
			Id:    "media_overview_showNetwork",
			Value: "true",
		},
		{
			Id:    "media_overview_showProfile",
			Value: "true",
		},
		{
			Id:    "media_overview_showSeasonCount",
			Value: "true",
		},
		{
			Id:    "media_overview_showPath",
			Value: "false",
		},
		{
			Id:    "media_overview_showSizeOnDisk",
			Value: "true",
		},
		{
			Id:    "queue_filter",
			Value: "all",
		},
		{
			Id:    "queue_page_size",
			Value: "12",
		},
		{
			Id:    "history_filter",
			Value: "all",
		},
		{
			Id:    "history_page_size",
			Value: "15",
		},
		{
			Id:    "events_filter",
			Value: "all",
		},
		{
			Id:    "events_page_size",
			Value: "15",
		},
		{
			Id:    "port",
			Value: "7889",
		},
	}
	for _, defaultSetting := range defaultSettings {
		db.Create(&defaultSetting)
	}
	db.Create(&models.Seed{Name: "SeedSettings"})

}
