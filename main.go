package main

import (
	"microservices/api/handler"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.GET("/users", handler.GetUsersNames)
	r.POST("/user/create", handler.NewUser)
	r.POST("user/login", handler.UserLogin)

	r.Run()

}
