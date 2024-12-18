package server

import (
	"microservices/repository"
	"microservices/services"
	"microservices/utils"
)

// init user repo
var repo repository.UserRepository = repository.DbConnect(utils.GetEnvFile().Name)

type CheckUser struct{}

type JWTVerify struct {
	TokenString string
}

type JWTReply struct {
	ID       uint
	UserName string
	Email    string
}

/* func init() {
	repo = repository.DbConnect(utils.GetEnvFile().Name)
} */

func (checkUser CheckUser) VerifyUserLogin(jwt *JWTVerify, reply *JWTReply) error {
	// RPC response to identify user from JWT
	user, err := services.GetUserFromJWT(repo, jwt.TokenString)
	if err != nil {
		return err
	}

	reply.UserName = user.Username
	reply.Email = user.Email
	reply.ID = user.ID

	return nil
}
