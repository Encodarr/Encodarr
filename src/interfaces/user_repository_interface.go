package interfaces

import "transfigurr/models"

type UserRepositoryInterface interface {
	GetUser() (models.User, error)
}
