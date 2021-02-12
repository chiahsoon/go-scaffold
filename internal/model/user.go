package model

import (
	"github.com/pkg/errors"
)

type User struct {
	Base
	Name     string
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}

type LoginRequest struct {
	Username string
	Password string
}

type SignupRequest struct {
	Email    string
	Name     string
	Username string
	Password string
}

func CreateUser(name, email, username, password string) (*User, error) {
	user := User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: password,
	}

	if ret := DB.Create(&user); ret.Error != nil {
		return nil, errors.Wrapf(ret.Error, DBError)
	}

	return &user, nil
}

func QueryUserByUsername(username string) (User, error) {
	queryString := "username = ?"
	return queryUser(queryString, username)
}

func queryUser(query string, args ...interface{}) (User, error) {
	var user User
	if ret := DB.Where(query, args...).First(&user); ret.Error != nil {
		return user, errors.Wrapf(ret.Error, DBError)
	}

	return user, nil
}
