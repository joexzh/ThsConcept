package repos

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/joexzh/dbh"

	"github.com/joexzh/ThsConcept/db"
	"github.com/joexzh/ThsConcept/model"
	"github.com/joexzh/ThsConcept/tmpl"
	"github.com/pkg/errors"
)

type StockMarketRepo struct {
	DB   *sql.DB
	Name string
}

func NewStockMarketRepo(db *sql.DB) *StockMarketRepo {
	return &StockMarketRepo{DB: db, Name: "StockMarketRepo"}
}

func (repo *StockMarketRepo) ZdtListDesc(ctx context.Context, start time.Time, limit int) ([]*model.ZDTHistory, error) {
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	list, err := dbh.QueryContext[*model.ZDTHistory](repo.DB, ctx, tmpl.SelectZdt,
		start, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	// reverse
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	return list, nil
}

func (repo *StockMarketRepo) InsertZdtList(ctx context.Context, list []*model.ZDTHistory) (int64, error) {
	if len(list) < 1 {
		return 0, nil
	}

	conn, err := repo.DB.Conn(ctx)
	if err != nil {
		return 0, errors.Wrap(err, repo.Name)
	}
	defer conn.Close()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, errors.Wrap(err, repo.Name)
	}
	defer tx.Rollback()

	if err != nil {
		return 0, errors.Wrap(err, repo.Name)
	}

	total, err := dbh.BulkInsertContext(tx, ctx, 1000, list...)
	if err != nil {
		return 0, errors.Wrap(err, repo.Name)
	}
	if err = tx.Commit(); err != nil {
		return 0, errors.Wrap(err, repo.Name)
	}
	return total, nil
}

func scsToConceptMap(scs []*model.ConceptStockView) map[string]*model.Concept {
	conceptMap := make(map[string]*model.Concept)
	for _, sc := range scs {
		concept, ok := conceptMap[sc.ConceptId]
		if !ok {
			concept = &model.Concept{
				Id:        sc.ConceptId,
				Name:      sc.ConceptName,
				PlateId:   sc.ConceptPlateId,
				Define:    sc.ConceptDefine,
				UpdatedAt: sc.ConceptUpdatedAt,
				Stocks:    make([]*model.ConceptStock, 0),
			}
			conceptMap[sc.ConceptId] = concept
		}
		stock := &model.ConceptStock{
			StockCode:   sc.StockCode,
			StockName:   sc.StockName,
			Description: sc.Description,
			UpdatedAt:   sc.UpdatedAt,
			ConceptId:   sc.ConceptId,
		}
		concept.Stocks = append(concept.Stocks, stock)
	}
	return conceptMap
}

type UpdateConceptResult struct {
	ConceptConceptInserted int64
	ConceptConceptUpdated  int64
	ConceptConceptDeleted  int64
	ConceptStockInserted   int64
	ConceptStockUpdated    int64
	ConceptStockDeleted    int64
}

func (repo *StockMarketRepo) UpdateConcept(ctx context.Context, newcs ...*model.Concept) (UpdateConceptResult, error) {
	result := UpdateConceptResult{}

	conn, err := repo.DB.Conn(ctx)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	defer conn.Close()
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	defer tx.Rollback()

	// query old concepts and stocks from db for update
	oldscs, err := dbh.QueryContext[*model.ConceptStockView](repo.DB, ctx, tmpl.SelectFromConceptStockView+" for update")
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	oldcsMap := scsToConceptMap(oldscs)

	// 1. delete un-exist concepts
	cids := make([]string, 0, len(newcs))
	for _, c := range newcs {
		cids = append(cids, c.Id)
	}
	listSql, vals := db.ArgList(cids...)
	ret, err := tx.Exec("DELETE FROM concept_concept WHERE id NOT IN "+listSql, vals...)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	ra, _ := ret.RowsAffected()
	result.ConceptConceptDeleted = ra

	insertConceptStmt, err := tx.Prepare(tmpl.InsertConcept)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	updateConceptStmt, err := tx.Prepare(tmpl.UpdateConcept)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	insertConceptStockStmt, err := tx.Prepare(tmpl.InsertConceptStock)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	updateConceptStockStmt, err := tx.Prepare(tmpl.UpdateConceptStock)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}

	for _, newc := range newcs {
		// 2. update concept_concept
		oldc, ok := oldcsMap[newc.Id]
		if !ok {
			// insert
			_, err = insertConceptStmt.Exec(newc.Args()...)
			if err != nil {
				return result, errors.Wrap(err, repo.Name)
			}
			result.ConceptConceptInserted++
		} else {
			if newc.Define != "" && !newc.Cmp(oldc) {
				// update
				_, err = updateConceptStmt.Exec(newc.Name, newc.PlateId, newc.Define, newc.UpdatedAt, newc.Id)
				result.ConceptConceptUpdated++
			}
		}

		// 3. delete un-exist concept_stock
		codes := make([]string, len(newc.Stocks))
		for i, stock := range newc.Stocks {
			codes[i] = stock.StockCode
		}
		listSql, vals = db.ArgList(codes...)
		vals = append(vals, newc.Id)
		ret, err = tx.Exec(fmt.Sprintf("DELETE FROM concept_stock WHERE stock_code NOT IN %s and concept_id=?", listSql), vals...)
		if err != nil {
			return result, errors.Wrap(err, repo.Name)
		}
		ra, _ = ret.RowsAffected()
		result.ConceptStockDeleted += ra

		// 4. update concept_stock
		var oldStockMap map[string]*model.ConceptStock
		if oldc != nil {
			oldStockMap = make(map[string]*model.ConceptStock, len(oldc.Stocks))
			for _, oldsc := range oldc.Stocks {
				oldStockMap[oldsc.StockCode] = oldsc
			}
		} else {
			oldStockMap = make(map[string]*model.ConceptStock)
		}
		for _, newStock := range newc.Stocks {
			oldStock, ok := oldStockMap[newStock.StockCode]
			if !ok {
				// insert
				_, err = insertConceptStockStmt.Exec(newStock.Args()...)
				if err != nil {
					return result, errors.Wrap(err, repo.Name)
				}
				result.ConceptStockInserted++
			} else {
				if !newStock.Cmp(oldStock) {
					// update
					_, err = updateConceptStockStmt.Exec(newStock.StockName, newStock.Description, newStock.UpdatedAt, newStock.StockCode, newc.Id)
					if err != nil {
						return result, errors.Wrap(err, repo.Name)
					}
					result.ConceptStockUpdated++
				}
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return result, errors.Wrap(err, repo.Name)
	}

	return result, nil
}

const conceptLimit = 500

func (repo *StockMarketRepo) QueryConceptStockByKw(ctx context.Context, stockKw string, conceptKw string, limit int) (
	[]*model.ConceptStockView, error) {

	switch {
	case stockKw != "" && conceptKw != "":
		return repo.QueryConceptStockByStockConceptKw(ctx, stockKw, conceptKw, limit)
	case stockKw != "":
		return repo.QeuryConceptStockByStockKw(ctx, stockKw, limit)
	case conceptKw != "":
		return repo.QueryConceptStockByConceptKw(ctx, conceptKw, limit)
	default:
		return repo.QueryConceptStockByUpdatedDesc(ctx, limit)
	}
}

func (repo *StockMarketRepo) QueryConceptStockByStockConceptKw(ctx context.Context, stockKw string, conceptKw string, limit int) (
	[]*model.ConceptStockView, error) {

	if limit < 1 || limit > conceptLimit {
		limit = conceptLimit
	}
	scs := make([]*model.ConceptStockView, 0)
	scs, err := dbh.QueryContext[*model.ConceptStockView](repo.DB, ctx, tmpl.SelectConceptStockViewByStockConceptKw,
		stockKw, stockKw, conceptKw, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QeuryConceptStockByStockKw(ctx context.Context, stockKw string, limit int) ([]*model.ConceptStockView, error) {
	if limit < 1 || limit > conceptLimit {
		limit = conceptLimit
	}
	scs := make([]*model.ConceptStockView, 0)
	scs, err := dbh.QueryContext[*model.ConceptStockView](repo.DB, ctx, tmpl.SelectConceptStockViewByStockKw,
		stockKw, stockKw, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QueryConceptStockByConceptKw(ctx context.Context, conceptKw string, limit int) ([]*model.ConceptStockView, error) {
	if limit < 1 || limit > conceptLimit {
		limit = conceptLimit
	}
	scs := make([]*model.ConceptStockView, 0)
	scs, err := dbh.QueryContext[*model.ConceptStockView](repo.DB, ctx, tmpl.SelectConceptStockViewByConceptKw,
		conceptKw, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QueryConceptStockByUpdatedDesc(ctx context.Context, limit int) ([]*model.ConceptStockView, error) {
	if limit < 1 || limit > conceptLimit {
		limit = conceptLimit
	}
	scs := make([]*model.ConceptStockView, 0)
	scs, err := dbh.QueryContext[*model.ConceptStockView](repo.DB, ctx, tmpl.SelectConceptStockViewByUpdateAtDesc,
		limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QueryConcepts(ctx context.Context, conceptKw string, limit int, includeStock bool) ([]*model.Concept, error) {
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	var conceptVal interface{}
	if conceptKw != "" {
		conceptVal = conceptKw
	}
	concepts, err := dbh.QueryContext[*model.Concept](repo.DB, ctx,
		tmpl.SelectConceptByName,
		conceptVal, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	if includeStock {
		for _, concept := range concepts {
			stocks, err := repo.QueryConceptStockByConceptId(ctx, concept.Id)
			if err != nil {
				return nil, errors.Wrap(err, repo.Name)
			}
			concept.Stocks = stocks
		}
	}
	return concepts, nil
}

func (repo *StockMarketRepo) QueryConceptStockByConceptId(ctx context.Context, conceptId string) ([]*model.ConceptStock, error) {
	scs, err := dbh.QueryContext[*model.ConceptStock](repo.DB, ctx, tmpl.SelectConceptStockByConceptId,
		conceptId)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QueryRealtimeArchive(ctx context.Context, userId int, limit int) ([]*model.RealtimeMessage, error) {
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	messages, err := dbh.QueryContext[*model.RealtimeMessage](repo.DB, ctx,
		tmpl.SelectRealtimeByUserId,
		userId, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return messages, nil
}

func (repo *StockMarketRepo) SaveRealtimeArchive(ctx context.Context, message *model.RealtimeMessage) (int64, error) {
	return dbh.InsertContext(repo.DB, ctx, message)
}

func (repo *StockMarketRepo) DeleteRealtimeArchive(ctx context.Context, userId int, seq string) (int64, error) {
	ret, err := repo.DB.ExecContext(ctx, tmpl.DeleteRealtimeByUserIdSeq, userId, seq)
	if err != nil {
		return 0, errors.Wrap(err, repo.Name)
	}
	ra, _ := ret.RowsAffected()
	return ra, nil
}
