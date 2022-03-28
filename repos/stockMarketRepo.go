package repos

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"

	"github.com/joexzh/ThsConcept/db"
	"github.com/joexzh/ThsConcept/model"
	"github.com/pkg/errors"
)

type StockMarketRepo struct {
	db *sqlx.DB
}

func NewStockMarketRepo() (*StockMarketRepo, error) {
	client, err := db.GetMysqlClient()
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	return &StockMarketRepo{client}, nil
}

func (repo *StockMarketRepo) ZdtListDesc(ctx context.Context, start time.Time, limit int) ([]model.ZDTHistory, error) {
	if limit < 1 {
		limit = db.Limit
	}

	var list []model.ZDTHistory

	err := repo.db.SelectContext(ctx, &list,
		"SELECT * FROM long_short WHERE date >= ? ORDER BY date DESC LIMIT ?", start, limit)
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	return list, nil
}

func (repo *StockMarketRepo) InsertZdtList(ctx context.Context, list []model.ZDTHistory) (int64, error) {
	if len(list) < 1 {
		return 0, nil
	}

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, errors.Wrap(err, db.Mysql)
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex("INSERT INTO long_short VALUES (?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return 0, errors.Wrap(err, db.Mysql)
	}

	var total int64 = 0
	for _, zdt := range list {
		ret, err := stmt.ExecContext(ctx, zdt.Date, zdt.LongLimitCount, zdt.ShortLimitCount, zdt.StopTradeCount, zdt.Amount,
			zdt.SHLongCount, zdt.SHEvenCount, zdt.SHShortCount, zdt.SZLongCount, zdt.SZEvenCount, zdt.SZShortCount)
		if err != nil {
			return 0, errors.Wrap(err, db.Mysql)
		}
		cnt, _ := ret.RowsAffected()
		total += cnt
	}
	if err = tx.Commit(); err != nil {
		return 0, errors.Wrap(err, db.Mysql)
	}
	return total, nil
}

const stockModelSql = `SELECT
	s.code AS stock_code,
	s.name AS stock_name,
	c.id AS concept_id,
	c.name AS concept_name,
	sc.description,
	sc.updated_at 
FROM
	concept_stock AS s
	INNER JOIN concept_stock_concept AS sc ON sc.stock_code = s.CODE
	INNER JOIN concept_concept AS c ON c.id = sc.concept_id`

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

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return result, errors.Wrap(err, db.Mysql)
	}
	defer tx.Rollback()

	// query old concepts and stocks from db for update
	var oldcs []*model.Concept
	oldcsMap := make(map[string]*model.Concept)
	err = tx.SelectContext(ctx, &oldcs, "SELECT * FROM concept_concept for update")
	if err != nil {
		return result, errors.Wrap(err, db.Mysql)
	}
	stockByConceptIdStmt, err := tx.Preparex(stockModelSql + " WHERE c.id = ? for update")
	if err != nil {
		return result, errors.Wrap(err, db.Mysql)
	}
	for _, oldc := range oldcs {
		var stocks []*model.ConceptStock
		err = stockByConceptIdStmt.SelectContext(ctx, &stocks, oldc.Id)
		if err != nil {
			return result, errors.Wrap(err, db.Mysql)
		}
		oldc.Stocks = stocks
		oldcsMap[oldc.Id] = oldc
	}

	// 1. delete un-exist concepts
	cids := make([]string, 0, len(newcs))
	for _, c := range newcs {
		cids = append(cids, c.Id)
	}
	listSql, vals := db.ParamList(cids...)
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
	listSql, vals = db.ParamList(distinctCodes...)
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
	rows, err := tx.QueryxContext(ctx, "SELECT * FROM concept_stock")
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
			_, err = insertConceptStmt.Exec(newc.Id, newc.Name, newc.PlateId, newc.Define, newc.UpdatedAt)
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
		listSql, vals = db.ParamList(codes...)
		vals = append(vals, newc.Id)
		ret, err = tx.Exec(fmt.Sprintf("DELETE FROM concept_stock_concept WHERE stock_code NOT IN %s and concept_id=?", listSql), vals...)
		if err != nil {
			return result, errors.Wrap(err, db.Mysql)
		}
		ra, _ = ret.RowsAffected()
		result.ConceptStockConceptDeleted += ra

		// 6. update concept_stock_concept
		oldScMap := make(map[string]*model.ConceptStock, len(oldc.Stocks))
		for _, oldsc := range oldc.Stocks {
			oldScMap[oldsc.StockCode] = oldsc
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

func (repo *StockMarketRepo) QueryStockConcept(ctx context.Context, stockKw string, conceptKw string, limit int) (
	[]*model.Concept, error) {

	if limit < 0 || limit > 1000 {
		limit = 1000
	}
	// query stock concept
	const scSql = stockModelSql + `
WHERE
	s.CODE = ? 
	OR s.NAME = ? 
	OR c.NAME = ? 
ORDER BY
	sc.updated_at DESC 
	LIMIT ?`
	rows, err := repo.db.QueryxContext(ctx, scSql, stockKw, stockKw, conceptKw, limit)
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	stockMapByConcept := make(map[string][]*model.ConceptStock)
	for rows.Next() {
		var stock model.ConceptStock
		err = rows.StructScan(&stock)
		if err != nil {
			_ = rows.Close()
			return nil, errors.Wrap(err, db.Mysql)
		}
		stockMapByConcept[stock.ConceptId] = append(stockMapByConcept[stock.ConceptId], &stock)
	}

	// query concept
	conceptIds := make([]string, 0, len(stockMapByConcept))
	for conceptId := range stockMapByConcept {
		conceptIds = append(conceptIds, conceptId)
	}
	listSql, vals := db.ParamList(conceptIds...)
	rows, err = repo.db.QueryxContext(ctx,
		fmt.Sprintf("SELECT * FROM concept_concept WHERE id IN %s order by updated_at desc", listSql), vals...)
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	var concepts []*model.Concept
	for rows.Next() {
		var concept model.Concept
		err = rows.StructScan(&concept)
		if err != nil {
			_ = rows.Close()
			return nil, errors.Wrap(err, db.Mysql)
		}
		concept.Stocks = stockMapByConcept[concept.Id]
		concepts = append(concepts, &concept)
	}
	return concepts, nil
}

func (repo *StockMarketRepo) QueryRealtime() {
	// TODO query realtime
}

func (repo *StockMarketRepo) SaveRealtime() {
	// TODO save realtime
}

func (repo *StockMarketRepo) DeleteRealtime() {
	// TODO delete realtime
}

func (repo *StockMarketRepo) FuzzyStockKw() {
	// TODO fuzzy stock key word
}

func (repo *StockMarketRepo) FuzzyConceptKw() {
	// TODO fuzzy concept key word
}
