package handlers

import (
	"github.com/pkg/errors"
	"net/http"

	"github.com/chiahsoon/go_scaffold/internal"
	"github.com/chiahsoon/go_scaffold/internal/models"
	"github.com/chiahsoon/go_scaffold/web/helper"
	"github.com/chiahsoon/go_scaffold/web/view"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenCookieKeyName = "access_token"
	RefreshTokenCookieKeyName = "refresh_token"
)

func Home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "<h1>Welcome!</h1>")
}

func Login(ctx *gin.Context) {
	var req view.EmailLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrorToErrorResponse(ctx, models.NewBadRequestError(err.Error()))
		return
	}

	user, err := internal.UserService.GetUserByEmail(internal.DB, req.Email)
	if err != nil {
		helper.UnauthorizedResponse(ctx, err)
		return
	}

	accessToken, refreshToken, err := internal.AuthService.ExchangePasswordForTokenPair(req.Password, user)
	if err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}

	uwt := &models.UserRefreshToken{UserID: user.ID, RefreshToken: refreshToken}
	if err := internal.UserRefreshTokenService.Upsert(internal.DB, uwt); err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}


	helper.SetCookieSecurely(ctx, AccessTokenCookieKeyName, accessToken)
	helper.SetCookieSecurely(ctx, RefreshTokenCookieKeyName, refreshToken)
	helper.SuccessResponse(ctx, user)
}

func Signup(ctx *gin.Context) {
	var req view.SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.BadRequestResponse(ctx, err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		helper.ErrorToErrorResponse(ctx, models.NewInternalServerError(err.Error()))
		return
	}

	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
	}

	if err = internal.UserService.CreateUser(internal.DB, user); err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}

	accessToken, refreshToken, err := internal.AuthService.GetTokenPairForNewUser(&user)
	if err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}

	uwt := &models.UserRefreshToken{UserID: user.ID, RefreshToken: refreshToken}
	if err := internal.UserRefreshTokenService.Upsert(internal.DB, uwt); err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}

	helper.SetCookieSecurely(ctx, AccessTokenCookieKeyName, accessToken)
	helper.SetCookieSecurely(ctx, RefreshTokenCookieKeyName, refreshToken)
	helper.SuccessResponse(ctx, user)
}

func Logout(ctx *gin.Context) {
	if err := internal.AuthService.Logout(); err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}
	helper.RemoveCookieSafely(ctx, AccessTokenCookieKeyName)
	helper.RemoveCookieSafely(ctx, RefreshTokenCookieKeyName)
	helper.SuccessResponse(ctx, "User logged out.")
}

func CurrentUser(ctx *gin.Context) {
	accessTokenString := helper.GetValidCookie(ctx, AccessTokenCookieKeyName) // Assume no error
	userID, err := internal.AuthService.GetUserIDFromAccessToken(accessTokenString)
	if err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}

	user, err := internal.UserService.GetUserByID(internal.DB, userID)
	if err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}

	helper.SuccessResponse(ctx, user)
}

func IsAuthenticated(ctx *gin.Context) {
	accessTokenString := helper.GetValidCookie(ctx, AccessTokenCookieKeyName)

	if err := internal.AuthService.ValidateToken(accessTokenString); err == nil {
		return
	}

	ATExpired := internal.AuthService.AccessTokenHasExpired(accessTokenString)
	if !ATExpired {
		helper.InternalServerErrorResponse(ctx, errors.New("failed to parse jwt"))
		return
	}

	// Refresh access token
	refreshTokenString := helper.GetValidCookie(ctx, RefreshTokenCookieKeyName)
	RTExpired := internal.AuthService.RefreshTokenHasExpired(refreshTokenString)
	if RTExpired {
		helper.ErrorToErrorResponse(ctx, models.NewBadRequestError("refresh token has expired"))
		return
	}

	newAT, err := internal.AuthService.Refresh(accessTokenString, refreshTokenString)
	if err != nil {
		helper.ErrorToErrorResponse(ctx, err)
		return
	}

	helper.SetCookieSecurely(ctx, AccessTokenCookieKeyName, newAT)
	return
}
