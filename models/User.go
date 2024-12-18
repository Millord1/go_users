package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey" json:"-"`
	Username string `gorm:"type:varchar(45);uniqueIndex:username_unique;not null;<-" json:"username" fake:"{username}"`
	Email    string `gorm:"type:varchar(60);uniqueIndex:email_unique;not null;<-" json:"-" fake:"email"`
	Password string `gorm:"type:varchar(65);not null;<-" json:"-" fake:"password"`
	Totp     string `gorm:"type:varchar(60);<-" json:"-"`
}

func (user *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) HashTotp() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Totp), 14)
	if err != nil {
		return err
	}
	user.Totp = string(bytes)
	return nil
}

func (user *User) VerifyTotp() bool {
	return true
}

func (user User) VerifyPassword(pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pw))
	return err == nil
}

func (user *User) CheckPasswordIsHashed() error {
	// Check password length and first chars
	if !(user.Password[0:4] == "$2a$" && len(user.Password) > 50) {
		if err := user.HashPassword(); err != nil {
			return err
		}
	}
	return nil
}
