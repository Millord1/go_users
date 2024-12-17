package handler

import (
	"fmt"
	"microservices/models"
	"microservices/repository"
	"microservices/services"
	"microservices/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var repo repository.UserRepository

type loginForm struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=5"`
	Otp      string `form:"otp"`
	/* Otp      string `form:"otp" validate:"required,numeric"` */
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

	err := services.CreateNewUser(repo, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "Internal error", "message": "Internal error"})
	}

	// Terminal output
	services.EnableTwoFactorAuth(repo, &user)
	/* c.JSON(http.StatusCreated, dbUser.Username) */
}

func UserLogin(c *gin.Context) {

	var login loginForm
	err := c.ShouldBind(&login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "fail"})
		return
	}

	jwt, err := services.LoginUser(repo, login.Email, login.Password, login.Otp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": err})
	}
	fmt.Println(jwt)

	c.JSON(http.StatusAccepted, jwt)
}

func Activate2Fa(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")
	user, err := services.GetUserFromJWT(repo, tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "Bad request"})
	}

	services.EnableTwoFactorAuth(repo, user)
}

/* func EnterOTP(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	user, err := services.GetUserFromJWT(repo, tokenString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "Bad request"})
	}
	utils.VerifyOtp(user.Totp)
} */

func Test(c *gin.Context) {
	/* 	data := "testTruc" */
	/* 	enc, _ := encryption.EncryptData(data)
	   	fmt.Println(enc)
	   	log.Fatalln(encryption.DecryptData(enc)) */
	/* c.JSON(http.StatusAccepted) */
}
