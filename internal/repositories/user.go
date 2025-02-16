package repositories

import (
	"database/sql"
	"errors"
	"log"
	"transfigurr/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

var (
	ErrRecordNotFound = errors.New("record not found")
)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (repo *UserRepository) GetUser() (models.User, error) {
	var user models.User
	err := repo.DB.QueryRow(`
        SELECT id, username, password, secret 
        FROM users 
        LIMIT 1
    `).Scan(&user.Id, &user.Username, &user.Password, &user.Secret)

	if err == sql.ErrNoRows {
		log.Print(err)
		return models.User{}, ErrRecordNotFound
	}
	if err != nil {
		log.Print(err)
		return models.User{}, err
	}
	return user, nil
}

func (repo *UserRepository) UpsertUser(user models.User) error {
	var exists bool
	err := repo.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users LIMIT 1)").Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = repo.DB.Exec(`
            UPDATE users SET 
            username = ?, password = ?, secret = ?
            WHERE id = ?`,
			user.Username, user.Password, user.Secret, user.Id,
		)
	} else {
		_, err = repo.DB.Exec(`
            INSERT INTO users (username, password, secret)
            VALUES (?, ?, ?)`,
			user.Username, user.Password, user.Secret,
		)
	}
	return err
}

func (repo *UserRepository) DeleteUser() error {
	_, err := repo.DB.Exec("DELETE FROM users")
	return err
}
