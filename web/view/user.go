package view

import "github.com/chiahsoon/go_scaffold/internal/models"

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

type UserWithTokens struct {
	*models.User
	AccessToken  string
	RefreshToken string
}