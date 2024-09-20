package models

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GenericClient interface {
	GetClient() *mongo.Client
	Logger() echo.Logger
}

func InsertOne(client GenericClient, collection string, doc interface{}) (primitive.ObjectID, error) {
	result, err := client.GetClient().Database("paperback").Collection(collection).InsertOne(context.Background(), doc)
	if err != nil {
		if client.Logger() != nil {
			client.Logger().Warn(fmt.Sprintf("Error inserting document into paperback database: collection: %s, document: %+v, error: %s", collection, doc, err.Error()))
		}

		return primitive.ObjectID{}, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func DeleteOne(client GenericClient, collection string, filter bson.D) error {
	_, err := client.GetClient().Database("paperback").Collection(collection).DeleteOne(context.TODO(), filter)
	if err != nil {
		if client.Logger() != nil {
			client.Logger().Warn(fmt.Sprintf("Error deleting document from paperback database: collection: %s, filter: %+v, error: %s", collection, filter, err.Error()))
		}

		return err
	}

	return nil
}
