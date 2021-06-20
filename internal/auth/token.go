package auth

import (
	"github.com/chiahsoon/go_scaffold/internal/models"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenConfigSecretKeyName = "access_token_secret"
	AccessTokenExpiryMinutes       = 60 // 1 Hour
)

type AuthService struct{}

func (h AuthService) ExchangePasswordForAccessToken(password string, user *models.User) (string, error) {

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", models.NewUnauthorizedError(InvalidPassword)
	}

	tokenSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	token, err := GenerateToken(user.ID, tokenSecret, AccessTokenExpiryMinutes)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (h AuthService) Logout() error {
	return nil
}

func (h AuthService) GetAccessTokenForNewUser(user *models.User) (string, error) {
	accessTokenSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	token, err := GenerateToken(user.ID, accessTokenSecret, AccessTokenExpiryMinutes)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (h AuthService) ValidateToken(accessTokenString string) error {
	accessTokenSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	expired, err := HasExpired(accessTokenString, accessTokenSecret)
	if err != nil || expired {
		return models.NewUnauthorizedError(InvalidAccessToken)
	}

	return nil
}

func (h AuthService) GetUserIDFromAccessToken(accessToken string) (string, error) {
	accessTokenSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	claims, err := ExtractTokenClaims(accessToken, accessTokenSecret)
	if err != nil {
		return "", err
	}

	return claims[ClaimsUserIDKeyName].(string), nil
}
