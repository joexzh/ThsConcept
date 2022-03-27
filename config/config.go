package config

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	Db               = "ThsConcept"
	CollConcept      = "concepts"
	CollStockConcept = "stockConcepts"
	CollRealtime     = "realtime"

	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"

	RexStockSymbol = `~(?:[a-z]{2})?([0-9]+)`
	StockSymbolUrl = "http://www.shdjt.com/js/lib/astock.js"

	ConceptPageUrl       = "http://basic.10jqka.com.cn/%v/concept.html"
	ConceptAllUrl        = "http://q.10jqka.com.cn/gn/"
	ConceptDetailPageUrl = "http://q.10jqka.com.cn/gn/detail/code/%v/"
	RealTimeUrl          = "https://news.10jqka.com.cn/tapp/news/push/stock/"

	RexValidConceptPage = `<h1 style="margin:3px 0px 0px 0px">\s*\d{6}\s*</h1>`
	RexConceptId        = `cid="(\d*)"`
	RexConceptCode      = `http://q.10jqka.com.cn/gn/detail/code/(\d*)`
	RexConceptDefine    = `<h4>[^\x00-\xff]+</h4>\s*<p style="overflow-y:auto;">([^<>]+)</p>`
	ConceptApiUrl       = "http://basic.10jqka.com.cn/ajax/stock/conceptlist.php?cid=%v&code=601127"

	TimeLayoutDate = "2006-01-02"
	TimeLayoutHour = "2006-01-02 15:04:05"
	TimeLocal      = "Asia/Shanghai"

	// Throttle 并发 goroutine 的数量, 防止系统或网络崩溃, 不能为0, 死锁!
	Throttle = 3
	// SleepRandUpTo 每个goroutine随机间隔的最大毫秒, 最低是0
	SleepRandUpTo = 500
)

var loc *time.Location

type Env struct {
	MongoUser     string
	MongoPassword string
	MongoHostPort string
	MongoConnStr  string

	ServerPort string

	MysqlUser     string
	MysqlPassword string
	MysqlHost     string
	MysqlPort     string
	MysqlConnStr  string
}

var env = Env{}
var once sync.Once

func init() {
	loc, _ = time.LoadLocation(TimeLocal)
}

func GetEnv() *Env {
	once.Do(func() {
		if os.Getenv("MONGO_USER") == "" {
			panic("MONGO_USER is not set")
		}
		env.MongoUser = os.Getenv("MONGO_USER")

		if os.Getenv("MONGO_PASSWORD") == "" {
			panic("MONGO_PASSWORD is not set")
		}
		env.MongoPassword = os.Getenv("MONGO_PASSWORD")

		if os.Getenv("MONGO_HOST_PORT") == "" {
			panic("MONGO_HOST_PORT is not set")
		}
		env.MongoHostPort = os.Getenv("MONGO_HOST_PORT")

		if os.Getenv("mysql_user") == "" {
			panic("mysql_user is not set")
		}
		env.MysqlUser = os.Getenv("mysql_user")

		if os.Getenv("mysql_password") == "" {
			panic("mysql_password is not set")
		}
		env.MysqlPassword = os.Getenv("mysql_password")

		if os.Getenv("mysql_host") == "" {
			panic("mysql_host is not set")
		}
		env.MysqlHost = os.Getenv("mysql_host")

		if os.Getenv("mysql_port") == "" {
			panic("mysql_port is not set")
		}
		env.MysqlPort = os.Getenv("mysql_port")

		env.MongoConnStr = fmt.Sprintf(`mongodb://%s:%s@%s`, env.MongoUser, env.MongoPassword, env.MongoHostPort)
		env.MysqlConnStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/stock_market?parseTime=true", env.MysqlUser, env.MysqlPassword, env.MysqlHost, env.MysqlPort)
	})
	return &env
}

func ChinaLoc() *time.Location {
	return loc
}
