package handlers

import (
	"github.com/chiahsoon/go_scaffold/internal/models/auth"
	"net/http"

	"github.com/chiahsoon/go_scaffold/internal/models"
	"github.com/chiahsoon/go_scaffold/web/helper"
	"github.com/gin-gonic/gin"
)

const (
	AccessTokenCookieKeyName = "access_token"
)

var handler = auth.TokenHandler{}

func Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "<h1>Welcome!</h1>")
}

func Login(ctx *gin.Context) {
	var req models.EmailLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.BadRequestResponse(ctx, err)
		return
	}

	user, accessToken, err := handler.Login(req.Email, req.Password)
	if err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}

	helper.SetCookieSecurely(ctx, AccessTokenCookieKeyName, accessToken)
	helper.SuccessResponse(ctx, user)
}

func Signup(ctx *gin.Context) {
	var req models.SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.BadRequestResponse(ctx, err)
		return
	}

	user, accessToken, err := handler.Register(req.Name, req.Username, req.Email, req.Password)
	if err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}

	helper.SetCookieSecurely(ctx, AccessTokenCookieKeyName, accessToken)
	helper.SuccessResponse(ctx, user)
}

func Logout(ctx *gin.Context) {
	helper.RemoveCookieSafely(ctx, AccessTokenCookieKeyName)
	helper.SuccessResponse(ctx, "User logged out.")
}

func CurrentUser(ctx *gin.Context) {
	accessTokenString, _ := helper.GetValidCookie(ctx, AccessTokenCookieKeyName) // Assume no error
	user, err := handler.GetCurrentUser(accessTokenString)
	if err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}

	helper.SuccessResponse(ctx, user)
}

func IsAuthenticated(ctx *gin.Context) {
	accessTokenString, err := helper.GetValidCookie(ctx, AccessTokenCookieKeyName)
	if err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		ctx.Abort()
		return
	}

	if err := handler.ValidateToken(accessTokenString); err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		ctx.Abort()
		return
	}
}
