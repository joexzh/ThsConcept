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
	"github.com/pkg/errors"
)

type StockMarketRepo struct {
	db *sql.DB
}

func NewStockMarketRepo() (*StockMarketRepo, error) {
	client, err := db.GetMysqlClient()
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	return &StockMarketRepo{client}, nil
}

func (repo *StockMarketRepo) ZdtListDesc(ctx context.Context, start time.Time, limit int) ([]*model.ZDTHistory, error) {
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	list, err := dbh.QueryContext[*model.ZDTHistory](repo.db, ctx,
		"SELECT * FROM long_short WHERE date >= ? ORDER BY date DESC LIMIT ?",
		start, limit)
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
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

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, errors.Wrap(err, db.Mysql)
	}
	defer tx.Rollback()

	if err != nil {
		return 0, errors.Wrap(err, db.Mysql)
	}

	total, err := dbh.BulkInsertContext(tx, ctx, 1000, list...)
	if err != nil {
		return 0, errors.Wrap(err, db.Mysql)
	}
	if err = tx.Commit(); err != nil {
		return 0, errors.Wrap(err, db.Mysql)
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

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return result, errors.Wrap(err, db.Mysql)
	}
	defer tx.Rollback()

	// query old concepts and stocks from db for update
	oldscs, err := dbh.QueryContext[*model.ConceptStock](repo.db, ctx, `SELECT
	s.CODE AS stock_code,
	s.NAME AS stock_name,
	sc.updated_at,
	sc.description,
	c.id AS concept_id,
	c.NAME AS concept_name,
	c.plate_id AS concept_plate_id,
	c.define AS concept_define,
	c.updated_at AS concept_updated_at 
FROM
	concept_stock AS s
	INNER JOIN concept_stock_concept AS sc ON sc.stock_code = s.
	CODE INNER JOIN concept_concept AS c ON c.id = sc.concept_id
for update`)
	if err != nil {
		return result, errors.Wrap(err, db.Mysql)
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
		return result, errors.Wrap(err, db.Mysql)
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
		return result, errors.Wrap(err, db.Mysql)
	}
	ra, _ = ret.RowsAffected()
	result.ConceptStockDeleted = ra

	// 3. update concept_stock
	updateStockStmt, err := tx.Prepare("UPDATE concept_stock SET name = ? WHERE code = ?")
	if err != nil {
		return result, errors.Wrap(err, db.Mysql)
	}
	insertStockStmt, err := tx.Prepare("insert INTO concept_stock VALUES (?,?)")
	if err != nil {
		return result, errors.Wrap(err, db.Mysql)
	}
	distinctOldStocks := make(map[string]*model.ConceptStock)
	rows, err := tx.QueryContext(ctx, "SELECT * FROM concept_stock")
	if err != nil {
		return result, errors.Wrap(err, db.Mysql)
	}
	for rows.Next() {
		var oldStock model.ConceptStock
		err = rows.Scan(&oldStock.StockCode, &oldStock.StockName)
		if err != nil {
			return result, errors.Wrap(err, db.Mysql)
		}
		distinctOldStocks[oldStock.StockCode] = &oldStock
	}
	for _, newstock := range distinctNewStocks {
		oldstock, ok := distinctOldStocks[newstock.StockCode]
		if !ok {
			// insert
			_, err = insertStockStmt.Exec(newstock.StockCode, newstock.StockName)
			if err != nil {
				return result, errors.Wrap(err, db.Mysql)
			}
			result.ConceptStockInserted++
		} else {
			if newstock.StockName != "" && !newstock.CmpStock(oldstock) {
				// update
				_, err = updateStockStmt.Exec(newstock.StockName, newstock.StockCode)
				if err != nil {
					return result, errors.Wrap(err, db.Mysql)
				}
				log.Println("update stock", newstock.StockCode, newstock.StockName)
				result.ConceptStockUpdated++
			}
		}
	}

	insertConceptStmt, err := tx.Prepare("insert INTO concept_concept VALUES (?,?,?,?,?)")
	if err != nil {
		return result, errors.Wrap(err, db.Mysql)
	}
	updateConceptStmt, err := tx.Prepare("UPDATE concept_concept SET name = ?, plate_id = ?, define = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return result, errors.Wrap(err, db.Mysql)
	}
	insertScStmt, err := tx.Prepare("insert INTO concept_stock_concept VALUES (?,?,?,?)")
	if err != nil {
		return result, errors.Wrap(err, db.Mysql)
	}
	updateScStmt, err := tx.Prepare("UPDATE concept_stock_concept SET description = ?, updated_at = ? WHERE stock_code = ? AND concept_id = ?")
	if err != nil {
		return result, errors.Wrap(err, db.Mysql)
	}

	for _, newc := range newcs {
		// 4. update concept_concept
		oldc, ok := oldcsMap[newc.Id]
		if !ok {
			// insert
			_, err = insertConceptStmt.Exec(newc.Args()...)
			if err != nil {
				return result, errors.Wrap(err, db.Mysql)
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
			return result, errors.Wrap(err, db.Mysql)
		}
		ra, _ = ret.RowsAffected()
		result.ConceptStockConceptDeleted += ra

		// 6. update concept_stock_concept
		oldscSize := 0
		if oldc != nil {
			oldscSize = len(oldc.Stocks)
		}
		oldScMap := make(map[string]*model.ConceptStock, oldscSize)
		if oldc != nil {
			for _, oldsc := range oldc.Stocks {
				oldScMap[oldsc.StockCode] = oldsc
			}
		}
		for _, newsc := range newc.Stocks {
			oldsc, ok := oldScMap[newsc.StockCode]
			if !ok {
				// insert
				_, err = insertScStmt.Exec(newsc.StockCode, newc.Id, newsc.Description, newsc.UpdatedAt)
				if err != nil {
					return result, errors.Wrap(err, db.Mysql)
				}
				result.ConceptStockConceptInserted++
			} else {
				if newsc.Description != "" && !newsc.CmpConcept(oldsc) {
					// update
					_, err = updateScStmt.Exec(newsc.Description, newsc.UpdatedAt, newsc.StockCode, newc.Id)
					if err != nil {
						return result, errors.Wrap(err, db.Mysql)
					}
					result.ConceptStockConceptUpdated++
				}
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return result, errors.Wrap(err, db.Mysql)
	}

	return result, nil
}

func (repo *StockMarketRepo) QueryConceptStock(ctx context.Context, stockKw string, conceptKw string, limit int) (
	[]*model.ConceptStock, error) {

	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	const scSql = `SELECT
	s.CODE AS stock_code,
	s.NAME AS stock_name,
	sc.updated_at,
	sc.description,
	c.id AS concept_id,
	c.NAME AS concept_name,
	c.plate_id AS concept_plate_id,
	c.define AS concept_define,
	c.updated_at AS concept_updated_at 
FROM
	concept_stock AS s
	INNER JOIN concept_stock_concept AS sc ON sc.stock_code = s.
	CODE INNER JOIN concept_concept AS c ON c.id = sc.concept_id
WHERE
	(
		s.CODE = IFNULL(?, s.code)
		OR s.NAME = IFNULL(?, s.name)
	) 
	and c.NAME = IFNULL(?, c.name) 
ORDER BY
	sc.updated_at DESC 
	LIMIT ?`
	vals := make([]interface{}, 4)
	if stockKw == "" {
		vals[0] = nil
		vals[1] = nil
	} else {
		vals[0] = stockKw
		vals[1] = stockKw
	}
	if conceptKw == "" {
		vals[2] = nil
	} else {
		vals[2] = conceptKw
	}
	vals[3] = limit

	scs := make([]*model.ConceptStock, 0)
	scs, err := dbh.QueryContext[*model.ConceptStock](repo.db, ctx, scSql, vals...)
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
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
	concepts, err := dbh.QueryContext[*model.Concept](repo.db, ctx,
		"select * from concept_concept where name=IFNULL(?, name) order by updated_at desc limit ?",
		conceptVal, limit)
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	return concepts, nil
}

func (repo *StockMarketRepo) QueryStockByConceptId(ctx context.Context, conceptId string) ([]*model.ConceptStock, error) {
	scs, err := dbh.QueryContext[*model.ConceptStock](repo.db, ctx, `SELECT
	s.CODE AS stock_code,
	s.NAME AS stock_name,
	sc.updated_at,
	sc.description,
	c.id AS concept_id,
	c.NAME AS concept_name,
	c.plate_id AS concept_plate_id,
	c.define AS concept_define,
	c.updated_at AS concept_updated_at 
FROM
	concept_stock AS s
	INNER JOIN concept_stock_concept AS sc ON sc.stock_code = s.
	CODE INNER JOIN concept_concept AS c ON c.id = sc.concept_id
where c.id=?
order by sc.updated_at`, conceptId)
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	return scs, nil
}

func (repo *StockMarketRepo) QueryRealtimeArchive(ctx context.Context, userId int, limit int) ([]*model.RealtimeMessage, error) {
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	messages, err := dbh.QueryContext[*model.RealtimeMessage](repo.db, ctx,
		"select * from realtime_archive where user_id=? order by seq desc limit ?",
		userId, limit)
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	return messages, nil
}

func (repo *StockMarketRepo) SaveRealtimeArchive(ctx context.Context, message *model.RealtimeMessage) (int64, error) {
	return dbh.InsertContext(repo.db, ctx, message)
}

func (repo *StockMarketRepo) DeleteRealtimeArchive(ctx context.Context, userId int, seq string) (int64, error) {
	ret, err := repo.db.ExecContext(ctx, "delete from realtime_archive where user_id=? and seq=?", userId, seq)
	if err != nil {
		return 0, errors.Wrap(err, db.Mysql)
	}
	ra, _ := ret.RowsAffected()
	return ra, nil
}

func (repo *StockMarketRepo) FuzzyStockKw() {
	// TODO fuzzy stock key word
}

func (repo *StockMarketRepo) FuzzyConceptKw() {
	// TODO fuzzy concept key word
}
