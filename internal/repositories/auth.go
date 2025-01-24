package repositories

import (
	"errors"
	"log"
	"transfigurr/internal/models"

	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

func (repo *AuthRepository) GetUser() (models.User, error) {
	var user models.User

	if err := repo.DB.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, err
		}
		return models.User{}, err
	}
	return user, nil
}

func (repo *AuthRepository) CreateUser(user *models.User) error {
	if err := repo.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}
func (repo *AuthRepository) UpdateUser(user *models.User) error {
	var existingUser models.User

	result := repo.DB.Where("secret = ?", user.Secret).First(&existingUser)
	if result.Error != nil {
		log.Print("User not found:", result.Error)
		return result.Error
	}

	// Update the existing user with new values
	existingUser.Username = user.Username
	existingUser.Password = user.Password

	// Save the updated user
	// Update with WHERE condition
	return repo.DB.Model(&existingUser).
		Where("secret = ?", existingUser.Secret).
		Updates(map[string]interface{}{
			"username": user.Username,
			"password": user.Password,
		}).Error
}
