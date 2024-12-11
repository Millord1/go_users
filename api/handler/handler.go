package handler

import (
	"microservices/models"
	"microservices/repository"
	"microservices/services"
	"microservices/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var repo repository.UserRepository

func init() {
	repo = repository.DbConnect(utils.GetEnvFile().Name)
}

func GetUsersNames(c *gin.Context) {
	allUserNames, err := services.GetUserNames(repo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "Internal error", "message": "Internal error"})
	}
	c.JSON(http.StatusOK, allUserNames)
}

func NewUser(c *gin.Context) {
	user := models.User{
		Username: c.PostForm("username"),
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	dbUser, err := services.CreateNewUser(repo, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "Internal error", "message": "Internal error"})
	}

	c.JSON(http.StatusCreated, dbUser.Username)
}

func UserLogin(c *gin.Context) {
	email := c.PostForm("email")
	pw := c.PostForm("password")

	jwt, err := services.LoginUser(repo, email, pw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "Bad request", "message": "Login failed"})
	}

	c.JSON(http.StatusAccepted, jwt)
}
