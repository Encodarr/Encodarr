package repositories

import "transfigurr/internal/models"

type AuthRepositoryInterface interface {
	GetUser() (models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
}
