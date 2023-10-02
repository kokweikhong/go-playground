package service

import (
	"errors"

	"github.com/kokweikhong/go-playground/chi-middleware-jwt-auth/model"
)

type User struct{}

func NewUser() *User {
	return &User{}
}

// get user by username
func (u *User) GetUserByUsername(username string) (*model.User, error) {
	for _, user := range model.Users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("User not found")
}
