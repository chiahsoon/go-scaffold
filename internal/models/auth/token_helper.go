package auth

import (
	"github.com/chiahsoon/go_scaffold/internal"
	"time"

	"github.com/chiahsoon/go_scaffold/internal/models"
	"github.com/dgrijalva/jwt-go"
)

// TODO: Add logging

const (
	ClaimsUserIDKeyName = "userId"
	ClaimsExpiryKeyName = "exp"
)

var signingMethod = jwt.SigningMethodHS512

// Generates access token
func generateToken(userid string, tokenSecret string, expiryMinutes int) (string, error) {
	// TODO: Add parameter for claims

	claims := jwt.MapClaims{}
	claims[ClaimsUserIDKeyName] = userid
	claims[ClaimsExpiryKeyName] = time.Now().Add(time.Minute * time.Duration(expiryMinutes)).Format(time.UnixDate)
	at := jwt.NewWithClaims(signingMethod, claims)

	tokenSecretBuf := []byte(tokenSecret)
	token, err := at.SignedString(tokenSecretBuf)
	if err != nil {
		return "", models.NewInternalServerError(internal.JwtSigningError)
	}

	return token, nil
}

// Check if access token has expired
func hasExpired(tokenString, tokenSecret string) (bool, error) {
	claims, err := extractTokenClaims(tokenString, tokenSecret)
	if err != nil {
		return false, err
	}

	expiry := claims[ClaimsExpiryKeyName]
	expiryTime, err := time.Parse(time.UnixDate, expiry.(string))
	if err != nil {
		return false, models.NewInternalServerError(internal.TimeParseError)
	}

	if expiryTime.Before(time.Now()) {
		return true, nil
	}

	return false, nil
}

// Get claims in access token
func extractTokenClaims(tokenString, tokenSecret string) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	jwtToken, err := tokenStringToJwtToken(tokenString, tokenSecret)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return nil, models.NewUnauthorizedError(internal.InvalidJwtToken)
	}

	_, err = time.Parse(time.UnixDate, claims[ClaimsExpiryKeyName].(string))
	if err != nil {
		return nil, models.NewInternalServerError(internal.FailedToParseJwtToken)
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
			return nil, models.NewInternalServerError(internal.JwtSigningError)
		}
		return cookieTokenSecretBuf, nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, models.NewInternalServerError(internal.FailedToParseJwtToken)
	}

	return token, nil
}
