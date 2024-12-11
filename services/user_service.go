package services

import (
	"errors"
	"microservices/models"
	"microservices/repository"
	"microservices/utils"
)

func CreateNewUser(repo repository.UserRepository, user *models.User) (*models.User, error) {

	userToPush, hashErr := user.CheckPasswordIsHashed()
	if hashErr != nil {
		return user, hashErr
	}
	return repo.Save(*userToPush)
}

func UpdateUser(repo repository.UserRepository, user *models.User) (*models.User, error) {
	toUpdate := models.User{
		Username: user.Username,
		Email:    user.Email,
	}

	return repo.Update(&toUpdate)
}

/* func UpdatePassword(repo repository.UserRepository, user *models.User, oldPw string) (*models.User, error) {
	if !user.VerifyPassword(oldPw) {
		err := errors.New("wrong password")
		log.Fatalln(err)
		return user, err
	}

	return repo.Update(user)
} */

func LoginUser(repo repository.UserRepository, email string, pw string) (string, error) {
	user, err := repo.FindByMail(email)
	if err != nil {
		return "", err
	}

	if !user.VerifyPassword(pw) {
		return "", errors.New("password doesn't match for user " + user.Username)
	}

	return utils.CreateToken(user)
}

func GetUserNames(repo repository.UserRepository) (*[]models.User, error) {
	return repo.FindAllNames()
}
