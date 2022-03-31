package db

import (
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joexzh/ThsConcept/config"
)

var mysqlConfig *DBConfig
var mysqlConfigOnce sync.Once

func NewMysqlConfig() *DBConfig {
	mysqlConfigOnce.Do(func() {
		mysqlConfig = &DBConfig{
			Driver:       "mysql",
			DSN:          config.GetEnv().MysqlConnStr,
			MaxLifetime:  3 * time.Minute,
			MaxIdletime:  10 * time.Second,
			MaxOpenConns: 10,
			MaxIdleConns: 10,
		}
	})
	return mysqlConfig
}
