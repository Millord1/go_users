package services

import (
	"errors"
	"microservices/models"
	"microservices/repository"
	"microservices/utils"

	"github.com/xlzd/gotp"
)

func CreateNewUser(repo repository.UserRepository, user *models.User) (*models.User, error) {

	userToPush, hashErr := user.CheckPasswordIsHashed()
	if hashErr != nil {
		return user, hashErr
	}
	return repo.Save(*userToPush)
}

func UpdateUser(repo repository.UserRepository, user *models.User, pw string) (*models.User, error) {
	if !user.VerifyPassword(pw) {
		return user, errors.New("Unauthorized")
	}

	toUpdate := models.User{
		Username: user.Username,
		Email:    user.Email,
	}

	return repo.Update(&toUpdate)
}

func GetUserFromMail(repo repository.UserRepository, email string) (*models.User, error) {
	user, err := repo.FindByMail(email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetUserFromJWT(repo repository.UserRepository, tokenString string) (*models.User, error) {
	expectedUser, err := utils.GetUserData(tokenString)
	if err != nil {
		return nil, err
	}

	user, userErr := repo.FindByMail(expectedUser.Email)
	if userErr != nil || user.Username != expectedUser.Username {
		return nil, err
	}

	return user, nil
}

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

func EnableTwoFactorAuth(repo repository.UserRepository, user *models.User) (*models.User, error) {
	randSecret := gotp.RandomSecret(16)
	utils.GenerateTOTPWithSecret(user, randSecret)

	// Update user with encrypted TOTP
	_, err := repo.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserNames(repo repository.UserRepository) (*[]models.User, error) {
	return repo.FindAllNames()
}
