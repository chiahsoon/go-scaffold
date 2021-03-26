package helper

import (
	"github.com/chiahsoon/go_scaffold/internal"
	"net/http"

	"github.com/chiahsoon/go_scaffold/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	DomainConfigKeyName       = "domain"
	CookieSecureConfigKeyName = "cookie_secure"
)

func SetCookieSecurely(ctx *gin.Context, cookieKeyName, token string) {
	domain := viper.GetString(DomainConfigKeyName)
	secure := viper.GetBool(CookieSecureConfigKeyName)
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(cookieKeyName, token, 0, "/", domain, secure, true)
}

func GetValidCookie(ctx *gin.Context, cookieKeyName string) (string, error) {
	cookieToken, err := ctx.Cookie(cookieKeyName)
	if err != nil || cookieToken == "" {
		return "", models.NewUnauthorizedError(internal.InvalidCookieValue)
	}

	return cookieToken, nil
}

func RemoveCookieSafely(ctx *gin.Context, cookieKeyName string) {
	domain := viper.GetString(DomainConfigKeyName)
	secure := viper.GetBool(CookieSecureConfigKeyName)
	ctx.SetCookie(cookieKeyName, "", -1, "/", domain, secure, true)
}
