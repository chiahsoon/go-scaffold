package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TODO: Add logging

const (
	ClaimsUserIDKeyName = "userId"
	ClaimsExpiryKeyName = "exp"
)

var signingMethod = jwt.SigningMethodHS512

// Returns user id from access token
func GetUserIDFromAccessToken(accessToken, accessTokenSecret string) (string, error) {
	claims, err := ExtractTokenClaims(accessToken, accessTokenSecret)
	if err != nil {
		return "", err
	}

	return claims[ClaimsUserIDKeyName].(string), nil
}

// Generates access or refresh token
func GenerateToken(userid string, tokenSecret string, isAccessToken bool, expiryMinutes int) (string, error) {
	// TODO: Add parameter for claims

	claims := jwt.MapClaims{}
	if isAccessToken {
		claims[ClaimsUserIDKeyName] = userid
	}

	claims[ClaimsExpiryKeyName] = time.Now().Add(time.Minute * time.Duration(expiryMinutes)).Format(time.UnixDate)
	at := jwt.NewWithClaims(signingMethod, claims)

	tokenSecretBuf := []byte(tokenSecret)
	token, err := at.SignedString(tokenSecretBuf)
	if err != nil {
		return "", NewInternalServerError(JwtSigningError)
	}

	return token, nil
}

// Returns new access token in exchange for expired access token and valid refresh token
func Refresh(accessTokenString, accessTokenSecret,
	refreshTokenString, refreshTokenSecret string, accessTokenExpiryMinutes int) (string, error) {
	accessTokenClaims, err := ExtractTokenClaims(accessTokenString, accessTokenSecret)
	if err != nil {
		return "", err
	}

	userID := accessTokenClaims[ClaimsUserIDKeyName].(string)
	if expired, err := HasExpired(refreshTokenString, refreshTokenSecret); err != nil {
		return "", err
	} else if expired {
		err = NewUnauthorizedError(ExpiredJwtToken)
		return "", err
	}

	token, err := GenerateToken(userID, accessTokenSecret, true, accessTokenExpiryMinutes)
	if err != nil {
		return "", nil
	}

	return token, nil
}

// Check if access or refresh token has expired
func HasExpired(tokenString, tokenSecret string) (bool, error) {
	claims, err := ExtractTokenClaims(tokenString, tokenSecret)
	if err != nil {
		return false, err
	}

	expiry := claims[ClaimsExpiryKeyName]
	expiryTime, err := time.Parse(time.UnixDate, expiry.(string))
	if err != nil {
		return false, NewInternalServerError(TimeParseError)
	}

	if expiryTime.Before(time.Now()) {
		return true, nil
	}

	return false, nil
}

// Get claims in access or refresh token
func ExtractTokenClaims(tokenString, tokenSecret string) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	jwtToken, err := tokenStringToJwtToken(tokenString, tokenSecret)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return nil, NewUnauthorizedError(InvalidJwtToken)
	}

	_, err = time.Parse(time.UnixDate, claims[ClaimsExpiryKeyName].(string))
	if err != nil {
		return nil, NewInternalServerError(FailedToParseJwtToken)
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
			return nil, NewInternalServerError(JwtSigningError)
		}
		return cookieTokenSecretBuf, nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, NewInternalServerError(FailedToParseJwtToken)
	}

	return token, nil
}
