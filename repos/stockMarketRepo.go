package repos

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joexzh/ThsConcept/config"
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

const scModelSql = `SELECT
	s.code AS stock_code,
	s.name AS stock_name,
	s.pinyin_first_letter,
	s.pinyin_normal,
	c.id AS concept_id,
	c.name AS concept_name,
	sc.description,
	sc.updated_at 
FROM
	concept_stock AS s
	INNER JOIN concept_stock_concept AS sc ON sc.stock_code = s.CODE
	INNER JOIN concept_concept AS c ON c.id = sc.concept_id`

func (repo *StockMarketRepo) UpdateConcept(ctx context.Context, newcs ...*model.Concept) (map[string]int64, error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	defer tx.Rollback()

	results := make(map[string]int64, 4)

	// 1. delete un-exist concepts
	cids := make([]string, 0, len(newcs))
	for _, c := range newcs {
		cids = append(cids, c.Id)
	}
	listSql, vals := db.ParamList(cids)
	ret, err := tx.Exec("DELETE FROM concept_concept WHERE id NOT IN "+listSql, vals...)
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	ra, _ := ret.RowsAffected()
	results["concept_deleted"] = ra

	// 2. delete un-exist stocks
	var scodes []string
	for _, c := range newcs {
		for _, s := range c.Stocks {
			scodes = append(scodes, s.StockCode)
		}
	}
	listSql, vals = db.ParamList(scodes)
	_, err = tx.Exec("DELETE FROM concept_stock WHERE code NOT IN "+listSql, vals...)
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}

	// query all db concepts and stocks for later compare
	dbConceptsMap := make(map[string]*model.Concept)
	dbStocksMap := make(map[string]*model.ConceptStock)
	rows, err := tx.QueryxContext(ctx, "SELECT * FROM concept_concept")
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	for rows.Next() {
		var concept model.Concept
		err = rows.StructScan(&concept)
		if err != nil {
			_ = rows.Close()
			return nil, errors.Wrap(err, db.Mysql)
		}
		concept.UpdatedAt = concept.UpdatedAt.In(config.ChinaLoc()).Add(-time.Hour * 8)
		dbConceptsMap[concept.Id] = &concept
	}
	rows, err = tx.QueryxContext(ctx, "SELECT * FROM concept_stock")
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	for rows.Next() {
		var stock model.ConceptStock
		err = rows.Scan(&stock.StockCode, &stock.StockName, &stock.PinyinFirstLetter, &stock.PinyinNormal)
		if err != nil {
			_ = rows.Close()
			return nil, errors.Wrap(err, db.Mysql)
		}
		dbStocksMap[stock.StockCode] = &stock
	}

	// prepared stmts
	replaceConceptStmt, err := tx.Prepare("REPLACE INTO concept_concept VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	stockByConceptIdStmt, err := tx.Preparex(scModelSql + ` WHERE c.id = ?`)
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	scReplaceStmt, err := tx.Prepare("REPLACE INTO concept_stock_concept VALUES (?,?,?,?)")
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	stockReplaceStmt, err := tx.Prepare("REPLACE INTO concept_stock VALUES (?,?,?,?)")
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}

	for _, concept := range newcs {
		// 3. replace un-exist or outdated concept_concept
		dbConcept, ok := dbConceptsMap[concept.Id]
		if !ok || !concept.Cmp(dbConcept) {
			ret, err = replaceConceptStmt.Exec(
				concept.Id, concept.Name, concept.PinyinFirstLetter, concept.PinyinNormal, concept.PlateId, concept.Define, concept.UpdatedAt)
			if err != nil {
				return nil, errors.Wrap(err, db.Mysql)
			}
			ra, _ = ret.RowsAffected()
			results["concept_updated"] = results["concept_updated"] + ra
		}

		// 4. replace un-exist or outdated concept_stock
		distinctStockMap := make(map[string]*model.ConceptStock)
		for _, stock := range concept.Stocks {
			distinctStockMap[stock.StockCode] = stock
		}
		for stockCode, stock := range distinctStockMap {
			dbstock, ok := dbStocksMap[stockCode]
			if !ok || !stock.CmpStock(dbstock) {
				ret, err = stockReplaceStmt.Exec(stockCode, stock.StockName, stock.PinyinFirstLetter, stock.PinyinNormal)
				if err != nil {
					return nil, errors.Wrap(err, db.Mysql)
				}
			}
		}

		// 5. delete un-exist concept_stock_concept
		distinctStockCodes := make([]string, 0, len(concept.Stocks))
		for code := range distinctStockMap {
			distinctStockCodes = append(distinctStockCodes, code)

		}
		listSql, vals = db.ParamList(distinctStockCodes)
		vals = append(vals[:1], vals)
		vals[0] = concept.Id
		ret, err = tx.Exec("delete from concept_stock_concept where concept_id = ? and stock_code not in "+listSql, vals...)
		if err != nil {
			return nil, errors.Wrap(err, db.Mysql)
		}
		ra, _ = ret.RowsAffected()
		results["stock_deleted"] = results["stock_deleted"] + ra

		// query concept_stock_concept by concept_id then put in map
		rows, err = stockByConceptIdStmt.QueryxContext(ctx, concept.Id)
		if err != nil {
			return nil, errors.Wrap(err, db.Mysql)
		}
		dbStockMapByCodeAndId := make(map[string]*model.ConceptStock)
		for rows.Next() {
			var dbstock model.ConceptStock
			err = rows.StructScan(&dbstock)
			if err != nil {
				_ = rows.Close()
				return nil, errors.Wrap(err, db.Mysql)
			}
			dbstock.UpdatedAt = dbstock.UpdatedAt.In(config.ChinaLoc()).Add(-time.Hour * 8)
			dbStockMapByCodeAndId[dbstock.StockCode+dbstock.ConceptId] = &dbstock
		}

		// 6. compare concept_stock_concept, if not exist or outdated, replace it
		for _, stock := range concept.Stocks {
			dbstock, ok := dbStockMapByCodeAndId[stock.StockCode+stock.ConceptId]
			if !ok || !stock.CmpConcept(dbstock) {
				ret, err = scReplaceStmt.Exec(stock.StockCode, stock.ConceptId, stock.Description, stock.UpdatedAt)
				if err != nil {
					return nil, errors.Wrap(err, db.Mysql)
				}
				ra, _ = ret.RowsAffected()
				results["stock_updated"] = results["stock_updated"] + ra
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}

	return results, nil
}

func (repo *StockMarketRepo) QueryStockConcept(ctx context.Context, stockKw string, conceptKw string, limit int) (
	[]*model.Concept, error) {

	if limit < 0 || limit > 1000 {
		limit = 1000
	}
	// query stock concept
	const scSql = scModelSql + `
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
		stock.UpdatedAt = stock.UpdatedAt.In(config.ChinaLoc()).Add(-time.Hour * 8)
		stockMapByConcept[stock.ConceptId] = append(stockMapByConcept[stock.ConceptId], &stock)
	}

	// query concept
	conceptIds := make([]string, 0, len(stockMapByConcept))
	for conceptId := range stockMapByConcept {
		conceptIds = append(conceptIds, conceptId)
	}
	listSql, vals := db.ParamList(conceptIds)
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
		concept.UpdatedAt = concept.UpdatedAt.In(config.ChinaLoc()).Add(-time.Hour * 8)
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
