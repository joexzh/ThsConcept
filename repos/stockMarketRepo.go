package repos

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	list, err := dbh.QueryContext(repo.DB, ctx, tmpl.SelectZdt,
		func() *model.ZDTHistory { return new(model.ZDTHistory) },
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

func scsToConceptMap(scs []*model.ConceptStock) map[string]*model.Concept {
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
		concept.Stocks = append(concept.Stocks, sc)
	}
	return conceptMap
}

type UpdateConceptResult struct {
	ConceptConceptInserted      int64
	ConceptConceptUpdated       int64
	ConceptConceptDeleted       int64
	ConceptStockInserted        int64
	ConceptStockUpdated         int64
	ConceptStockDeleted         int64
	ConceptStockConceptInserted int64
	ConceptStockConceptUpdated  int64
	ConceptStockConceptDeleted  int64
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
	oldscs, err := dbh.QueryContext(repo.DB, ctx, tmpl.SelectAllSc,
		func() *model.ConceptStock { return new(model.ConceptStock) })
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

	// 2. delete un-exist stocks
	var distinctCodes []string
	distinctNewStocks := make(map[string]*model.ConceptStock)
	for _, newc := range newcs {
		for _, newstock := range newc.Stocks {
			stock, ok := distinctNewStocks[newstock.StockCode]
			if !ok {
				distinctNewStocks[newstock.StockCode] = newstock
				distinctCodes = append(distinctCodes, newstock.StockCode)
			} else {
				// new stock name may be empty, double check it
				if stock.StockName == "" {
					distinctNewStocks[newstock.StockCode] = newstock
				}
			}
		}
	}
	listSql, vals = db.ArgList(distinctCodes...)
	ret, err = tx.Exec("DELETE FROM concept_stock WHERE code NOT IN "+listSql, vals...)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	ra, _ = ret.RowsAffected()
	result.ConceptStockDeleted = ra

	// 3. update concept_stock
	updateStockStmt, err := tx.Prepare(tmpl.UpdateConceptStock)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	insertStockStmt, err := tx.Prepare(tmpl.InsertConceptStock)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	distinctOldStocks := make(map[string]*model.ConceptStock)
	rows, err := tx.QueryContext(ctx, tmpl.SelectAllConceptStock)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	for rows.Next() {
		var oldStock model.ConceptStock
		err = rows.Scan(&oldStock.StockCode, &oldStock.StockName)
		if err != nil {
			return result, errors.Wrap(err, repo.Name)
		}
		distinctOldStocks[oldStock.StockCode] = &oldStock
	}
	for _, newstock := range distinctNewStocks {
		oldstock, ok := distinctOldStocks[newstock.StockCode]
		if !ok {
			// insert
			_, err = insertStockStmt.Exec(newstock.StockCode, newstock.StockName)
			if err != nil {
				return result, errors.Wrap(err, repo.Name)
			}
			result.ConceptStockInserted++
		} else {
			if newstock.StockName != "" && !newstock.CmpStock(oldstock) {
				// update
				_, err = updateStockStmt.Exec(newstock.StockName, newstock.StockCode)
				if err != nil {
					return result, errors.Wrap(err, repo.Name)
				}
				log.Println("update stock", newstock.StockCode, newstock.StockName)
				result.ConceptStockUpdated++
			}
		}
	}

	insertConceptStmt, err := tx.Prepare(tmpl.InsertConcept)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	updateConceptStmt, err := tx.Prepare(tmpl.UpdateConcept)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	insertScStmt, err := tx.Prepare(tmpl.InsertConceptStockConcept)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}
	updateScStmt, err := tx.Prepare(tmpl.UpdateConceptStockConcept)
	if err != nil {
		return result, errors.Wrap(err, repo.Name)
	}

	for _, newc := range newcs {
		// 4. update concept_concept
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

		// 5. delete un-exist concept_stock_concept
		codes := make([]string, len(newc.Stocks))
		for i, stock := range newc.Stocks {
			codes[i] = stock.StockCode
		}
		listSql, vals = db.ArgList(codes...)
		vals = append(vals, newc.Id)
		ret, err = tx.Exec(fmt.Sprintf("DELETE FROM concept_stock_concept WHERE stock_code NOT IN %s and concept_id=?", listSql), vals...)
		if err != nil {
			return result, errors.Wrap(err, repo.Name)
		}
		ra, _ = ret.RowsAffected()
		result.ConceptStockConceptDeleted += ra

		// 6. update concept_stock_concept
		var oldScMap map[string]*model.ConceptStock
		if oldc != nil {
			oldScMap = make(map[string]*model.ConceptStock, len(oldc.Stocks))
			for _, oldsc := range oldc.Stocks {
				oldScMap[oldsc.StockCode] = oldsc
			}
		} else {
			oldScMap = make(map[string]*model.ConceptStock)
		}
		for _, newsc := range newc.Stocks {
			oldsc, ok := oldScMap[newsc.StockCode]
			if !ok {
				// insert
				_, err = insertScStmt.Exec(newsc.StockCode, newc.Id, newsc.Description, newsc.UpdatedAt)
				if err != nil {
					return result, errors.Wrap(err, repo.Name)
				}
				result.ConceptStockConceptInserted++
			} else {
				if newsc.Description != "" && !newsc.CmpConcept(oldsc) {
					// update
					_, err = updateScStmt.Exec(newsc.Description, newsc.UpdatedAt, newsc.StockCode, newc.Id)
					if err != nil {
						return result, errors.Wrap(err, repo.Name)
					}
					result.ConceptStockConceptUpdated++
				}
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return result, errors.Wrap(err, repo.Name)
	}

	return result, nil
}

func (repo *StockMarketRepo) QueryScByKw(ctx context.Context, stockKw string, conceptKw string, limit int) (
	[]*model.ConceptStock, error) {

	switch {
	case stockKw != "" && conceptKw != "":
		return repo.QueryScByStockConceptKw(ctx, stockKw, conceptKw, limit)
	case stockKw != "":
		return repo.QeuryScByStockKw(ctx, stockKw, limit)
	case conceptKw != "":
		return repo.QueryScByConceptKw(ctx, conceptKw, limit)
	default:
		return repo.QueryScByUpdatedDesc(ctx, limit)
	}
}

func (repo *StockMarketRepo) QueryScByStockConceptKw(ctx context.Context, stockKw string, conceptKw string, limit int) (
	[]*model.ConceptStock, error) {

	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	scs := make([]*model.ConceptStock, 0)
	scs, err := dbh.QueryContext(repo.DB, ctx, tmpl.SelectScByStockConceptKw,
		func() *model.ConceptStock {
			return new(model.ConceptStock)
		},
		stockKw, conceptKw, conceptKw, conceptKw, conceptKw, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QeuryScByStockKw(ctx context.Context, stockKw string, limit int) ([]*model.ConceptStock, error) {
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	scs := make([]*model.ConceptStock, 0)
	scs, err := dbh.QueryContext(repo.DB, ctx, tmpl.SelectScByStockKw,
		func() *model.ConceptStock {
			return new(model.ConceptStock)
		},
		stockKw, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QueryScByConceptKw(ctx context.Context, conceptKw string, limit int) ([]*model.ConceptStock, error) {
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	scs := make([]*model.ConceptStock, 0)
	scs, err := dbh.QueryContext(repo.DB, ctx, tmpl.SelectScByConceptKw,
		func() *model.ConceptStock {
			return new(model.ConceptStock)
		},
		conceptKw, conceptKw, conceptKw, conceptKw, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QueryScByUpdatedDesc(ctx context.Context, limit int) ([]*model.ConceptStock, error) {
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	scs := make([]*model.ConceptStock, 0)
	scs, err := dbh.QueryContext(repo.DB, ctx, tmpl.SelectScByUpdateAtDesc,
		func() *model.ConceptStock {
			return new(model.ConceptStock)
		},
		limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QueryConcepts(ctx context.Context, conceptKw string, limit int) ([]*model.Concept, error) {
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	var conceptVal interface{}
	if conceptKw != "" {
		conceptVal = conceptKw
	}
	concepts, err := dbh.QueryContext(repo.DB, ctx,
		tmpl.SelectConceptByName,
		func() *model.Concept {
			return new(model.Concept)
		},
		conceptVal, limit)
	if err != nil {
		return nil, errors.Wrap(err, repo.Name)
	}
	return concepts, nil
}

func (repo *StockMarketRepo) QueryStockByConceptId(ctx context.Context, conceptId string) ([]*model.ConceptStock, error) {
	scs, err := dbh.QueryContext(repo.DB, ctx, tmpl.SelectScByConceptId,
		func() *model.ConceptStock {
			return new(model.ConceptStock)
		},
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
	messages, err := dbh.QueryContext(repo.DB, ctx,
		tmpl.SelectRealtimeByUserId,
		func() *model.RealtimeMessage {
			return new(model.RealtimeMessage)
		},
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
