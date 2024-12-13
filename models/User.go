package models

import (
	"microservices/encryption"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(45);uniqueIndex:username_unique;not null;<-" json:"username" fake:"{username}"`
	Email    string `gorm:"type:varchar(60);uniqueIndex:email_unique;not null;<-" json:"-" fake:"email"`
	Password string `gorm:"type:varchar(65);not null;<-" json:"-" fake:"password"`
	Totp     string `gorm:"type:varchar(60);<-" json:"-"`
}

func (user *User) HashPassword() (*User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(bytes)
	return user, err
}

func (user *User) HashTotp() (*User, error) {
	enc, err := encryption.EncryptData(user.Totp)
	user.Totp = enc
	return user, err
}

func (user User) VerifyPassword(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pw))
	return err == nil
}

func (user *User) CheckPasswordIsHashed() (*User, error) {
	// Check password length and first chars
	if !(user.Password[0:4] == "$2a$" && len(user.Password) > 50) {
		return user.HashPassword()
	}
	return user, nil
}
