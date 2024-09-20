package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"paperback-server/internal/db"
)

type Book struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	AccountID      primitive.ObjectID `bson:"accountId" json:"accountId"`
	Title          string             `bson:"title" json:"title"`
	Author         Author             `bson:"author" json:"author"`
	Series         Series             `bson:"series" json:"series"`
	Collections    []Collection       `bson:"collections" json:"collections"`
	SeriesNumber   int                `bson:"seriesNumber" json:"seriesNumber"`
	CoverImagePath string             `bson:"coverImagePath" json:"coverImagePath"`
}

type BookClient struct {
	*db.Client
}

func (c *BookClient) GetBookByID(id primitive.ObjectID) (Book, error) {
	client := c.GetClient()

	var book Book
	filter := bson.M{"_id": id}
	err := client.Database("paperback").Collection("books").FindOne(context.TODO(), filter).Decode(&book)

	if err != nil {
		return Book{}, err
	}

	return book, nil
}

func (c *BookClient) CreateBook(book Book) (Book, error) {
	client := c.GetClient()

	if book.AccountID.IsZero() {
		return Book{}, errors.New("accountId is required")
	}

	result, err := client.Database("paperback").Collection("books").InsertOne(context.TODO(), book)
	if err != nil {
		c.Client.Logger().Error("Failed to create book: ", err)
		return Book{}, err
	}

	book.ID = result.InsertedID.(primitive.ObjectID)

	return book, nil
}
