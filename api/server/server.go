package server

import (
	"microservices/repository"
	"microservices/services"
	"microservices/utils"
)

var repo repository.UserRepository

type CheckUser struct{}

type JWTVerify struct {
	TokenString string
}

type JWTReply struct {
	UserName string
	Email    string
}

func init() {
	repo = repository.DbConnect(utils.GetEnvFile().Name)
}

func (checkUser CheckUser) VerifyUserLogin(jwt *JWTVerify, reply *JWTReply) error {
	user, err := services.GetUserFromJWT(repo, jwt.TokenString)
	if err != nil {
		return err
	}

	reply.UserName = user.Username
	reply.Email = user.Email

	return nil
}
