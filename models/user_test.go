package models

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func checkisHashed(pw string) bool {
	return pw[0:4] == "$2a$" && len(pw) > 50
}

func TestPwHashed(t *testing.T) {

	pw := gofakeit.Password(true, true, true, true, false, 10)
	mockUser := User{
		Username: gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: pw,
	}

	hashErr := mockUser.HashPassword()
	if assert.NoError(t, hashErr) {
		assert.True(t, checkisHashed(mockUser.Password))
	}

	err := mockUser.CheckPasswordIsHashed()
	if assert.NoError(t, err) {
		assert.True(t, checkisHashed(mockUser.Password))
	}

	assert.NotEqual(t, pw, mockUser.Password)
	assert.True(t, mockUser.VerifyPassword(pw))
}
