package handlers

import (
	"net/http"

	"github.com/chiahsoon/go_scaffold/internal/models"
	"github.com/chiahsoon/go_scaffold/internal/models/users"
	"github.com/chiahsoon/go_scaffold/web/helper"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "<h1>Welcome!</h1>")
}

func Login(ctx *gin.Context) {
	var req users.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.BadRequestResponse(ctx, err)
		return
	}

	user, err := users.QueryUserByUsername(req.Username)
	if err != nil {
		helper.UnauthorizedResponse(ctx, err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		helper.UnauthorizedResponse(ctx, models.NewUnauthorizedError(models.InvalidPassword))
		return
	}

	// Generate Access Token
	if err := helper.GenerateAndSetAccessTokenInCookie(ctx, user.ID); err != nil {
		helper.InternalServerErrorResponse(ctx, err)
		return
	}

	// Generate Refresh Token
	if err := helper.GenerateAndSetRefreshTokenInCookie(ctx, user.ID); err != nil {
		helper.InternalServerErrorResponse(ctx, err)
		return
	}

	helper.SuccessResponse(ctx, user)
}

func Signup(ctx *gin.Context) {
	var req users.SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.BadRequestResponse(ctx, err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		helper.InternalServerErrorResponse(ctx, models.NewInternalServerError(models.BcryptHashError))
		return
	}

	user, err := users.CreateUser(req.Name, req.Email, req.Username, string(hash))
	if err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}

	// Generate Access Token
	if err := helper.GenerateAndSetAccessTokenInCookie(ctx, user.ID); err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}

	// Generate Refresh Token
	if err := helper.GenerateAndSetRefreshTokenInCookie(ctx, user.ID); err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}

	helper.SuccessResponse(ctx, user)
}

func Logout(ctx *gin.Context) {
	helper.RemoveTokenInCookie(ctx, helper.AccessTokenCookieKeyName)
	helper.RemoveTokenInCookie(ctx, helper.RefreshTokenCookieKeyName)
	helper.SuccessResponse(ctx, "User logged out.")
}
