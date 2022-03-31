package db

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/dbh"
)

type DBConfig struct {
	Driver       string
	DSN          string
	MaxLifetime  time.Duration
	MaxIdletime  time.Duration
	MaxOpenConns int
	MaxIdleConns int
}

func NewMysqlConfig() *DBConfig {
	return &DBConfig{
		Driver:       "mysql",
		DSN:          config.GetEnv().MysqlConnStr,
		MaxLifetime:  3 * time.Minute,
		MaxIdletime:  10 * time.Second,
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
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
		if err = db.Ping(); err != nil {
			log.Fatal(err)
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
