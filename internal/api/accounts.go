package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
}

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

	token, err := LoginRequest(username, password)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	}

	return c.String(http.StatusOK, fmt.Sprintf("{\"token\": \"%s\"}", token))
}

func TestToken(c echo.Context) error {
	// Test the token
	//
	// Test the token.
	//
	// Responses:
	//   200: testTokenResponse
	//   400: errorResponse
	//   500: errorResponse
	fmt.Println("received test token request")
	raw_auth := c.Request().Header.Get("Authorization")
	fmt.Printf("Authorization: <%s>\n", raw_auth)

	account, ok := isAuthorized(raw_auth)
	if !ok {
		return c.String(http.StatusForbidden, "error: not authorized")
	}

	return c.String(http.StatusOK, fmt.Sprintf("{\"account\": \"%s\"}", account))
}
