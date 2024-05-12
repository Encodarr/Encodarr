package seeds

import (
	"transfigurr/models"

	"github.com/jinzhu/gorm"
)

func SeedProfileAudioLanguages(db *gorm.DB) {
	if !db.HasTable(&models.Seed{}) {
		db.CreateTable(&models.Seed{})
		db.Model(&models.Seed{}).AddUniqueIndex("idx_name", "name")
	}

	var seed models.Seed
	db.Where("name = ?", "SeedProfileAudioLanguages").First(&seed)
	if seed.Name == "SeedProfileAudioLanguages" {
		return
	}
	defaultProfileAudioLanguages := []models.ProfileAudioLanguage{
		{
			ProfileId: "1",
			Language:  "all",
		},
		{
			ProfileId: "2",
			Language:  "all",
		},
		{
			ProfileId: "3",
			Language:  "all",
		},
		{
			ProfileId: "4",
			Language:  "all",
		},
		{
			ProfileId: "5",
			Language:  "all",
		},
		{
			ProfileId: "6",
			Language:  "all",
		},
		{
			ProfileId: "7",
			Language:  "all",
		},
	}

	for _, defaultAudioLanguage := range defaultProfileAudioLanguages {
		db.Create(&defaultAudioLanguage)
	}
	db.Create(&models.Seed{Name: "SeedProfileAudioLanguages"})
}
