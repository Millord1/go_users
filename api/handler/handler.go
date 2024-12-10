package handler

import (
	"microservices/models"
	"microservices/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsersNames(c *gin.Context) {
	us := services.NewUserService()
	allUserNames, err := us.GetUserNames()
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

	us := services.NewUserService()

	dbUser, err := us.CreateNewUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "Internal error", "message": "Internal error"})
	}

	c.JSON(http.StatusCreated, dbUser.Username)
}

func UserLogin(c *gin.Context) {
	email := c.PostForm("email")
	pw := c.PostForm("password")

	us := services.NewUserService()
	jwt, err := us.LoginUser(email, pw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "Bad request", "message": "Login failed"})
	}

	c.JSON(http.StatusAccepted, jwt)
}
