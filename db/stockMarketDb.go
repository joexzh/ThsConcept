package db

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joexzh/ThsConcept/config"
)

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
