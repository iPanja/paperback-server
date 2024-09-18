package api

import (
	"context"
	"fmt"
	"net/http"
	"paperback-server/internal/db"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var jwt_key = []byte("mysecret321")

const TokenDuration = 1 // in hours
type token string

type my_jwt struct {
	jwt.RegisteredClaims
	Account Account `json:"account" bson:"account"`
}

// Give the user a new access token (API key)
func issueAccessToken(account Account) (string, error) {
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

// Validate the token.
// Returns the username if the token is valid.
// Enforces expiration time.
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

// Extract the token from the authorization header.
// Returns the account and true if the token is valid, false otherwise
func IsAuthorized(authorization_header string) (Account, bool) {
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

// Returns a jwt token, true if the user is authenticated
// Returns false if the user is not authenticated
func LoginRequest(username string, password string) (string, bool) {
	// Authenticate the user
	// Find the user in the database
	type unsafeAccount struct {
		ID       primitive.ObjectID `bson:"_id" json:"id"`
		Username string             `bson:"username" json:"username"`
		Password string             `bson:"password" json:"password"`
	}

	// Get client instance
	client, err := db.GetClient()
	if err != nil {
		fmt.Println("Failed to get client ", err)
		return "", false
	}
	defer client.Disconnect(context.TODO())

	// Find account
	var account unsafeAccount
	filter := bson.D{{Key: "username", Value: username}}
	err = client.Database("paperback").Collection("accounts").FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		fmt.Println("Failed to find account ", err)
		return "", false
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		fmt.Println("Failed to compare passwords ", err)
		return "", false
	}

	// Issue token
	token, err := issueAccessToken(Account{account.ID, account.Username})
	return token, err == nil
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
	username, password := c.FormValue("username"), c.FormValue("password")

	token, ok := LoginRequest(username, password)
	if !ok {
		return c.String(http.StatusInternalServerError, "Internal server error, please try again later")
	}

	return c.String(http.StatusOK, fmt.Sprintf("{\"token\": \"%s\"}", token))
}

func HashPassword(c echo.Context) error {
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

// Middlweware to check the token
func EnforceLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		raw_auth := c.Request().Header.Get("Authorization")
		acc, ok := IsAuthorized(raw_auth)
		if !ok {
			return c.String(http.StatusForbidden, "error: not authorized")
		}

		c.Set("account", acc)
		return next(c)
	}
}
