package db

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connection_string = "mongodb://localhost:27017/"

// Client exists to bundle an (echo) HTTP and database context together.
//
// The only way to access a database client/cursor is through this struct.
// While you technically can create a nil context, it is not advised.
type Client struct {
	Context echo.Context
}

// Create a new client and connect to the server
func getClient() (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connection_string).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c Client) Logger() echo.Logger {
	return c.Logger()
}

func (c Client) GetClient() *mongo.Client {
	client, err := getClient()
	if err != nil && c.Context != nil {
		c.Context.Logger().Error(err)
		c.Context.String(500, "error: failed to connect to the database")
	}

	return client
}
