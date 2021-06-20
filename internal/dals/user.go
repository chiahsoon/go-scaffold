package dals

import (
	"github.com/chiahsoon/go_scaffold/internal/models"
	"gorm.io/gorm"
)

type UserDAL struct{}

func (dal UserDAL) CreateUser(tx *gorm.DB, user models.User) error {
	if ret := tx.Create(&user); ret.Error != nil {
		return models.NewInternalServerError(ret.Error.Error())
	}

	return nil
}

func (dal UserDAL) GetUserByEmail(tx *gorm.DB, email string) (*models.User, error) {
	queryString := "email = ?"
	return queryUser(tx, queryString, email)
}

func (dal UserDAL) GetUserByID(tx *gorm.DB, id string) (*models.User, error) {
	queryString := "id = ?"
	return queryUser(tx, queryString, id)
}

func queryUser(tx *gorm.DB, query string, args ...interface{}) (*models.User, error) {
	var user models.User
	if ret := tx.Where(query, args...).First(&user); ret.Error != nil {
		return nil, ret.Error
	}

	return &user, nil
}
