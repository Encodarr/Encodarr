package seeds

import (
	"crypto/rand"
	"encoding/hex"
	"transfigurr/internal/config"
	"transfigurr/internal/models"

	"gorm.io/gorm"
)

func generateSecretKey() (string, error) {
	bytes := make([]byte, config.SecretKeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func SeedUser(db *gorm.DB) {
	if !db.Migrator().HasTable(&models.Seed{}) {
		db.Migrator().CreateTable(&models.Seed{})
		db.Migrator().CreateIndex(&models.Seed{}, "idx_name")
	}

	var seed models.Seed
	db.Where("name = ?", "SeedUser").First(&seed)
	if seed.Name == "SeedUser" {
		return
	}

	secret, err := generateSecretKey()
	if err != nil {
		panic(err)
	}

	defaultSystems := []models.User{
		{
			Username: "",
			Password: "",
			Secret:   secret,
		},
	}

	for _, defaultSystem := range defaultSystems {
		db.Create(&defaultSystem)
	}

	db.Create(&models.Seed{Name: "SeedUser"})
}
