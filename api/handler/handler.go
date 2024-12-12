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

type LoginForm struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=5"`
}

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

	var login LoginForm
	err := c.ShouldBind(&login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "fail"})
		return
	}

	jwt, err := services.LoginUser(repo, login.Email, login.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": err})
	}

	c.JSON(http.StatusAccepted, jwt)

	/*
		 	jwt, err := services.LoginUser(repo, email, pw)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": "Bad request", "message": "Login failed"})
			}

			c.JSON(http.StatusAccepted, jwt)
	*/
}

func Test(c *gin.Context) {
	c.JSON(http.StatusAccepted, "Hello, auth world")
}
