package main

import (
	"net/http"

	"paperback-server/internal/api"

	"github.com/labstack/echo/v4"
)

const version = "1.0.0"

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\nVersion: "+version)
	})

	// Account endpoints
	e.GET("/account", api.ViewAccount, api.EnforceLogin)

	// Authentication endpoints
	e.POST("/login", api.Login)
	e.POST("/hash", api.HashPassword)

	// Book endpoints

	// Author endpoints

	// Series endpoints

	// Collection endpoints

	// File endpoints

	// Device endpoints

	// TESTING
	e.POST("/test_token", api.TestToken)

	e.Logger.Fatal(e.Start(":1323"))
}
