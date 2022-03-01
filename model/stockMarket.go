package model

import "time"

// ZDTHistory 涨跌停历史
type ZDTHistory struct {
	Date            time.Time `json:"date"`            // 日期
	LongLimitCount  uint16    `json:"longLimitCount"`  // 涨停只数
	ShortLimitCount uint16    `json:"shortLimitCount"` // 跌停只数
	StopTradeCount  uint16    `json:"stopTradeCount"`  // 停牌只数
	Amount          float64   `json:"amount"`          // 两市交易(亿)
	SHLongCount     uint16    `json:"shLongCount"`     // 沪市上涨只数
	SHEvenCount     uint16    `json:"shEvenCount"`     // 沪市平盘只数
	SHShortCount    uint16    `json:"shShortCount"`    // 沪市下跌只数
	SZLongCount     uint16    `json:"szLongCount"`     // 深市上涨只数
	SZEvenCount     uint16    `json:"szEvenCount"`     // 深市平盘只数
	SZShortCount    uint16    `json:"szShortCount"`    // 深市下跌只数
}
