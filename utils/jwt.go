package utils

import (
	"log"
	"microservices/models"
	"net/http"
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

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, http.ErrAbortHandler
	}

	return token, nil
}

func GetUserData(tokenString string) (*models.User, error) {
	_, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	_, claimsErr := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if claimsErr != nil {
		return nil, err
	}

	return &models.User{
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
	}, nil
}
