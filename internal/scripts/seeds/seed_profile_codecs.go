package seeds

import (
	"transfigurr/internal/models"

	"gorm.io/gorm"
)

func SeedProfileCodecs(db *gorm.DB) {
	if !db.Migrator().HasTable(&models.Seed{}) {
		db.Migrator().CreateTable(&models.Seed{})
		db.Migrator().CreateIndex(&models.Seed{}, "idx_name")
	}

	var seed models.Seed
	db.Where("name = ?", "SeedProfileCodecs").First(&seed)
	if seed.Name == "SeedProfileCodecs" {
		return
	}
	defaultProfileCodecs := []models.ProfileCodec{
		{
			ProfileId: 0,
			CodecId:   "Any",
		},
		{
			ProfileId: 1,
			CodecId:   "Any",
		},
		{
			ProfileId: 2,
			CodecId:   "Any",
		},
		{
			ProfileId: 3,
			CodecId:   "Any",
		},
		{
			ProfileId: 4,
			CodecId:   "Any",
		},
		{
			ProfileId: 5,
			CodecId:   "Any",
		},
		{
			ProfileId: 6,
			CodecId:   "Any",
		},
	}
	for _, defaultProfileCodec := range defaultProfileCodecs {
		db.Create(&defaultProfileCodec)
	}
	db.Create(&models.Seed{Name: "SeedProfileCodecs"})
}
