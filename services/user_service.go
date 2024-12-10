package services

import (
	"errors"
	"microservices/models"
	"microservices/repository"
	"microservices/utils"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{repo: repository.DbConnect(utils.GetEnvFile().Name)}
}

func (s *UserService) LoginUser(email string, pw string) (string, error) {
	user, err := s.repo.GetUserByMail(email)
	if err != nil {
		return "", err
	}

	if !user.VerifyPassword(pw) {
		return "", errors.New("password doesn't match for user " + user.Username)
	}

	return utils.CreateToken(user)
}

func (s *UserService) GetUserNames() (*[]models.User, error) {
	return s.repo.GetAllUserNames()
}

func (s *UserService) CreateNewUser(user *models.User) (*models.User, error) {
	return s.repo.CreateNewUser(user)
}
