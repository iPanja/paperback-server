package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"paperback-server/internal/db"
	"paperback-server/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// go::embed jwt.key
var JwtKey []byte

const TokenDuration = 1 // in hours

type myJwt struct {
	jwt.RegisteredClaims
	Account   models.Account `json:"account" bson:"account"`
	TokenType string         `json:"token_type" bson:"token_type"`
}

// Give the user a new access token (API key)
// token_type is either "access" or "refresh"
func issueAccessToken(account models.Account, token_type string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		myJwt{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "paperback-auth-server",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * TokenDuration)),
			},
			Account:   account,
			TokenType: token_type,
		})

	s, err := t.SignedString(JwtKey)

	if err != nil {
		return "", err
	}

	return s, nil
}

func generateAndStoreTokenPair(account models.Account) (string, string, bool) {
	// Issue token
	accessToken, errA := issueAccessToken(models.Account{ID: account.ID, Username: account.Username}, "access")
	refreshToken, errB := issueAccessToken(models.Account{ID: account.ID, Username: account.Username}, "refresh")
	if errA != nil || errB != nil {
		return "", "", false
	}

	// Stpre the refresh token in the database, delete the old one
	ac := models.AccountClient{Client: &db.Client{Context: nil}}
	ac.StoreRefreshToken(account, refreshToken)

	return accessToken, refreshToken, true
}

// Validate the token.
// Returns the account, token_type if valid; error otherwise
// Enforces expiration time.
func ValidateToken(token string) (models.Account, string, error) {
	// Parse
	claims := myJwt{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		return models.Account{}, "", err
	}

	// TODO: Check if the account is still valid and get updated account information straight from the DB
	return claims.Account, claims.TokenType, nil
}

// Extract the token from the authorization header.
// Returns the account and true if the token is valid, false otherwise
// Enforces token type if non-empty
func IsAuthorized(authorizationHeader string, token_type string) (models.Account, bool) {
	if len(authorizationHeader) < len("Bearer ") {
		return models.Account{}, false
	}

	// Check the token
	token := authorizationHeader[len("Bearer "):]
	if acc, tokenType, err := ValidateToken(token); err != nil && (tokenType != "" && tokenType != "access") {
		fmt.Println("ERROR: ", err)
		return models.Account{}, false
	} else {
		return acc, true
	}
}

// Returns a jwt token, true if the user is authenticated
// Returns false if the user is not authenticated
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
	accessToken, refreshToken, ok := generateAndStoreTokenPair(account.Account)

	return accessToken, refreshToken, ok
}

// API endpoint for logging in
// Returns a token if successful
// Endpoint: POST /login
func Login(c echo.Context) error {
	// Login to the server
	//
	// Login to the server.
	//
	// Responses:
	//   200: loginResponse
	//   400: errorResponse
	//   500: errorResponse
	fmt.Println("received login request")
	c.Logger().Info("received login request")
	username, password := c.FormValue("username"), c.FormValue("password")

	accessToken, refreshToken, ok := LoginRequest(c, username, password)
	if !ok {
		c.Logger().Info("failed login attempt")
		return c.String(http.StatusInternalServerError, "Internal server error, please try again later")
	}

	jwtBundle := map[string]string{
		"access":  accessToken,
		"refresh": refreshToken,
	}

	bytes, err := json.Marshal(jwtBundle)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal server error, please try again later")
	}
	return c.String(http.StatusOK, fmt.Sprintf("{\"tokens\": \"%s\"}", bytes))
}

func RefreshTokens(c echo.Context) error {
	// Refresh the tokens
	//
	// Refresh the tokens.
	//
	// Responses:
	//   200: refreshTokensResponse
	//   400: errorResponse
	//   500: errorResponse
	fmt.Println("received refresh tokens request")
	authorizationHeader := c.Request().Header.Get("Authorization")
	if len(authorizationHeader) < len("Bearer ") {
		return c.String(http.StatusBadRequest, "error: invalid request, missing refresh token")
	}

	// Validate the token
	suppliedRefreshToken := authorizationHeader[len("Bearer "):]
	acc, tokenType, err := ValidateToken(suppliedRefreshToken)
	if err != nil || tokenType != "refresh" {
		c.Logger().Info("failed to validate token: ", suppliedRefreshToken, err)
		return c.String(http.StatusForbidden, "error: not authorized")
	}

	// Compare supplied refresh token against entry in DB
	ac := models.AccountClient{Client: &db.Client{Context: c}}
	if rt, ok := ac.FetchRefreshToken(acc); !ok || rt != suppliedRefreshToken {
		c.Logger().Info("failed to validate token, does not match refresh token in DB: ", suppliedRefreshToken, err)
		return c.String(http.StatusForbidden, "error: not authorized")
	}

	// Issue a new access, refresh token and delete the old refresh token
	if accessToken, refreshToken, ok := generateAndStoreTokenPair(acc); !ok {
		c.Logger().Warn("failed to generate and store new tokens")
		return c.String(http.StatusInternalServerError, "Internal server error, please try again later")
	} else {
		jwtBundle := map[string]string{
			"access":  accessToken,
			"refresh": refreshToken,
		}

		bytes, err := json.Marshal(jwtBundle)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal server error, please try again later")
		}
		return c.String(http.StatusOK, fmt.Sprintf("{\"tokens\": \"%s\"}", bytes))
	}

}

func TestHashPassword(c echo.Context) error {
	// Hash a password
	//
	// Hash a password.
	//
	// Responses:
	//   200: hashPasswordResponse
	//   400: errorResponse
	//   500: errorResponse
	fmt.Println("received hash password request")
	password := c.FormValue("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error, please try again later")
	}

	return c.String(http.StatusOK, fmt.Sprintf("{\"hash\": \"%s\"}", hash))
}
