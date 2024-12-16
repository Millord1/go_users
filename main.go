package main

import (
	"microservices/api/handler"
	"microservices/api/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()

	authGroup := r.Group("/user", middlewares.AuthMiddleware())

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, "Helo world")
	})

	r.GET("/users", handler.GetUsersNames)
	r.POST("/create/user", handler.NewUser)
	r.POST("/login", handler.UserLogin)

	authGroup.GET("/test", handler.Test)
	authGroup.GET("/two_auth/activate", handler.Activate2Fa)
	/* authGroup.GET("two_auth/login", handler.EnterOTP) */

	r.Run()

}
