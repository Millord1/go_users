package unit_test

import (
	"microservices/models"
	"testing"
)

func TestPwHashed(t *testing.T) {

	pw := "BadPw1234?"
	mockUser := models.User{
		Username: "millord",
		Email:    "test@test.com",
		Password: pw,
	}

	user, err := mockUser.CheckPasswordIsHashed()
	if err != nil {
		t.Error(err)
	}

	if pw == user.Password {
		t.Error("password not properly hashed")
	}

	if !user.VerifyPassword(pw) {
		t.Error("password cant be verified")
	}

}

func TestUserLogin(t *testing.T) {

}
