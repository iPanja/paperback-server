package main

import (
	"net/http"

	"paperback-server/internal/api"
	"paperback-server/internal/db"

	"github.com/labstack/echo/v4"
)

const version = "1.0.0"

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\nVersion: "+version)
	})

	// Account endpoints
	e.POST("/login", api.Login)

	// Book endpoints

	// Author endpoints

	// Series endpoints

	// Collection endpoints

	// File endpoints

	// Device endpoints

	// TESTING
	e.GET("/test", db.DBTest)
	e.POST("/test_token", api.TestToken)

	e.Logger.Fatal(e.Start(":1323"))
}
