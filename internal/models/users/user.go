package users

import (
	"github.com/chiahsoon/go_scaffold/internal/models"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type User struct {
	models.Base
	Name     string
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string `json:"-"`
}

func (u User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", u.Name)
	enc.AddString("username", u.Username)
	enc.AddString("email", u.Email)
	return nil
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
	user := User {
		Name:     name,
		Username: username,
		Email:    email,
		Password: password,
	}

	if ret := models.DB.Create(&user); ret.Error != nil {
		zap.L().Error(
			models.DBError,
			zap.Object("user", &user),
		)
		return nil, models.NewInternalServerError(models.DBError)
	}

	return &user, nil
}

func QueryUserByUsername(username string) (User, error) {
	queryString := "username = ?"
	return queryUser(queryString, username)
}

func queryUser(query string, args ...interface{}) (User, error) {
	var user User
	if ret := models.DB.Where(query, args...).First(&user); ret.Error != nil {
		zap.L().Error(
			models.DBError,
			zap.String("query", query),
		)
		return user, models.NewInternalServerError(models.DBError)
	}

	return user, nil
}
