package repos

import (
	"context"
	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type DisConnectable interface {
	CloseConnection(ctx context.Context) error
}

type MongodbRepo struct {
	Client *mongo.Client
}

func NewRepo(ctx context.Context) (*MongodbRepo, error) {
	client, err := db.NewMongoClient(ctx, config.GetEnv().MongoConnStr)
	if err != nil {
		return nil, err
	}
	return &MongodbRepo{Client: client}, nil
}

func (r *MongodbRepo) CloseConnection(ctx context.Context) error {
	if r.Client == nil {
		return nil
	}
	return r.Client.Disconnect(ctx)
}
