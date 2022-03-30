package db

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/dbh"
)

const (
	Mysql = "mysql"
	Limit = 10000
)

type mysqlClient struct {
	pool *sql.DB
	err  error
}

var _mysqlClient *mysqlClient
var once sync.Once

func GetMysqlClient() (*sql.DB, error) {
	once.Do(func() {
		db, err := newMysqlClient(config.GetEnv().MysqlConnStr)
		_mysqlClient = &mysqlClient{
			pool: db,
			err:  err,
		}
	})
	return _mysqlClient.pool, _mysqlClient.err
}

// newMysqlClient create mysql client, not connect yet.
func newMysqlClient(dsn string) (*sql.DB, error) {
	pool, err := sql.Open(Mysql, dsn)
	if err != nil {
		return nil, err
	}
	pool.SetConnMaxLifetime(3 * time.Minute)
	pool.SetConnMaxIdleTime(10 * time.Second)
	pool.SetMaxOpenConns(10)
	pool.SetMaxIdleConns(10)

	return pool, nil
}

// ArgList generates sql list part, (?,?,...?), question mark is used to replace the value.
func ArgList[T any](params ...T) (string, []interface{}) {
	vals := make([]interface{}, len(params))
	for i, param := range params {
		vals[i] = param
	}
	return dbh.DefaultConfig.MarkInsertValueSql(len(params), 1), vals
}
