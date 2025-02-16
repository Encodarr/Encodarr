package repositories

import (
	"database/sql"
	"errors"
	"log"
	"transfigurr/internal/models"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

func (repo *AuthRepository) GetUser() (models.User, error) {
	var user models.User
	err := repo.DB.QueryRow(`
        SELECT id, username, password, secret 
        FROM users 
        LIMIT 1
    `).Scan(&user.Id, &user.Username, &user.Password, &user.Secret)

	if err == sql.ErrNoRows {
		return models.User{}, errors.New("record not found")
	}
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (repo *AuthRepository) CreateUser(user *models.User) error {
	_, err := repo.DB.Exec(`
        INSERT INTO users (username, password, secret)
        VALUES (?, ?, ?)`,
		user.Username, user.Password, user.Secret,
	)
	return err
}

func (repo *AuthRepository) UpdateUser(user *models.User) error {
	result, err := repo.DB.Exec(`
        UPDATE users 
        SET username = ?, password = ?
        WHERE secret = ?`,
		user.Username, user.Password, user.Secret,
	)
	if err != nil {
		log.Print("Error updating user:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
