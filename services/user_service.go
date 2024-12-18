package services

import (
	"errors"
	"microservices/encryption"
	"microservices/models"
	"microservices/repository"

	"microservices/utils"

	"github.com/xlzd/gotp"
)

func CreateNewUser(repo repository.UserRepository, user *models.User) error {

	hashErr := user.CheckPasswordIsHashed()
	if hashErr != nil {
		return hashErr
	}
	err := repo.Save(*user)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(repo repository.UserRepository, user *models.User, pw string) error {
	if !user.VerifyPassword(pw) {
		return errors.New("Unauthorized")
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

func LoginUser(repo repository.UserRepository, email string, pw string, otp string) (string, error) {
	user, err := repo.FindByMail(email)
	if err != nil {
		return "", err
	}

	if !user.VerifyPassword(pw) {
		return "", errors.New("password doesn't match for user " + user.Username)
	}

	userTotp, err := encryption.DecryptData(user.Totp)
	if err != nil {
		return "", err
	}

	if !utils.VerifyOtp(string(userTotp), otp) {
		return "", errors.New("2FA Auth failed")
	}

	return utils.CreateToken(user)
}

func EnableTwoFactorAuth(repo repository.UserRepository, user *models.User) error {
	randSecret := gotp.RandomSecret(16)
	utils.GenerateTOTPWithSecret(user, randSecret)

	// Update user with encrypted TOTP
	err := repo.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func GetUserNames(repo repository.UserRepository) (*[]models.User, error) {
	return repo.FindAllNames()
}
