package api

import (
	"context"
	"fmt"
	"net/http"
	"paperback-server/internal/db"
	"paperback-server/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// go::embed jwt.key
var jwt_key []byte

const TokenDuration = 1 // in hours

type my_jwt struct {
	jwt.RegisteredClaims
	Account models.Account `json:"account" bson:"account"`
}

// Give the user a new access token (API key)
func issueAccessToken(account models.Account) (string, error) {
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
func ValidateToken(token string) (models.Account, error) {
	// Parse
	claims := my_jwt{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwt_key), nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		return models.Account{}, err
	}

	// TODO: Check if the account is still valid and get updated account information straight from the DB
	return claims.Account, nil
}

// Extract the token from the authorization header.
// Returns the account and true if the token is valid, false otherwise
func IsAuthorized(authorization_header string) (models.Account, bool) {
	if len(authorization_header) < len("Bearer ") {
		return models.Account{}, false
	}

	// Check the token
	token := authorization_header[len("Bearer "):]
	if acc, err := ValidateToken(token); err != nil {
		fmt.Println("ERROR: ", err)
		return models.Account{}, false
	} else {
		return acc, true
	}
}

// Returns a jwt token, true if the user is authenticated
// Returns false if the user is not authenticated
func LoginRequest(ctx echo.Context, username string, password string) (string, bool) {
	// Authenticate the user
	// Find the user in the database
	type unsafeAccount struct {
		ID       primitive.ObjectID `bson:"_id" json:"id"`
		Username string             `bson:"username" json:"username"`
		Password string             `bson:"password" json:"password"`
	}

	// Get client instance
	client := db.Client{Context: ctx}.GetClient()
	defer client.Disconnect(context.TODO())

	// Find account
	var account unsafeAccount
	filter := bson.D{{Key: "username", Value: username}}
	err := client.Database("paperback").Collection("accounts").FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		fmt.Println("Failed to find account ", err)
		return "", false
	}

	// Compare passwords
	if !models.CheckPassword(account.Password, password) {
		return "", false
	}

	// Issue token
	token, err := issueAccessToken(models.Account{ID: account.ID, Username: account.Username})
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

	token, ok := LoginRequest(c, username, password)
	if !ok {
		return c.String(http.StatusInternalServerError, "Internal server error, please try again later")
	}

	return c.String(http.StatusOK, fmt.Sprintf("{\"token\": \"%s\"}", token))
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
