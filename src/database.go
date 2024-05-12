package main

import (
	"transfigurr/models"
	"transfigurr/seeds"

	"github.com/jinzhu/gorm"
)

func InitDB(db *gorm.DB) {
	db.AutoMigrate(&models.Episode{})
	db.AutoMigrate(&models.History{})
	db.AutoMigrate(&models.Event{})
	db.AutoMigrate(&models.Movie{})
	db.AutoMigrate(&models.Profile{})
	db.AutoMigrate(&models.ProfileAudioLanguage{})
	db.AutoMigrate(&models.ProfileCodec{})
	db.AutoMigrate(&models.ProfileSubtitleLanguage{})
	db.AutoMigrate(&models.Season{})
	db.AutoMigrate(&models.Series{})
	db.AutoMigrate(&models.Setting{})
	db.AutoMigrate(&models.System{})
	db.AutoMigrate(&models.User{})
}

func SeedDB(db *gorm.DB) {
	seeds.SeedProfileAudioLanguages(db)
	seeds.SeedProfileSubtitleLanguages(db)
	seeds.SeedProfileCodecs(db)
	seeds.SeedProfiles(db)
	seeds.SeedSettings(db)
	seeds.SeedSystems(db)
}
