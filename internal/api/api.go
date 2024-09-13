package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UploadBook(c echo.Context) error {
	// Upload a book
	//
	// Upload a book to the server.
	//
	// Responses:
	//   200: bookResponse
	//   400: errorResponse
	//   500: errorResponse
	fmt.Println("received upload request")
	return c.String(http.StatusOK, "receiving upload request...")
}
