package seeds

import (
	"transfigurr/models"

	"gorm.io/gorm"
)

func SeedProfileSubtitleLanguages(db *gorm.DB) {
	if !db.Migrator().HasTable(&models.Seed{}) {
		db.Migrator().CreateTable(&models.Seed{})
		db.Migrator().CreateIndex(&models.Seed{}, "idx_name")
	}

	var seed models.Seed
	db.Where("name = ?", "SeedProfileSubtitleLanguages").First(&seed)
	if seed.Name == "SeedProfileSubtitleLanguages" {
		return
	}
	defaultProfileSubtitleLanguages := []models.ProfileSubtitleLanguage{
		{
			ProfileId: 0,
			Language:  "all",
		},
		{
			ProfileId: 1,
			Language:  "all",
		},
		{
			ProfileId: 2,
			Language:  "all",
		},
		{
			ProfileId: 3,
			Language:  "all",
		},
		{
			ProfileId: 4,
			Language:  "all",
		},
		{
			ProfileId: 5,
			Language:  "all",
		},
		{
			ProfileId: 6,
			Language:  "all",
		},
	}

	for _, defaultSubtitleLanguage := range defaultProfileSubtitleLanguages {
		db.Create(&defaultSubtitleLanguage)
	}
	db.Create(&models.Seed{Name: "SeedProfileSubtitleLanguages"})
}
