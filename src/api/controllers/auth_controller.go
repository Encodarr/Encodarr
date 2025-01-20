package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"transfigurr/constants"
	"transfigurr/interfaces"
	"transfigurr/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	Repo interfaces.AuthRepositoryInterface
}

func NewAuthController(repo interfaces.AuthRepositoryInterface) *AuthController {
	return &AuthController{
		Repo: repo,
	}
}

func generateSecretKey() (string, error) {
	bytes := make([]byte, constants.SecretKeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (ctrl *AuthController) GetActivated(w http.ResponseWriter, r *http.Request) {
	user, err := ctrl.Repo.GetUser()
	response := map[string]bool{"activated": err == nil && user.Username != "" && user.Password != ""}
	json.NewEncoder(w).Encode(response)
}

func (ctrl *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var reqUser models.User
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	user, err := ctrl.Repo.GetUser()
	if err == nil && user.Username != "" && user.Password != "" {
		http.Error(w, "User already registered", http.StatusForbidden)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	user.Username = string(reqUser.Username)
	// Create the user
	if err := ctrl.Repo.UpdateUser(&user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "User registered successfully"}
	json.NewEncoder(w).Encode(response)
}

func (ctrl *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := ctrl.Repo.GetUser()
	if err != nil || user.Username != loginData.Username {
		log.Println("User not found or username mismatch")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		log.Printf("Login failed for user: %s", loginData.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": loginData.Username,
	})

	tokenString, err := token.SignedString([]byte(user.Secret))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": tokenString}
	json.NewEncoder(w).Encode(response)
}

func (ctrl *AuthController) LoginToken(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	user, err := ctrl.Repo.GetUser()
	if err != nil && user.Username == "" && user.Password == "" {
		http.Error(w, "You must register first", http.StatusForbidden)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(user.Secret), nil
	})

	if err != nil {
		http.Error(w, "Failed to parse token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		response := map[string]string{"message": "Welcome " + claims["username"].(string)}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	}
}
