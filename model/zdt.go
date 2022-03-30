package model

import "time"

// ZDTHistory 涨跌停历史
type ZDTHistory struct {
	Date            time.Time `json:"date" db:"date"`                         // 日期
	LongLimitCount  uint16    `json:"longLimitCount" db:"long_limit_count"`   // 涨停只数
	ShortLimitCount uint16    `json:"shortLimitCount" db:"short_limit_count"` // 跌停只数
	StopTradeCount  uint16    `json:"stopTradeCount" db:"stop_trade_count"`   // 停牌只数
	Amount          float64   `json:"amount" db:"amount"`                     // 两市交易(亿)
	SHLongCount     uint16    `json:"shLongCount" db:"sh_long_count"`         // 沪市上涨只数
	SHEvenCount     uint16    `json:"shEvenCount" db:"sh_even_count"`         // 沪市平盘只数
	SHShortCount    uint16    `json:"shShortCount" db:"sh_short_count"`       // 沪市下跌只数
	SZLongCount     uint16    `json:"szLongCount" db:"sz_long_count"`         // 深市上涨只数
	SZEvenCount     uint16    `json:"szEvenCount" db:"sz_even_count"`         // 深市平盘只数
	SZShortCount    uint16    `json:"szShortCount" db:"sz_short_count"`       // 深市下跌只数
}

func (z *ZDTHistory) Args() []any {
	return []any{&z.Date, &z.LongLimitCount, &z.ShortLimitCount, &z.StopTradeCount, &z.Amount,
		&z.SHLongCount, &z.SHEvenCount, &z.SHShortCount, &z.SZLongCount, &z.SZEvenCount, &z.SZShortCount}
}
func (z *ZDTHistory) Columns() []string {
	return []string{"date",
		"long_limit_count",
		"short_limit_count",
		"stop_trade_count",
		"amount",
		"sh_long_count",
		"sh_even_count",
		"sh_short_count",
		"sz_long_count",
		"sz_even_count",
		"sz_short_count"}
}
func (z *ZDTHistory) TableName() string {
	return "long_short"
}
