package api

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwt_key = []byte("mysecret321")

const TokenDuration = 1 // in hours

// Give the user a new access token (API key)
func IssueAccessToken(username string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "paperback-auth-server",
			"sub": username,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour * TokenDuration).Unix(),
		})

	s, err := t.SignedString(jwt_key)

	if err != nil {
		return "", err
	}

	return s, nil
}

// Validate the token
// Returns the username if the token is valid
// Enforces expiration time
func ValidateToken(token string) (string, error) {
	// Parse
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwt_key), nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		return "", err
	}

	// Return the username
	if claims["sub"] == nil {
		return "", errors.New("username not found")
	}

	return claims["sub"].(string), err
}

func LoginRequest(username string, password string) (string, error) {
	// Authenticate the user
	if username != "iPanja" {
		return "", errors.New("username not found")
	}

	return IssueAccessToken(username)
}
