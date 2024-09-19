package models

import (
	"context"
	"fmt"
	"paperback-server/internal/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AccountClient struct {
	*db.Client
}

type Account struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
}

func (c AccountClient) FetchAccountByUsername(username string) (Account, bool) {
	client := c.GetClient()

	// Fetch the account
	var account Account
	filter := bson.D{{Key: "username", Value: username}}
	err := client.Database("paperback").Collection("accounts").FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		return Account{}, false
	}

	return account, true
}

func (a *Account) HashPassword(plaintext string) (string, error) {
	if len(plaintext) == 0 {
		return "", nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPassword(hash string, plaintext string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(hash))
	if err != nil {
		fmt.Println("Failed to compare passwords ", err)
		return false
	}

	return true
}
