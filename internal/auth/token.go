package auth

import (
	"github.com/chiahsoon/go_scaffold/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

const (
	AccessTokenConfigSecretKeyName = "access_token_secret"
	AccessTokenExpiryMinutes       = 60 // 1 Hour
)

type AuthService struct{}

func (h AuthService) Logout() error {
	return nil
}

func (h AuthService) GetTokenPairForUser(user *models.User) (string, string, error) {
	accessToken, refreshToken, err := generateTokenPair(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (h AuthService) ValidateToken(accessTokenString string) error {
	accessTokenSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	err := isTokenValid(accessTokenString, accessTokenSecret)
	if err != nil {
		return err
	}

	return nil
}

func (h AuthService) GetUserIDFromAccessToken(accessToken string) (string, error) {
	accessTokenSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	jwtToken, err := tokenStringToJwtToken(accessToken, accessTokenSecret)
	if err != nil {
		return "", err
	}

	claims, err := extractTokenClaims(jwtToken)
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}

func (h AuthService) Refresh(accessTokenString, refreshTokenString string) (string, error) {
	refreshTokenSecret := viper.GetString(RefreshTokenConfigSecretKeyName)
	userID, err := h.GetUserIDFromAccessToken(accessTokenString)
	if err != nil {
		return "", err
	}

	if !isTokenSignatureValid(refreshTokenString, refreshTokenSecret) {
		return "", jwt.NewValidationError(InvalidJwtToken, jwt.ValidationErrorSignatureInvalid)
	}

	if !IsTokenExpired(refreshTokenString, refreshTokenSecret) {
		return "", jwt.NewValidationError(InvalidJwtToken, jwt.ValidationErrorExpired)
	}

	accessTokenSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	return generateToken(userID, accessTokenSecret, AccessTokenExpiryMinutes)
}


func (h AuthService) AccessTokenHasExpired(tokenString string) bool {
	accessTokenSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	return IsTokenExpired(tokenString, accessTokenSecret)
}

func (h AuthService) RefreshTokenHasExpired(tokenString string) bool {
	refreshTokenSecret := viper.GetString(RefreshTokenConfigSecretKeyName)
	return IsTokenExpired(tokenString, refreshTokenSecret)
}
