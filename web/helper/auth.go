package helper

import (
	"net/http"

	"github.com/chiahsoon/go_scaffold/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	AccessTokenCookieKeyName       = "access_token"
	AccessTokenConfigSecretKeyName = "access_token_secret"
	AccessTokenExpiryMinutes       = 60 // 1 Hour

	RefreshTokenCookieKeyName       = "refresh_token"
	RefreshTokenConfigSecretKeyName = "refresh_token_secret"
	RefreshTokenExpiryMinutes       = 10080 // 1 Day

	DomainConfigKeyName       = "domain"
	CookieSecureConfigKeyName = "cookie_secure"
)

func GenerateAndSetAccessTokenInCookie(ctx *gin.Context, userID string) error {
	return generateAndSetTokenInCookie(ctx, userID, AccessTokenCookieKeyName, AccessTokenConfigSecretKeyName,
		AccessTokenExpiryMinutes)
}

func GenerateAndSetRefreshTokenInCookie(ctx *gin.Context, userID string) error {
	return generateAndSetTokenInCookie(ctx, userID, RefreshTokenCookieKeyName, RefreshTokenConfigSecretKeyName,
		RefreshTokenExpiryMinutes)
}

func RemoveTokenInCookie(ctx *gin.Context, cookieKeyName string) {
	domain := viper.GetString(DomainConfigKeyName)
	secure := viper.GetBool(CookieSecureConfigKeyName)
	ctx.SetCookie(cookieKeyName, "", -1, "/", domain, secure, true)
}

func IsAuthorized(ctx *gin.Context) {
	accessTokenSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	accessTokenString, err := getTokenInCookie(ctx, AccessTokenCookieKeyName)
	if err != nil {
		ErrorToErrorResponse(ctx, err)
		return
	}
	expired, err := model.HasExpired(accessTokenString, accessTokenSecret)
	if err != nil {
		ErrorToErrorResponse(ctx, err)
		return
	}

	if expired {
		refreshTokenString, err := getTokenInCookie(ctx, RefreshTokenCookieKeyName)
		if err != nil {
			ErrorToErrorResponse(ctx, err)
			return
		}

		refreshTokenSecret := viper.GetString(RefreshTokenConfigSecretKeyName)
		token, err := model.Refresh(accessTokenString, accessTokenSecret,
			refreshTokenString, refreshTokenSecret, AccessTokenExpiryMinutes)
		if err != nil {
			ErrorToErrorResponse(ctx, err)
			return
		}

		setTokenInCookie(ctx, token, AccessTokenCookieKeyName)
	}
}

func GetCurrentUserID(ctx *gin.Context) (string, error) {
	accessToken, err := getTokenInCookie(ctx, AccessTokenCookieKeyName)
	if err != nil {
		return "", err
	}

	accessTokenSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	return model.GetUserIDFromAccessToken(accessToken, accessTokenSecret)
}

// =====================================================================================================================
// Helper Methods
// =====================================================================================================================

func generateAndSetTokenInCookie(ctx *gin.Context, userID string, cookieKeyName string, secretConfigKeyName string,
	expiryMinutes int) error {
	secret := viper.GetString(secretConfigKeyName)
	token, err := model.GenerateToken(userID, secret, cookieKeyName == AccessTokenCookieKeyName, expiryMinutes)
	if err != nil {
		return err
	}

	setTokenInCookie(ctx, token, cookieKeyName)
	return nil
}

func setTokenInCookie(ctx *gin.Context, token, cookieKeyName string) {
	domain := viper.GetString(DomainConfigKeyName)
	secure := viper.GetBool(CookieSecureConfigKeyName)
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(cookieKeyName, token, 0, "/", domain, secure, true)
}

func getTokenInCookie(ctx *gin.Context, cookieKeyName string) (string, error) {
	cookieToken, err := ctx.Cookie(cookieKeyName)
	if err != nil || cookieToken == "" {
		if cookieKeyName == AccessTokenCookieKeyName {
			return "", model.NewBadRequestError(model.EmptyAccessToken)
		}
		return "", model.NewBadRequestError(model.EmptyRefreshToken)
	}

	return cookieToken, nil
}
