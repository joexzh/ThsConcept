package repos

import (
	"context"
	"database/sql"
	"github.com/joexzh/ThsConcept/db"
	"github.com/joexzh/ThsConcept/model"
	"github.com/pkg/errors"
	"strings"
	"sync"
	"time"
)

const queryLongShortSql = `SELECT * FROM long_short WHERE date >= ? `
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

var (
	getRepoOnce = sync.Once{}
	repo        *StockMarketRepo
	repoErr     error
)

func GetStockMarketRepo() (*StockMarketRepo, error) {
	getRepoOnce.Do(func() {
		r, err := NewStockMarketRepo()
		if err != nil {
			repoErr = err
		}
		repo = r
	})
	return repo, repoErr
}

type DateOrder string

const (
	DateAsc  = DateOrder("")
	DateDesc = DateOrder("ORDER BY date DESC")
)

func (repo *StockMarketRepo) QueryLongShort(ctx context.Context, start time.Time, order DateOrder, limit int) ([]model.ZDTHistory, error) {
	sb := strings.Builder{}
	sb.WriteString(queryLongShortSql)
	vars := []interface{}{start}
	sb.WriteString(string(order))
	if limit > 0 {
		sb.WriteString(" LIMIT ?")
		vars = append(vars, limit)
	}

	rows, err := repo.QueryContext(ctx, sb.String(), vars...)
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

func (repo *StockMarketRepo) InsertLongShort(ctx context.Context, list []model.ZDTHistory) (int64, error) {
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
