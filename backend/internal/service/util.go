package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrUsernameExists = errors.New("username already exists")

func bcryptHash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
