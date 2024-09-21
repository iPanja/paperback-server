// Package api is responsible for handling all HTTP requests
package api

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"paperback-server/internal/service"
)

// BadRequest returns a generic bad request error (400).
var BadRequest = func() (int, string) {
	return http.StatusBadRequest, "error: invalid request"
}

// InternalServerError returns a generic internal server error (500).
var InternalServerError = func() (int, string) {
	return http.StatusInternalServerError, "server error, please try again later"
}

// Login - /auth/login endpoint
//
// Attempts to log in the user given the credentials supplied in the request form.
func Login(c echo.Context) error {
	// Responses:
	//   200: {tokens: {"access_token" : "", "refresh_token": ""}
	//   500: authentication issue
	fmt.Println("received login request")
	c.Logger().Info("received login request")
	username, password := c.FormValue("username"), c.FormValue("password")

	accessToken, refreshToken, ok := service.LoginRequest(c, username, password)
	if !ok {
		c.Logger().Info("failed login attempt")
		return c.String(InternalServerError())
	}

	jwtBundle := map[string]string{
		"access":  accessToken,
		"refresh": refreshToken,
	}

	bytes, err := json.Marshal(jwtBundle)
	if err != nil {
		c.Logger().Error(err)
		return c.String(InternalServerError())
	}

	return c.String(http.StatusOK, fmt.Sprintf("{\"tokens\": \"%s\"}", bytes))
}

// RefreshTokens - /auth/refresh endpoint
//
// Issues a new access, refresh token to the user and invalidates the old refresh token if applicable.
func RefreshTokens(c echo.Context) error {
	// Responses:
	//   200: {tokens: {"access_token" : "", "refresh_token": ""}
	//   400: missing refresh token
	//   500: authentication issue
	fmt.Println("received refresh tokens request")
	suppliedRefreshToken, ok := service.ExtractToken(c.Request())
	if !ok {
		return c.String(BadRequest())
	}

	acc, ok := service.ValidateRefreshToken(suppliedRefreshToken)
	if !ok {
		c.Logger().Warn("failed to validate refresh token")
		return c.String(BadRequest())
	}

	// Issue a new access and refresh token; delete the old refresh token
	accessToken, refreshToken, ok := service.GenerateAndStoreTokenPair(acc)

	if !ok {
		c.Logger().Warn("failed to generate and store new tokens")
		return c.String(InternalServerError())
	}

	jwtBundle := map[string]string{
		"access":  accessToken,
		"refresh": refreshToken,
	}

	bytes, err := json.Marshal(jwtBundle)
	if err != nil {
		return c.String(InternalServerError())
	}
	return c.String(http.StatusOK, fmt.Sprintf("{\"tokens\": \"%s\"}", bytes))
}

func TestHashPassword(c echo.Context) error {
	// Responses:
	//   200: {"hash": ""}
	//   500: bcrypt error
	fmt.Println("received hash password request")
	password := c.FormValue("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(InternalServerError())
	}

	return c.String(http.StatusOK, fmt.Sprintf("{\"hash\": \"%s\"}", hash))
}
