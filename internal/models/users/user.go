package users

import "github.com/chiahsoon/go_scaffold/internal/models"

type User struct {
	models.Base
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

	if ret := models.DB.Create(&user); ret.Error != nil {
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
		return user, models.NewInternalServerError(models.DBError)
	}

	return user, nil
}
