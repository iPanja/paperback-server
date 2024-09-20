package api

import (
	"fmt"
	"net/http"
	"paperback-server/internal/db"
	"paperback-server/internal/models"

	"github.com/labstack/echo/v4"
)

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
	rawAuth := c.Request().Header.Get("Authorization")
	fmt.Printf("Authorization: <%s>\n", rawAuth)

	account, ok := IsAuthorized(rawAuth, "") // allow any type of token
	if !ok {
		return c.String(http.StatusForbidden, "error: not authorized")
	}

	return c.String(http.StatusOK, fmt.Sprintf("{\"account\": \"%s\"}", account))
}

func ViewAccount(c echo.Context) error {
	acc := c.Get("account").(models.Account)
	fmt.Println("acc: ", acc)
	return c.String(http.StatusOK, fmt.Sprintf("View account, %s. You are authenticated!", acc.Username))
}

func ViewAccountByUsername(c echo.Context) error {
	username := c.Param("username")
	ac := models.AccountClient{Client: &db.Client{Context: c}}
	if account, ok := ac.FetchAccountByUsername(username); !ok {
		return c.String(http.StatusNotFound, "Account not found")
	} else {
		return c.String(http.StatusOK, fmt.Sprintf("View account, %s.\n%+v", account.Username, account))
	}
}
