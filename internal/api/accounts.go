package api

import (
	"context"
	"fmt"
	"net/http"
	"paperback-server/internal/db"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
}

func FetchAccount(username string) (Account, bool) {
	// Fetch the account from the database
	//
	// Fetch the account from the database.
	//
	// Responses:
	//   200: accountResponse
	//   400: errorResponse
	//   500: errorResponse

	// Get client instance
	client, err := db.GetClient()
	if err != nil {
		return Account{}, false
	}
	defer client.Disconnect(context.TODO())

	// Fetch the account
	var account Account
	filter := bson.D{{Key: "username", Value: username}}
	err = client.Database("paperback").Collection("accounts").FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		return Account{}, false
	}

	return account, true
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

	account, ok := IsAuthorized(raw_auth)
	if !ok {
		return c.String(http.StatusForbidden, "error: not authorized")
	}

	return c.String(http.StatusOK, fmt.Sprintf("{\"account\": \"%s\"}", account))
}
