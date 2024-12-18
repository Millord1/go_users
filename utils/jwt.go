package utils

import (
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
	// Used to encode JWT
	envFile := GetEnvFile().Name
	if err := godotenv.Load(envFile); err != nil {
		logger.Sugar.Fatal("Error loading " + envFile + " file")
	}

	return os.Getenv("PASSPHRASE")
}

func CreateToken(user *models.User) (string, error) {
	// Create a new JWT to return to user for auth
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.Username,
			"email":    user.Email,
			"exp":      time.Now().Add(time.Hour * 4).Unix(),
		})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logger.Sugar.Fatal(err)
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	// Verify JWT validity
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		logger.Sugar.Fatal(err)
		return nil, err
	}

	if !token.Valid {
		logger.Sugar.Fatal("Invalid token", "token", tokenString)
		return nil, http.ErrAbortHandler
	}

	return token, nil
}

func GetUserData(tokenString string) (*models.User, error) {
	// Get user from JWT
	_, err := VerifyToken(tokenString)
	if err != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	claims := jwt.MapClaims{}
	_, claimsErr := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if claimsErr != nil {
		logger.Sugar.Error(err)
		return nil, err
	}

	return &models.User{
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
	}, nil
}
