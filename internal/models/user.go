package models

type User struct {
	Base
	Name     string `json:"name"`
	Username string `gorm:"unique" json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`
}

type UserWithTokens struct {
	*User
	AccessToken  string
	RefreshToken string
}