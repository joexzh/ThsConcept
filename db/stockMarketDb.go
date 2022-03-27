package db

import (
	"github.com/jmoiron/sqlx"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joexzh/ThsConcept/config"
)

const (
	Mysql = "mysql"
	Limit = 10000
)

type mysqlClient struct {
	pool *sqlx.DB
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

func GetMysqlClient() (*sqlx.DB, error) {
	return _mysqlClient.pool, _mysqlClient.err
}

// newMysqlClient create mysql client, not connect yet.
func newMysqlClient(dsn string) (*sqlx.DB, error) {
	pool, err := sqlx.Open(Mysql, dsn)
	if err != nil {
		return nil, err
	}
	pool.SetConnMaxLifetime(0)
	pool.SetConnMaxIdleTime(5 * time.Second)
	pool.SetMaxOpenConns(10)

	return pool, nil
}

// ParamList generates sql list part, (?,?,...?), question mark is used to replace the value.
func ParamList(params ...interface{}) (string, []interface{}) {
	var b strings.Builder
	b.WriteString("(")
	for i, _ := range params {
		b.WriteString("?")
		if i < len(params)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString(")")
	return b.String(), params
}
