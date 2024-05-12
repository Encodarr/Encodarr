package seeds

import (
	"transfigurr/models"

	"github.com/jinzhu/gorm"
)

func SeedProfileCodecs(db *gorm.DB) {
	if !db.HasTable(&models.Seed{}) {
		db.CreateTable(&models.Seed{})
		db.Model(&models.Seed{}).AddUniqueIndex("idx_name", "name")
	}

	var seed models.Seed
	db.Where("name = ?", "SeedProfileCodecs").First(&seed)
	if seed.Name == "SeedProfileCodecs" {
		return
	}
	defaultProfileCodecs := []models.ProfileCodec{
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
		{
			ProfileId: 7,
			CodecId:   "Any",
		},
	}
	for _, defaultProfileCodec := range defaultProfileCodecs {
		db.Create(&defaultProfileCodec)
	}
	db.Create(&models.Seed{Name: "SeedProfileCodecs"})
}
