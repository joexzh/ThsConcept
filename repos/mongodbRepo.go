package repos

import (
	"github.com/joexzh/ThsConcept/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongodbRepo struct {
	Client *mongo.Client
}

func NewMongoDbRepo() (*MongodbRepo, error) {
	client, err := db.GetMongoClient()
	if err != nil {
		return nil, err
	}
	return &MongodbRepo{Client: client}, nil
}
