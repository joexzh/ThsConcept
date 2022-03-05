package repos

import (
	"context"
	"errors"
	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/model"
	"github.com/joexzh/ThsConcept/realtime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RealtimeRepo struct {
	*MongodbRepo
	collRealtime *mongo.Collection
	collConcept  *mongo.Collection
}

func NewRealtimeRepo(ctx context.Context) (*RealtimeRepo, error) {
	r, err := NewRepo(ctx)
	if err != nil {
		return nil, err
	}
	repo := RealtimeRepo{
		MongodbRepo:  r,
		collRealtime: r.Client.Database(config.Db).Collection(config.CollRealtime),
		collConcept:  r.Client.Database(config.Db).Collection(config.CollConcept),
	}
	return &repo, nil
}

func (r *RealtimeRepo) Query(ctx context.Context, queryDoc *bson.D, opts ...*options.FindOptions) ([]realtime.SavedMessage, error) {
	cursor, err := r.collRealtime.Find(ctx, queryDoc, opts...)
	if err != nil {
		return nil, err
	}

	var list []realtime.SavedMessage
	if err = cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *RealtimeRepo) QuerySaveMessageDesc(ctx context.Context, userId string) ([]realtime.SavedMessage, error) {
	queryDoc := bson.D{{"userId", userId}}
	opt := options.Find().SetSort(bson.D{{"message.ctime", -1}})
	opt.SetCollation(&options.Collation{Locale: "en_US", NumericOrdering: true})
	return r.Query(ctx, &queryDoc, opt)
}

func (r *RealtimeRepo) SaveMessage(ctx context.Context, userId string, message *realtime.Message) error {
	msg := realtime.SavedMessage{
		UserId:  userId,
		Message: *message,
	}
	ret, err := r.collRealtime.InsertOne(ctx, &msg)
	if err != nil {
		return err
	}
	if ret.InsertedID == nil {
		return errors.New("not inserted")
	}
	return nil
}

func (r *RealtimeRepo) DelSaveMessage(ctx context.Context, userId string, objId string) error {
	id, err := primitive.ObjectIDFromHex(objId)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", id}, {"userId", userId}}
	ret, err := r.collRealtime.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if ret.DeletedCount < 1 {
		return errors.New("not deleted")
	}
	return nil
}

func (r *RealtimeRepo) GetAllConceptNames(ctx context.Context) ([]string, error) {
	projection := bson.D{{"_id", 0}, {"conceptName", 1}}
	cursor, err := r.collConcept.Find(ctx, bson.D{}, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	var concepts []model.Concept
	if err = cursor.All(ctx, &concepts); err != nil {
		return nil, err
	}
	names := make([]string, 0, 300)
	for _, concept := range concepts {
		names = append(names, concept.ConceptName)
	}
	return names, nil
}
