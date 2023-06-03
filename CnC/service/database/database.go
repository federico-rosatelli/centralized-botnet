package database

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type AppDatabase interface {
	// GetName() (string, error)
	// SetName(name string) error

	Ping() error
	Collection(typeCli int) *mongo.Collection
}

type appDB struct {
	client  *mongo.Client
	zombies *mongo.Collection
}

var Ctx = context.TODO()

func InitDatabase(client *mongo.Client) (AppDatabase, error) {
	collectionUsers := client.Database("botnet").Collection("zombies")
	if collectionUsers == nil {
		return nil, errors.New("Error Creating users Collection")
	}
	return &appDB{
		client:  client,
		zombies: collectionUsers,
	}, nil
}
