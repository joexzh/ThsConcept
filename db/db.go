package db

import (
	"database/sql"
	"sync"
	"time"

	"github.com/joexzh/dbh"
)

type DBConfig struct {
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	MaxIdletime  time.Duration
	Driver       string
	DSN          string
}

var (
	_db        *sql.DB
	_err       error
	_newDbOnce sync.Once
)

func NewDB(config *DBConfig) (*sql.DB, error) {
	_newDbOnce.Do(func() {
		db, err := sql.Open(config.Driver, config.DSN)
		if err != nil {
			_err = err
			return
		}
		db.SetConnMaxLifetime(config.MaxLifetime)
		db.SetConnMaxIdleTime(config.MaxIdletime)
		db.SetMaxOpenConns(config.MaxOpenConns)
		db.SetMaxIdleConns(config.MaxIdleConns)
		_db = db
	})
	return _db, _err
}

// ArgList generates sql list part, (?,?,...?), question mark is used to replace the value.
func ArgList[T any](params ...T) (string, []interface{}) {
	vals := make([]interface{}, len(params))
	for i, param := range params {
		vals[i] = param
	}
	return dbh.DefaultConfig.MarkInsertValueSql(len(params), 1), vals
}
