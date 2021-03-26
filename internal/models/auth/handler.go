package auth

import (
	"github.com/chiahsoon/go_scaffold/internal"
	"github.com/chiahsoon/go_scaffold/internal/models"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenConfigSecretKeyName = "access_token_secret"
	AccessTokenExpiryMinutes       = 60 // 1 Hour
)

type TokenHandler struct{}

func (h *TokenHandler) Login(email, password string) (*models.User, string, error) {
	user, err := models.QueryUserByEmail(email)
	if err != nil {
		return user, "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, "", models.NewUnauthorizedError(internal.InvalidPassword)
	}

	tokenSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	token, err := generateToken(user.ID, tokenSecret, AccessTokenExpiryMinutes)
	if err != nil {
		return user, "", err
	}

	return user, token, nil
}

func (h *TokenHandler) Logout() error {
	return nil
}

func (h *TokenHandler) Register(name, username, email, password string) (*models.User, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, "", models.NewInternalServerError(internal.BcryptHashError)
	}

	user, err := models.CreateUser(name, email, username, string(hash))
	return user, "", nil
}

func (h *TokenHandler) ValidateToken(accessTokenString string) error {
	accessTokenSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	expired, err := hasExpired(accessTokenString, accessTokenSecret)
	if err != nil || expired {
		return models.NewUnauthorizedError(internal.InvalidAccessToken)
	}

	return nil
}

func (h *TokenHandler) GetCurrentUserID(accessToken string) (string, error) {
	accessTokenSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	claims, err := extractTokenClaims(accessToken, accessTokenSecret)
	if err != nil {
		return "", err
	}

	return claims[ClaimsUserIDKeyName].(string), nil
}

func (h *TokenHandler) GetCurrentUser(accessToken string) (*models.User, error) {
	userID, err := h.GetCurrentUserID(accessToken)
	if err != nil {
		return nil, err
	}

	return models.QueryUserByID(userID)
}
