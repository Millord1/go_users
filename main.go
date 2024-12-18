package main

import (
	"fmt"
	"microservices/api/handler"
	"microservices/api/middlewares"
	"microservices/api/server"
	"net"
	"net/http"
	"net/rpc"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {

	var wg sync.WaitGroup
	wg.Add(2)

	go RPCServe(":8000")
	go HTTPServe(":8080")

	wg.Wait()
}

func RPCServe(port string) {
	// Run RPC server
	check := new(server.CheckUser)
	rpc.Register(check)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	fmt.Println("Service listening on 8000")
	rpc.Accept(listener)
}

func HTTPServe(port string) {
	// Run HTTP server
	r := gin.New()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, "Helo world")
	})

	r.GET("/users", handler.GetUsersNames)
	r.POST("/create/user", handler.NewUser)
	r.POST("/login", handler.UserLogin)

	authGroup := r.Group("/user", middlewares.AuthMiddleware())
	authGroup.GET("/test", handler.Test)
	authGroup.GET("/two_auth/activate", handler.Activate2Fa)
	/* authGroup.GET("two_auth/login", handler.EnterOTP) */

	r.Run(port)

}
