package db

import (
	"context"
	"fmt"
	"net/http"

	"paperback-server/internal/api"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connection_string = "mongodb://localhost:27017/"

func DBTest(c echo.Context) error {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connection_string).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {

			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	filter := bson.D{{"username", "iPanja"}}
	var result api.Account
	err = client.Database("paperback").Collection("accounts").FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	c.String(http.StatusOK, fmt.Sprintf("PINGED! %s\n", result))

	return nil
}
