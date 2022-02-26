package config

import (
	"fmt"
	"os"
	"sync"
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

var env = Env{
	MongoUser:     os.Getenv("MONGO_USER"),
	MongoPassword: os.Getenv("MONGO_PASSWORD"),
	MongoHostPort: os.Getenv("MONGO_HOST_PORT"),

	ServerPort: os.Getenv("SERVER_PORT"),

	MysqlUser:     os.Getenv("mysql_user"),
	MysqlPassword: os.Getenv("mysql_password"),
	MysqlHost:     os.Getenv("mysql_host"),
	MysqlPort:     os.Getenv("mysql_port"),
}

var once sync.Once

func GetEnv() Env {
	once.Do(func() {
		env.MongoConnStr = fmt.Sprintf(`mongodb://%s:%s@%s`, env.MongoUser, env.MongoPassword, env.MongoHostPort)
		env.MysqlConnStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/stock_market?parseTime=true&loc=Asia%%2FShanghai", env.MysqlUser, env.MysqlPassword, env.MysqlHost, env.MysqlPort)
	})

	return env
}
