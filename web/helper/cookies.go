package helper

import (
	"net/http"

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

func GetValidCookie(ctx *gin.Context, cookieKeyName string) string {
	cookieToken, err := ctx.Cookie(cookieKeyName)
	if err != nil || cookieToken == "" {
		return ""
	}

	return cookieToken
}

func RemoveCookieSafely(ctx *gin.Context, cookieKeyName string) {
	domain := viper.GetString(DomainConfigKeyName)
	secure := viper.GetBool(CookieSecureConfigKeyName)
	ctx.SetCookie(cookieKeyName, "", -1, "/", domain, secure, true)
}
