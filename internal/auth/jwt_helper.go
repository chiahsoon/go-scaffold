package auth

import (
	"fmt"
	"github.com/spf13/viper"
	"time"

	"github.com/chiahsoon/go_scaffold/internal/models"
	"github.com/dgrijalva/jwt-go"
)

// TODO: Add logging

var signingMethod = jwt.SigningMethodHS512

const (
	RefreshTokenConfigSecretKeyName = "refresh_token_secret"
	RefreshTokenExpiryMinutes       = 525600 // 1 Year
)


func GenerateToken(userID string, tokenSecret string, expiryMinutes int) (string, error) {
	// TODO: Add parameter for claims
	claims := jwt.StandardClaims{
		Id:        userID,
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(expiryMinutes)).Unix(),
		Subject:   userID,
	}
	at := jwt.NewWithClaims(signingMethod, claims)

	tokenSecretBuf := []byte(tokenSecret)
	token, err := at.SignedString(tokenSecretBuf)
	if err != nil {
		return "", models.NewInternalServerError(JwtSigningError)
	}

	return token, nil
}

func generateTokenPair(user *models.User) (string, string, error) {
	// Generate Access Token
	atSecret := viper.GetString(AccessTokenConfigSecretKeyName)
	accessToken, err := GenerateToken(user.ID, atSecret, AccessTokenExpiryMinutes)
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	rtSecret := viper.GetString(RefreshTokenConfigSecretKeyName)
	refreshToken, err := GenerateToken(user.ID, rtSecret, RefreshTokenExpiryMinutes)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}


func tokenStringToJwtToken(tokenString, tokenSecret string) (*jwt.Token, error) {
	cookieTokenSecretBuf := []byte(tokenSecret)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.NewInternalServerError(JwtSigningError)
		}
		return cookieTokenSecretBuf, nil
	}

	return jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, keyFunc)
}

func extractTokenClaims(jwtToken *jwt.Token) (jwt.StandardClaims, error) {
	// https://stackoverflow.com/questions/52460230/dgrijalva-jwt-go-can-cast-claims-to-mapclaims-but-not-standardclaims
	claims, ok := jwtToken.Claims.(*jwt.StandardClaims)
	fmt.Println("HERE: ", claims, ok)
	if !ok {
		return jwt.StandardClaims{}, models.NewUnauthorizedError(InvalidJwtToken)
	}

	return *claims, nil
}

func isTokenValid(tokenString, tokenSecret string) error {
	_, err := tokenStringToJwtToken(tokenString, tokenSecret)
	return err
}

func isTokenSignatureValid(tokenString, tokenSecret string) bool {
	token, err := tokenStringToJwtToken(tokenString, tokenSecret)
	if token == nil {
		return false
	}

	if token.Valid {
		return true
	}

	if ve, ok := err.(*jwt.ValidationError); ok && (ve.Errors&jwt.ValidationErrorMalformed == 0 && ve.Errors&jwt.ValidationErrorSignatureInvalid == 0 && ve.Errors&jwt.ValidationErrorUnverifiable == 0) {
		return true
	}

	return false
}

func IsTokenExpired(tokenString, tokenSecret string) bool {
	token, err := tokenStringToJwtToken(tokenString, tokenSecret)
	if token == nil {
		return false
	}

	if token.Valid {
		return false
	}

	if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
		return true
	}

	return false
}

func isTokenMalformed(tokenString, tokenSecret string) bool {
	token, err := tokenStringToJwtToken(tokenString, tokenSecret)
	if token == nil {
		return false
	}

	if token.Valid {
		return false
	}

	if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&jwt.ValidationErrorMalformed != 0 {
		return true
	}

	return false
}

func isTokenPremature(tokenString, tokenSecret string) bool {
	token, err := tokenStringToJwtToken(tokenString, tokenSecret)
	if token == nil {
		return false
	}

	if token.Valid {
		return false
	}

	if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
		return true
	}

	return false
}

// func Example(tokenString, tokenSecret string) bool {
// 	// var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"
//
// 	token, err := tokenStringToJwtToken(tokenString, tokenSecret)
// 	if token == nil {
// 		return false
// 	}
//
// 	if token.Valid {
// 		// fmt.Println("You look nice today")
// 	} else if ve, ok := err.(*jwt.ValidationError); ok {
// 		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
// 			// fmt.Println("That's not even a token")
// 		} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
// 			// Token is expired
// 			// fmt.Println("Timing is everything")
// 		} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
// 			// Token is not active yet
// 			// fmt.Println("Timing is everything")
// 		} else {
// 			// fmt.Println("Couldn't handle this token:", err)
// 		}
// 	} else {
// 		// fmt.Println("Couldn't handle this token:", err)
// 	}
// }
