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

	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"

	RexStockSymbol = `~(?:[a-z]{2})?([0-9]+)`
	StockSymbolUrl = "http://www.shdjt.com/js/lib/astock.js"

	ConceptPageUrl       = "http://basic.10jqka.com.cn/%v/concept.html"
	ConceptAllUrl        = "http://q.10jqka.com.cn/gn/"
	ConceptDetailPageUrl = "http://q.10jqka.com.cn/gn/detail/code/%v/"
	RexValidConceptPage  = `<h1 style="margin:3px 0px 0px 0px">\s*\d{6}\s*</h1>`
	RexConceptId         = `cid="(\d*)"`
	RexConceptCode       = `http://q.10jqka.com.cn/gn/detail/code/(\d*)`
	RexConceptDefine     = `<h4>[^\x00-\xff]+</h4>\s*<p style="overflow-y:auto;">([^<>]+)</p>`
	ConceptApiUrl        = "http://basic.10jqka.com.cn/ajax/stock/conceptlist.php?cid=%v&code=601127"

	TimeLayoutDate = "2006-01-02"
	TimeLayoutHour = "2006-01-02 15:04:05"
	TimeLocal      = "Asia/Shanghai"

	// Throttle 并发 goroutine 的数量, 防止系统或网络崩溃, 不能为0, 死锁!
	Throttle = 3
	// SleepRandUpTo 每个goroutine随机间隔的最大毫秒, 最低是0
	SleepRandUpTo = 500
)

var env *Env
var once sync.Once

type Env struct {
	MongoUser     string
	MongoPassword string
	MongoHostPort string
	MongoConnStr  string
}

func GetEnv() *Env {
	once.Do(func() {
		var (
			user     = os.Getenv("MONGO_USER")
			password = os.Getenv("MONGO_PASSWORD")
			hostPort = os.Getenv("MONGO_HOST_PORT")
		)
		env = &Env{
			MongoUser:     user,
			MongoPassword: password,
			MongoHostPort: hostPort,
			MongoConnStr:  fmt.Sprintf(`mongodb://%v:%v@%v`, user, password, hostPort),
		}
	})

	return env
}
