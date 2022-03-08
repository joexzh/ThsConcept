package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joexzh/ThsConcept/config"
	"time"
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

func init() {
	db, err := newMysqlClient(config.GetEnv().MysqlConnStr)
	_mysqlClient = &mysqlClient{
		pool: db,
		err:  err,
	}
}

func GetMysqlClient() (*sql.DB, error) {
	return _mysqlClient.pool, _mysqlClient.err
}

// newMysqlClient create mysql client, not connect yet.
func newMysqlClient(dsn string) (*sql.DB, error) {
	pool, err := sql.Open(Mysql, dsn)
	if err != nil {
		return nil, err
	}
	pool.SetConnMaxLifetime(3 * time.Minute)
	pool.SetConnMaxIdleTime(10)
	pool.SetMaxOpenConns(10)

	return pool, nil
}
