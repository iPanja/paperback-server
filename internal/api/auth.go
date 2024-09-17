package api

import (
	"context"
	"errors"
	"fmt"
	"paperback-server/internal/db"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

var jwt_key = []byte("mysecret321")

const TokenDuration = 1 // in hours
type token string

type my_jwt struct {
	jwt.RegisteredClaims
	Account Account `json:"account" bson:"account"`
}

// Give the user a new access token (API key)
func IssueAccessToken(account Account) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		my_jwt{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "paperback-auth-server",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * TokenDuration)),
			},
			Account: account,
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
func ValidateToken(token string) (Account, error) {
	// Parse
	claims := my_jwt{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwt_key), nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		return Account{}, err
	}

	return claims.Account, nil
}

// Extract the account from the authorization header
// returns the account and true if the token is valid, false otherwise
func isAuthorized(authorization_header string) (Account, bool) {
	if len(authorization_header) < len("Bearer ") {
		return Account{}, false
	}

	// Check the token
	token := authorization_header[len("Bearer "):]
	if acc, err := ValidateToken(token); err != nil {
		fmt.Println("ERROR: ", err)
		return Account{}, false
	} else {
		return acc, true
	}
}

func LoginRequest(username string, password string) (string, error) {
	// Authenticate the user
	if username != "iPanja" {
		return "", errors.New("username not found")
	}

	// Find the user in the database
	client, err := db.GetClient()
	if err != nil {
		return "Internal server error, please try again later", err
	}
	defer client.Disconnect(context.TODO())

	var account Account
	filter := bson.D{{"username", username}}
	err = client.Database("paperback").Collection("accounts").FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		return "Internal server error, please try again later", err
	}

	// Issue token
	return IssueAccessToken(account)
}
