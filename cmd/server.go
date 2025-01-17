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

	// Upload endpoints
	e.GET("/upload", api.UploadBook)

	e.Logger.Fatal(e.Start(":1323"))
}
