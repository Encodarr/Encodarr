package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"transfigurr/constants"
	"transfigurr/models"
	"transfigurr/repository"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	Repo *repository.AuthRepository
}

func NewAuthController(repo *repository.AuthRepository) *AuthController {
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
	username := c.PostForm("username")
	password := c.PostForm("password")

	_, err := ctrl.Repo.GetUser()
	if err == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User already registered"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	secret, err := generateSecretKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating secret key"})
		return
	}
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
		Secret:   secret,
	}
	if err := ctrl.Repo.CreateUser(user); err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	user, err := ctrl.Repo.GetUser()
	if err != nil || user.Username != username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
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
