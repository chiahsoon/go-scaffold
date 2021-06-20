package models

import (
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
