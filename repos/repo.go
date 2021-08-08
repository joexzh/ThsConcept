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

type Repo struct {
	Client *mongo.Client
}

func NewRepo(ctx context.Context) (*Repo, error) {
	client, err := db.NewMongoClient(ctx, config.GetEnv().MongoConnStr)
	if err != nil {
		return nil, err
	}
	return &Repo{Client: client}, nil
}

func (r *Repo) CloseConnection(ctx context.Context) error {
	if r.Client == nil {
		return nil
	}
	return r.Client.Disconnect(ctx)
}
