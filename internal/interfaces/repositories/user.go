package repositories

import "transfigurr/internal/models"

type UserRepositoryInterface interface {
	GetUser() (models.User, error)
}
