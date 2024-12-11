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

	hashedPwUser, hashErr := mockUser.HashPassword()
	if assert.NoError(t, hashErr) {
		assert.True(t, checkisHashed(hashedPwUser.Password))
	}

	user, err := mockUser.CheckPasswordIsHashed()
	if assert.NoError(t, err) {
		assert.True(t, checkisHashed(user.Password))
	}

	assert.NotEqual(t, pw, user.Password)
	assert.True(t, user.VerifyPassword(pw))
}
