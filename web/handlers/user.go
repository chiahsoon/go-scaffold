package handlers

import (
	"net/http"

	"github.com/chiahsoon/go_scaffold/internal/model"
	"github.com/chiahsoon/go_scaffold/web/helper"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "<h1>Welcome!</h1>")
}

func Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.BadRequestResponse(ctx, err)
		return
	}

	user, err := model.QueryUserByUsername(req.Username)
	if err != nil {
		helper.UnauthorizedResponse(ctx, err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		helper.UnauthorizedResponse(ctx, errors.Wrapf(err, model.InvalidPassword))
		return
	}

	// Generate Access Token
	if err := helper.GenerateAndSetTokenInCookie(ctx, user.ID, helper.AccessTokenCookieKeyName,
		helper.AccessTokenConfigSecretKeyName, helper.AccessTokenExpiryMinutes); err != nil {
		helper.InternalServerErrorResponse(ctx, err)
		return
	}

	// Generate Refresh Token
	if err := helper.GenerateAndSetTokenInCookie(ctx, user.ID, helper.RefreshTokenCookieKeyName,
		helper.RefreshTokenConfigSecretKeyName, helper.RefreshTokenExpiryMinutes); err != nil {
		helper.InternalServerErrorResponse(ctx, err)
		return
	}

	helper.SuccessResponse(ctx, user)
}

func Signup(ctx *gin.Context) {
	var req model.SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.BadRequestResponse(ctx, err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		helper.InternalServerErrorResponse(ctx, errors.Wrapf(err, model.BcryptHashError))
		return
	}

	user, err := model.CreateUser(req.Name, req.Email, req.Username, string(hash))
	if err != nil {
		helper.InternalServerErrorResponse(ctx, err)
		return
	}

	// Generate Access Token
	if err := helper.GenerateAndSetTokenInCookie(ctx, user.ID, helper.AccessTokenCookieKeyName,
		helper.AccessTokenConfigSecretKeyName, helper.AccessTokenExpiryMinutes); err != nil {
		helper.InternalServerErrorResponse(ctx, err)
		return
	}

	// Generate Refresh Token
	if err := helper.GenerateAndSetTokenInCookie(ctx, user.ID, helper.RefreshTokenCookieKeyName,
		helper.RefreshTokenConfigSecretKeyName, helper.RefreshTokenExpiryMinutes); err != nil {
		helper.InternalServerErrorResponse(ctx, err)
		return
	}

	helper.SuccessResponse(ctx, user)
}

func Logout(ctx *gin.Context) {
	helper.RemoveTokenInCookie(ctx, helper.AccessTokenCookieKeyName)
	helper.RemoveTokenInCookie(ctx, helper.RefreshTokenCookieKeyName)
	helper.SuccessResponse(ctx, "User logged out.")
}
