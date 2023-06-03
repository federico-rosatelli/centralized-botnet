package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *appDB) Ping() error {
	return db.client.Ping(Ctx, nil)
}

func (db *appDB) Collection(typeCli int) *mongo.Collection {
	var collection mongo.Collection
	switch typeCli {
	case 0:
		collection = *db.zombies
	}
	return &collection
}
