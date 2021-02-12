package helper

import (
	"net/http"
	"time"

	"github.com/chiahsoon/go_scaffold/internal/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	AccessTokenCookieKeyName       = "access_token"
	AccessTokenConfigSecretKeyName = "access_token_secret"
	AccessTokenExpiryMinutes       = 60 // 1 Hour

	RefreshTokenCookieKeyName       = "refresh_token"
	RefreshTokenConfigSecretKeyName = "refresh_token_secret"
	RefreshTokenExpiryMinutes       = 10080 // 1 Day
	ClaimsAuthorizedKeyName         = "authorized"
	ClaimsUserIDKeyName             = "userId"
	ClaimsExpiryKeyName             = "exp"

	DomainConfigKeyName       = "domain"
	CookieSecureConfigKeyName = "cookie_secure"
)

var signingMethod = jwt.SigningMethodHS512

func GenerateAndSetTokenInCookie(ctx *gin.Context, userID string, cookieKeyName string, secretConfigKeyName string,
	expiryMinutes int) error {
	token, err := generateToken(userID, secretConfigKeyName, expiryMinutes)
	if err != nil {
		return err
	}

	domain := viper.GetString(DomainConfigKeyName)
	secure := viper.GetBool(CookieSecureConfigKeyName)
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(cookieKeyName, token, 0, "/", domain, secure, true)
	return nil
}

func RemoveTokenInCookie(ctx *gin.Context, cookieKeyName string) {
	domain := viper.GetString(DomainConfigKeyName)
	secure := viper.GetBool(CookieSecureConfigKeyName)
	ctx.SetCookie(cookieKeyName, "", -1, "/", domain, secure, true)
}

func IsAuthorized(endpoint func(ctx *gin.Context)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if err := validateAccessToken(ctx); err != nil {
			return
		}

		endpoint(ctx)
	}
}

func validateAccessToken(ctx *gin.Context) error {
	accessTokenClaims, err := extractTokenClaims(ctx, AccessTokenCookieKeyName, AccessTokenConfigSecretKeyName)
	if err != nil {
		return err
	}

	if expired, err := hasExpired(accessTokenClaims[ClaimsExpiryKeyName]); err != nil {
		return err
	} else if expired {
		userID := accessTokenClaims[ClaimsUserIDKeyName].(string)
		return refresh(ctx, userID)
	}

	return nil
}

func refresh(ctx *gin.Context, userID string) error {
	refreshTokenClaims, err := extractTokenClaims(ctx, RefreshTokenCookieKeyName, RefreshTokenConfigSecretKeyName)
	if err != nil {
		return err
	}

	if expired, err := hasExpired(refreshTokenClaims[ClaimsExpiryKeyName]); err != nil {
		return err
	} else if expired {
		err = errors.New(model.ExpiredJwtToken)
		return err
	}

	err = GenerateAndSetTokenInCookie(ctx, userID, AccessTokenCookieKeyName, AccessTokenConfigSecretKeyName,
		AccessTokenExpiryMinutes)
	if err != nil {
		InternalServerErrorResponse(ctx, err)
		return err
	}

	return nil
}

func decodeToken(cookieToken, configKeyName string) (*jwt.Token, error) {
	secretKey := []byte(viper.GetString(configKeyName))

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return secretKey, nil
	}

	token, err := jwt.Parse(cookieToken, keyFunc)
	if err != nil {
		return nil, errors.Wrapf(err, model.FailedToParseJwtToken)
	}

	return token, nil
}

func extractTokenClaims(ctx *gin.Context, cookieKeyName, configKeyName string) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	cookieAccessToken, err := getTokenInCookie(ctx, cookieKeyName)
	if err != nil {
		UnauthorizedResponse(ctx, err)
		return results, err
	}

	jwtToken, err := decodeToken(cookieAccessToken, configKeyName)
	if err != nil {
		InternalServerErrorResponse(ctx, err)
		return results, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		UnauthorizedResponse(ctx, errors.New(model.InvalidJwtToken))
		return results, err
	}

	_, err = time.Parse(time.UnixDate, claims[ClaimsExpiryKeyName].(string))
	if err != nil {
		InternalServerErrorResponse(ctx, errors.Wrapf(err, model.TimeParseError))
		return results, err
	}

	for claim, value := range claims {
		results[claim] = value
	}

	return results, nil
}

func getTokenInCookie(ctx *gin.Context, cookieKeyName string) (string, error) {
	cookieToken, err := ctx.Cookie(cookieKeyName)
	if err != nil {
		if cookieKeyName == AccessTokenCookieKeyName {
			return "", errors.New(model.EmptyAccessToken)
		}
		return "", errors.New(model.EmptyRefreshToken)
	}

	if cookieToken == "" {
		if cookieKeyName == AccessTokenCookieKeyName {
			return "", errors.New(model.EmptyAccessToken)
		}
		return "", errors.New(model.EmptyRefreshToken)
	}

	return cookieToken, nil
}

func hasExpired(expiry interface{}) (bool, error) {
	expiryTime, err := time.Parse(time.UnixDate, expiry.(string))
	if err != nil {
		return true, errors.Wrapf(err, model.TimeParseError)
	}

	if expiryTime.Before(time.Now()) {
		return true, nil
	}

	return false, nil
}

func generateToken(userid string, configSecretKeyName string, expiryMinutes int) (string, error) {
	claims := jwt.MapClaims{}
	claims[ClaimsAuthorizedKeyName] = true
	if configSecretKeyName == AccessTokenConfigSecretKeyName {
		claims[ClaimsUserIDKeyName] = userid
	}

	claims[ClaimsExpiryKeyName] = time.Now().Add(time.Minute * time.Duration(expiryMinutes)).Format(time.UnixDate)
	at := jwt.NewWithClaims(signingMethod, claims)

	secretKey := []byte(viper.GetString(configSecretKeyName))
	token, err := at.SignedString(secretKey)
	if err != nil {
		return "", errors.Wrapf(err, model.JwtSigningError)
	}

	return token, nil
}

func GetCurrentUserID(ctx *gin.Context) (string, error) {
	accessTokenClaims, err := extractTokenClaims(ctx, AccessTokenCookieKeyName, AccessTokenConfigSecretKeyName)
	if err != nil {
		return "", err
	}

	return accessTokenClaims[ClaimsUserIDKeyName].(string), nil
}
