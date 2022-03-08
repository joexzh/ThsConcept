package repos

import (
	"context"
	"database/sql"
	"github.com/joexzh/ThsConcept/db"
	"github.com/joexzh/ThsConcept/model"
	"github.com/pkg/errors"
	"strings"
	"time"
)

const queryLongShortSql = `SELECT * FROM long_short WHERE date >= ? ORDER BY date DESC LIMIT ?`
const insertLongShortSql = `INSERT INTO long_short VALUES `

func prepareInsertLongShort(list []model.ZDTHistory) (string, []interface{}) {
	var vals []interface{}
	var builder strings.Builder
	builder.WriteString(insertLongShortSql)
	for _, zdt := range list {
		builder.WriteString("(?,?,?,?,?,?,?,?,?,?,?),")
		vals = append(vals, zdt.Date, zdt.LongLimitCount, zdt.ShortLimitCount, zdt.StopTradeCount, zdt.Amount,
			zdt.SHLongCount, zdt.SHEvenCount, zdt.SHShortCount, zdt.SZLongCount, zdt.SZEvenCount, zdt.SZShortCount)
	}

	_sql := builder.String()
	_sql = _sql[:len(_sql)-1]

	return _sql, vals
}

type StockMarketRepo struct {
	*sql.DB
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
	rows, err := repo.QueryContext(ctx, queryLongShortSql, start, limit)
	if err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}
	defer rows.Close()

	var list []model.ZDTHistory
	for rows.Next() {
		var zdt model.ZDTHistory
		err = rows.Scan(&zdt.Date, &zdt.LongLimitCount, &zdt.ShortLimitCount, &zdt.StopTradeCount, &zdt.Amount,
			&zdt.SHLongCount, &zdt.SHEvenCount, &zdt.SHShortCount, &zdt.SZLongCount, &zdt.SZEvenCount, &zdt.SZShortCount)
		if err != nil {
			return nil, errors.Wrap(err, db.Mysql)
		}
		list = append(list, zdt)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, db.Mysql)
	}

	return list, nil
}

func (repo *StockMarketRepo) InsertZdtList(ctx context.Context, list []model.ZDTHistory) (int64, error) {
	if len(list) < 1 {
		return 0, nil
	}
	_sql, vals := prepareInsertLongShort(list)
	ret, err := repo.ExecContext(ctx, _sql, vals...)
	if err != nil {
		return 0, errors.Wrap(err, db.Mysql)
	}
	rows, err := ret.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, db.Mysql)
	}
	return rows, nil
}
