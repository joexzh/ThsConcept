//go:build wireinject

package repos

import (
	"github.com/google/wire"
	"github.com/joexzh/ThsConcept/db"
)

var set = wire.NewSet(NewStockMarketRepo, db.NewDB, db.NewMysqlConfig)

func InitStockMarketRepo() (*StockMarketRepo, error) {
	panic(wire.Build(set))
}
