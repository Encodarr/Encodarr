package controllers

import (
	"errors"
	"log"
	"net/http"
	"transfigurr/interfaces"
	"transfigurr/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	Repo interfaces.UserRepositoryInterface
}

func NewUserController(repo interfaces.UserRepositoryInterface) *UserController {
	return &UserController{
		Repo: repo,
	}
}

func (ctrl UserController) GetUsers(c *gin.Context) {
	users, err := ctrl.Repo.GetUser()
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
		} else {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving users"})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func (ctrl UserController) UpdateUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user, err := ctrl.Repo.GetUser()
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found."})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user.Username = username
	user.Password = string(passwordHash)
	if err := ctrl.Repo.UpdateUser(user); err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}
