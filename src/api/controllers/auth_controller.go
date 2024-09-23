package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"transfigurr/constants"
	"transfigurr/interfaces"
	"transfigurr/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

func (ctrl *AuthController) GetActivated(c *gin.Context) {
	_, err := ctrl.Repo.GetUser()
	if err != nil {
		c.JSON(200, gin.H{"activated": false})
		return
	}
	c.JSON(200, gin.H{"activated": true})
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	log.Printf("Username: %s, Password: %s", user.Username, user.Password)

	// Check if user already exists
	_, err := ctrl.Repo.GetUser()
	if err == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User already registered"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Create the user
	if err := ctrl.Repo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	log.Printf("Username: %s, Password: %s", loginData.Username, loginData.Password)

	user, err := ctrl.Repo.GetUser()
	if err != nil || user.Username != loginData.Username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": loginData.Username,
	})
	tokenString, err := token.SignedString([]byte(user.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (ctrl *AuthController) LoginToken(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	user, err := ctrl.Repo.GetUser()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You must register first"})
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(user.Secret), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token: " + err.Error()})
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome " + claims["username"].(string)})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	}
}
