package auth

import (
	"time"

	"github.com/spf13/viper"

	"github.com/chiahsoon/go_scaffold/internal/models"
	"github.com/dgrijalva/jwt-go"
)

// TODO: Add logging

const (
	ClaimsUserIDKeyName = "userId"
	ClaimsExpiryKeyName = "exp"

	RefreshTokenConfigSecretKeyName = "refresh_token_secret"
	RefreshTokenExpiryMinutes       = 525600 // 1 Year
)

var signingMethod = jwt.SigningMethodHS512

func generateToken(userid string, tokenSecret string, expiryMinutes int) (string, error) {
	// TODO: Add parameter for claims

	claims := jwt.MapClaims{}
	claims[ClaimsUserIDKeyName] = userid
	claims[ClaimsExpiryKeyName] = time.Now().Add(time.Minute * time.Duration(expiryMinutes)).Format(time.UnixDate)
	jwtToken := jwt.NewWithClaims(signingMethod, claims)

	tokenSecretBuf := []byte(tokenSecret)
	token, err := jwtToken.SignedString(tokenSecretBuf)
	if err != nil {
		return "", models.NewInternalServerError(JwtSigningError)
	}

	return token, nil
}

func generateTokenPair(user *models.User) (string, string, error) {
	// Generate Access Token
	atSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	accessToken, err := generateToken(user.ID, atSecret, AccessTokenExpiryMinutes)
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	rtSecret := viper.GetString(RefreshTokenConfigSecretKeyName)
	refreshToken, err := generateToken(user.ID, rtSecret, RefreshTokenExpiryMinutes)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func tokenIsWellSigned(tokenString, tokenSecret string) error {
	_, err := ExtractTokenClaims(tokenString, tokenSecret)
	return err
}

func HasExpired(tokenString, tokenSecret string) (bool, error) {
	claims, err := ExtractTokenClaims(tokenString, tokenSecret)

	if err != nil {
		return false, err
	}

	expiry := claims[ClaimsExpiryKeyName]
	expiryTime, err := time.Parse(time.UnixDate, expiry.(string))
	if err != nil {
		return false, models.NewInternalServerError(TimeParseError)
	}

	if expiryTime.Before(time.Now()) {
		return true, nil
	}

	return false, nil
}

func ExtractTokenClaims(tokenString, tokenSecret string) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	jwtToken, err := tokenStringToJwtToken(tokenString, tokenSecret)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return nil, models.NewUnauthorizedError(InvalidJwtToken)
	}

	_, err = time.Parse(time.UnixDate, claims[ClaimsExpiryKeyName].(string))
	if err != nil {
		return nil, models.NewInternalServerError(FailedToParseJwtToken)
	}

	for claim, value := range claims {
		results[claim] = value
	}

	return results, nil
}

// =====================================================================================================================
// Helper Methods
// =====================================================================================================================

func tokenStringToJwtToken(tokenString, tokenSecret string) (*jwt.Token, error) {
	cookieTokenSecretBuf := []byte(tokenSecret)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.NewInternalServerError(JwtSigningError)
		}
		return cookieTokenSecretBuf, nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, models.NewInternalServerError(FailedToParseJwtToken)
	}

	return token, nil
}
