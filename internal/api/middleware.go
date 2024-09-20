package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Middlweware to check the token
func EnforceLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		raw_auth := c.Request().Header.Get("Authorization")
		acc, ok := IsAuthorized(raw_auth, "access") // only allow access tokens
		if !ok {
			return c.String(http.StatusForbidden, "error: not authorized")
		}

		c.Set("account", acc)
		return next(c)
	}
}
