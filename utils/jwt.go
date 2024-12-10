package utils

import (
	"errors"
	"log"
	"microservices/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

var secretKey string = getPassPhrase()

type JWTChecker interface {
	CreateToken(username string) (string, error)
	VerifyToken(tokenString string) error
}

func getPassPhrase() string {
	envFile := GetEnvFile().Name
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}

	return os.Getenv("PASSPHRASE")
}

func CreateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.Username,
			"email":    user.Email,
			"exp":      time.Now().Add(time.Hour * 4).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}
