package repos

import (
	"context"
	"fmt"
	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConceptRepo struct {
	*MongodbRepo
}

func NewConceptRepo(ctx context.Context) (*ConceptRepo, error) {
	r, err := NewRepo(ctx)
	if err != nil {
		return nil, err
	}
	return &ConceptRepo{r}, nil
}

func (r *ConceptRepo) Query(ctx context.Context, queryStr string, opts ...*options.FindOptions) ([]model.Concept, error) {

	conceptsColl := r.Client.Database(config.Db).Collection(config.CollConcept)
	var queryDoc bson.D
	err := bson.UnmarshalExtJSON([]byte(queryStr), false, &queryDoc)
	if err != nil {
		return nil, err
	}
	cursor, err := conceptsColl.Find(ctx, queryDoc, opts...)
	if err != nil {
		return nil, err
	}

	var concepts []model.Concept
	if err = cursor.All(ctx, &concepts); err != nil {
		return nil, err
	}
	return concepts, nil
}

func (r *ConceptRepo) QueryByConceptName(ctx context.Context, conceptName string) ([]model.Concept, error) {
	queryStr := fmt.Sprintf(`{"conceptName":"%v"}`, conceptName)
	return r.Query(ctx, queryStr)
}

func (r *ConceptRepo) QueryByConceptNameRex(ctx context.Context, conceptName string) ([]model.Concept, error) {
	queryStr := fmt.Sprintf(`{"conceptName": {"$regex": "%v"}}`, conceptName)
	return r.Query(ctx, queryStr)
}

func (r *ConceptRepo) QueryByConceptId(ctx context.Context, conceptId string) (*model.Concept, error) {
	queryStr := fmt.Sprintf(`{"conceptId":"%v"}`, conceptId)
	concepts, err := r.Query(ctx, queryStr)
	if err != nil {
		return nil, err
	}
	if len(concepts) > 0 {
		return &concepts[0], nil
	} else {
		return nil, nil
	}
}

func (r *ConceptRepo) Update(ctx context.Context, concepts ...model.Concept) (int64, error) {
	conceptsColl := r.Client.Database(config.Db).Collection(config.CollConcept)
	opts := options.Replace().SetUpsert(true)

	var updated int64

	for _, concept := range concepts {
		// compare old and new doc is the same because we need ACTUAL lastModified.
		var oldConcept model.Concept
		single := conceptsColl.FindOne(ctx, bson.M{"conceptId": concept.ConceptId})
		if single.Err() == nil { // match an old concept
			if err := single.Decode(&oldConcept); err != nil {
				return 0, err
			}
			if oldConcept.Compare(&concept) {
				continue
			}
		}
		concept.SetLastModifiedNow()
		ret, err := conceptsColl.ReplaceOne(ctx, bson.M{"conceptId": concept.ConceptId}, concept, opts)
		if err != nil {
			return 0, err
		}
		updated += ret.ModifiedCount + ret.UpsertedCount
	}
	return updated, nil
}

func (r *ConceptRepo) DeleteUnMatch(ctx context.Context, concepts ...model.Concept) (int64, error) {
	conceptsColl := r.Client.Database(config.Db).Collection(config.CollConcept)

	cids := make([]string, 0, len(concepts))
	for _, concept := range concepts {
		cids = append(cids, concept.ConceptId)
	}

	ret, err := conceptsColl.DeleteMany(ctx, bson.M{"conceptId": bson.M{"$nin": cids}})
	if err != nil {
		return 0, err
	}
	return ret.DeletedCount, nil
}

func (r *ConceptRepo) UpdateConceptColl(ctx context.Context, concepts ...model.Concept) (int64, int64, error) {
	deletedNum, err := r.DeleteUnMatch(ctx, concepts...)
	if err != nil {
		return 0, 0, err
	}

	matched, err := r.Update(ctx, concepts...)
	if err != nil {
		return 0, 0, err
	}
	return deletedNum, matched, nil
}

// stockConcept collection ↓

func (r *ConceptRepo) QueryScDesc(ctx context.Context, stockName string, conceptNameRex string, limit int) ([]model.StockConcept, error) {
	coll := r.Client.Database(config.Db).Collection(config.CollStockConcept)
	if limit < 1 || limit > 1000 {
		limit = 1000
	}

	match := bson.D{}
	if stockName != "" {
		match = append(match, bson.E{"stockName", stockName})
	}
	if conceptNameRex != "" {
		match = append(match, bson.E{"conceptName", bson.M{"$regex": conceptNameRex}})
	}
	queryStage := bson.D{{"$match", match}}
	sortStage := bson.D{{"$sort", bson.M{"lastModified": -1}}}
	limitStage := bson.D{{"$limit", limit}}

	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{queryStage, sortStage, limitStage})
	if err != nil {
		return nil, err
	}
	var scs []model.StockConcept
	if err = cursor.All(ctx, &scs); err != nil {
		return nil, err
	}
	return scs, nil
}

func (r *ConceptRepo) SortLatestStockConcept(ctx context.Context, limit int) ([]model.StockConcept, error) {
	coll := r.Client.Database(config.Db).Collection(config.CollStockConcept)

	if limit == 0 {
		limit = 1000
	} else if limit > 1000 {
		limit = 1000
	}

	sortStage := bson.D{{"$sort", bson.M{"lastModified": -1}}}
	limitStage := bson.D{{"$limit", limit}}
	cursor, err := coll.Aggregate(ctx, mongo.Pipeline{sortStage, limitStage})
	if err != nil {
		return nil, err
	}
	var scs []model.StockConcept
	if err = cursor.All(ctx, &scs); err != nil {
		return nil, err
	}
	return scs, nil
}

func (r *ConceptRepo) DelUnExistStockConcept(ctx context.Context, scs ...model.StockConcept) (int64, error) {
	conceptsColl := r.Client.Database(config.Db).Collection(config.CollStockConcept)

	ids := make([]string, 0, len(scs))
	for _, sc := range scs {
		ids = append(ids, sc.Id)
	}

	ret, err := conceptsColl.DeleteMany(ctx, bson.M{"_id": bson.M{"$nin": ids}})
	if err != nil {
		return 0, err
	}
	return ret.DeletedCount, nil
}

func (r *ConceptRepo) UpsertStockConcept(ctx context.Context, scs ...model.StockConcept) (int64, error) {
	conceptsColl := r.Client.Database(config.Db).Collection(config.CollStockConcept)

	opts := options.Replace().SetUpsert(true)
	var updated int64

	for _, sc := range scs {
		singleRet := conceptsColl.FindOne(ctx, bson.M{"_id": sc.Id})

		if singleRet.Err() == nil { // found
			var oldSc model.StockConcept
			if err := singleRet.Decode(&oldSc); err != nil {
				return 0, err
			}
			if sc.Compare(&oldSc) {
				continue
			}
		}

		ret, err := conceptsColl.ReplaceOne(ctx, bson.M{"_id": sc.Id}, sc, opts)
		if err != nil {
			return 0, err
		}
		updated += ret.ModifiedCount + ret.UpsertedCount
	}
	return updated, nil
}

func (r *ConceptRepo) UpdateStockConcept(ctx context.Context, concepts ...model.Concept) (int64, int64, error) {
	scsLen := 0
	for _, concept := range concepts {
		scsLen += len(concept.Stocks)
	}

	scs := make([]model.StockConcept, 0, scsLen)
	for _, concept := range concepts {
		for _, stock := range concept.Stocks {
			sc := model.NewStockConcept(stock, concept.ConceptId, concept.ConceptName)
			scs = append(scs, *sc)
		}
	}

	delCount, err := r.DelUnExistStockConcept(ctx, scs...)
	if err != nil {
		return 0, 0, nil
	}
	upsertCount, err := r.UpsertStockConcept(ctx, scs...)
	if err != nil {
		return 0, 0, nil
	}
	return delCount, upsertCount, nil
}

// stockConcept collection ↑
