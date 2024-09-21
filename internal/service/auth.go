package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"paperback-server/internal/db"
	"paperback-server/internal/models"
	"time"
)

// go::embed jwt.key
var JwtKey []byte

const TokenDuration = 1 // in hours
type myJwt struct {
	jwt.RegisteredClaims
	Account   models.Account `json:"account" bson:"account"`
	TokenType string         `json:"token_type" bson:"token_type"`
}

// IssueAccessToken creates a signed JWT with the account information embedded in it.
// Additionally, it specifies a string field denoting the token type (access, refresh).
func IssueAccessToken(account models.Account, tokenType string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		myJwt{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "paperback-auth-server",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * TokenDuration)),
			},
			Account:   account,
			TokenType: tokenType,
		})

	s, err := t.SignedString(JwtKey)

	if err != nil {
		return "", err
	}

	return s, nil
}

// ExtractToken returns the bearer (JWT) token inside the HTTP request's authorization header
func ExtractToken(r *http.Request) (string, bool) {
	authorizationHeader := r.Header.Get("Authorization")
	if len(authorizationHeader) < len("Bearer ") {
		return "", false
	}

	return authorizationHeader[len("Bearer "):], true
}

// ValidateToken will return the account embedded in the JWT if valid.
// This function does enforce the token's expiration time.
func ValidateToken(token string) (models.Account, string, error) {
	// Parse
	claims := myJwt{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		return models.Account{}, "", err
	}

	return claims.Account, claims.TokenType, nil
}

// GenerateAndStoreTokenPair will generate an access and refresh token.
// It will also store the refresh token inside the database.
func GenerateAndStoreTokenPair(account models.Account) (string, string, bool) {
	// Issue token
	accessToken, errA := IssueAccessToken(models.Account{ID: account.ID, Username: account.Username}, "access")
	refreshToken, errB := IssueAccessToken(models.Account{ID: account.ID, Username: account.Username}, "refresh")
	if errA != nil || errB != nil {
		return "", "", false
	}

	// Store the refresh token in the database, delete the old one
	ac := models.AccountClient{Client: &db.Client{Context: nil}}
	ac.StoreRefreshToken(account, refreshToken)

	return accessToken, refreshToken, true
}

func ValidateRefreshToken(token string) (models.Account, bool) {
	// We can ensure the JWT is valid
	// i.e. valid signature, not expired, etc
	acc, tokenType, err := ValidateToken(token)
	if err != nil || tokenType != "refresh" {
		return models.Account{}, false
	}

	// Make sure it hasn't been revoked
	ac := models.AccountClient{Client: &db.Client{}}
	if rt, ok := ac.FetchRefreshToken(acc); !ok || rt != token {
		return models.Account{}, false
	}

	return acc, true
}

// LoginRequest will return a new access and refresh token for the account if the credentials are valid.
// The refresh token is also stored inside the database.
func LoginRequest(ctx echo.Context, username string, password string) (string, string, bool) {
	// Authenticate the user
	// Find the user in the database
	type unsafeAccount struct {
		models.Account `bson:",inline"`
		Password       string `bson:"password" json:"password"`
	}

	// Get client instance
	client := db.Client{Context: ctx}.GetClient()
	defer client.Disconnect(context.TODO())

	// Find account
	var account unsafeAccount
	filter := bson.D{{Key: "username", Value: username}}
	err := client.Database("paperback").Collection("accounts").FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		ctx.Logger().Info("Failed to find account ", err)
		return "", "", false
	}

	// Compare passwords
	if !models.CheckPassword(account.Password, password) {
		ctx.Logger().Info("Failed to authenticate user")
		return "", "", false
	}

	// Issue token
	accessToken, refreshToken, ok := GenerateAndStoreTokenPair(account.Account)

	return accessToken, refreshToken, ok
}
