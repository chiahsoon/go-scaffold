package models

import (
	"fmt"
	"github.com/chiahsoon/go_scaffold/internal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type User struct {
	Base
	Name     string `json:"name"`
	Username string `gorm:"unique" json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`
}

func (u User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", u.Name)
	enc.AddString("username", u.Username)
	enc.AddString("email", u.Email)
	return nil
}

type EmailLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignupRequest struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateUser(name, email, username, password string) (*User, error) {
	user := User{
		Name:     name,
		Username: username,
		Email:    email,
		Password: password,
	}

	if ret := DB.Create(&user); ret.Error != nil {
		zap.L().Error(
			internal.DBError,
			zap.Object("user", &user),
		)
		return nil, NewInternalServerError(internal.DBError)
	}

	return &user, nil
}

func QueryUserByEmail(email string) (*User, error) {
	queryString := "email = ?"
	return queryUser(queryString, email)
}

func QueryUserByID(id string) (*User, error) {
	queryString := "id = ?"
	return queryUser(queryString, id)
}

func queryUser(query string, args ...interface{}) (*User, error) {
	var user User
	if ret := DB.Where(query, args...).First(&user); ret.Error != nil {
		zap.L().Error(
			fmt.Sprintf(ret.Error.Error()),
			zap.String("query", query),
			zap.String("args", fmt.Sprintf("%v", args...)),
		)
		return nil, NewInternalServerError(internal.DBError)
	}

	return &user, nil
}
