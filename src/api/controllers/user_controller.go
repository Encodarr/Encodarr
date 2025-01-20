package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"transfigurr/interfaces"
	"transfigurr/repository"
)

type UserController struct {
	Repo interfaces.UserRepositoryInterface
}

func NewUserController(repo interfaces.UserRepositoryInterface) *UserController {
	return &UserController{
		Repo: repo,
	}
}

func (ctrl UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := ctrl.Repo.GetUser()
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			http.Error(w, "Users not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// func (ctrl UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	if err := r.ParseForm(); err != nil {
// 		http.Error(w, "Error parsing form data", http.StatusBadRequest)
// 		return
// 	}

// 	username := r.FormValue("username")
// 	password := r.FormValue("password")

// 	user, err := ctrl.Repo.GetUser()
// 	if err != nil {
// 		http.Error(w, "User not found", http.StatusForbidden)
// 		return
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		http.Error(w, "Error hashing password", http.StatusInternalServerError)
// 		return
// 	}

// 	user.Username = username
// 	user.Password = string(hashedPassword)

// 	if err := ctrl.Repo.UpdateUser(user); err != nil {
// 		http.Error(w, "Error updating user", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(user)
// }
