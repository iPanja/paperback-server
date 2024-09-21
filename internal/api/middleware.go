package api

import (
	"net/http"
	"paperback-server/internal/service"

	"github.com/labstack/echo/v4"
)

// EnforceLogin is a middleware to validate the access token supplied in the HTTP request.
//
// It will set the field "account" (models.Account) if valid.
func EnforceLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := service.ExtractToken(c.Request())
		if !ok {
			return c.String(BadRequest())
		}

		acc, tokenType, err := service.ValidateToken(token)
		if err != nil || tokenType != "access" {
			return c.String(http.StatusForbidden, "error: not authorized")
		}

		c.Set("account", acc)
		return next(c)
	}
}
