package db

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const Mysql = "mysql"

// NewMysqlClient create mysql client, not connect yet.
func NewMysqlClient(dsn string) (*sql.DB, error) {
	pool, err := sql.Open(Mysql, dsn)
	if err != nil {
		return nil, err
	}
	pool.SetConnMaxLifetime(3 * time.Minute)
	pool.SetConnMaxIdleTime(10)
	pool.SetMaxOpenConns(10)

	err = Ping(pool)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func Ping(pool *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
