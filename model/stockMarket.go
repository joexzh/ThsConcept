package model

import "time"

// ZDTHistory 涨跌停历史
type ZDTHistory struct {
	Date            time.Time // 日期
	LongLimitCount  uint16    // 涨停只数
	ShortLimitCount uint16    // 跌停只数
	StopTradeCount  uint16    // 停牌只数
	Amount          float64   // 两市交易(亿)
	SHLongCount     uint16    // 沪市上涨只数
	SHEvenCount     uint16    // 沪市平盘只数
	SHShortCount    uint16    // 沪市下跌只数
	SZLongCount     uint16    // 深市上涨只数
	SZEvenCount     uint16    // 深市平盘只数
	SZShortCount    uint16    // 深市下跌只数
}
