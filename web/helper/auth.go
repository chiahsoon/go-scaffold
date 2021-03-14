package helper

import (
	"net/http"

	"github.com/chiahsoon/go_scaffold/internal/models"
	"github.com/chiahsoon/go_scaffold/internal/models/auth"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	AccessTokenCookieKeyName       = "access_token"
	AccessTokenConfigSecretKeyName = "access_token_secret"
	AccessTokenExpiryMinutes       = 60 // 1 Hour

	DomainConfigKeyName       = "domain"
	CookieSecureConfigKeyName = "cookie_secure"
)

var (
	accessTokenSecret = ""
)

func Init() {
	accessTokenSecret = viper.GetString(AccessTokenConfigSecretKeyName)
}

func RemoveTokenInCookie(ctx *gin.Context, cookieKeyName string) {
	domain := viper.GetString(DomainConfigKeyName)
	secure := viper.GetBool(CookieSecureConfigKeyName)
	ctx.SetCookie(cookieKeyName, "", -1, "/", domain, secure, true)
}

func IsAuthorized(ctx *gin.Context) {
	accessTokenString, err := getTokenInCookie(ctx)
	if err != nil {
		ErrorToErrorResponse(ctx, err)
		ctx.Abort()
		return
	}

	expired, err := auth.HasExpired(accessTokenString, accessTokenSecret)
	if err != nil || expired {
		ErrorToErrorResponse(ctx, err)
		ctx.Abort()
		return
	}
}

func GetCurrentUserID(ctx *gin.Context) (string, error) {
	accessToken, err := getTokenInCookie(ctx)
	if err != nil {
		return "", err
	}

	return auth.GetUserIDFromAccessToken(accessToken, accessTokenSecret)
}

// =====================================================================================================================
// Helper Methods
// =====================================================================================================================

func GenerateAndSetAccessTokenInCookie(ctx *gin.Context, userID string) error {
	token, err := auth.GenerateToken(userID, accessTokenSecret, AccessTokenExpiryMinutes)
	if err != nil {
		return err
	}

	setTokenInCookie(ctx, token, AccessTokenCookieKeyName)
	return nil
}

func setTokenInCookie(ctx *gin.Context, token, cookieKeyName string) {
	domain := viper.GetString(DomainConfigKeyName)
	secure := viper.GetBool(CookieSecureConfigKeyName)
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(cookieKeyName, token, 0, "/", domain, secure, true)
}

func getTokenInCookie(ctx *gin.Context) (string, error) {
	cookieToken, err := ctx.Cookie(AccessTokenCookieKeyName)
	if err != nil || cookieToken == "" {
		return "", models.NewUnauthorizedError(models.EmptyAccessToken)
	}

	return cookieToken, nil
}
