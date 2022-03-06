package db

import (
	"context"
	"github.com/joexzh/ThsConcept/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDbClient struct {
	*mongo.Client
	err error
}

var _client *mongoDbClient

func init() {
	client, err := newMongoClient(context.Background(), config.GetEnv().MongoConnStr)
	_client = &mongoDbClient{Client: client, err: err}
}

func GetMongoClient() (*mongo.Client, error) {
	return _client.Client, _client.err
}

func newMongoClient(ctx context.Context, connStr string) (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(connStr)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}
