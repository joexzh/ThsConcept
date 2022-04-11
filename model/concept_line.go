package model

import (
	"time"

	"github.com/joexzh/dbh"
)

type ConceptLineDatePctChgOrderedView struct {
	Date  time.Time              `json:"date"`
	Lines []*ConceptLineWithName `json:"lines"`
}

type ConceptLineWithName struct {
	ConceptLine
	ConceptName string `json:"conceptName"`
}

func (c *ConceptLineWithName) Args() []any {
	return append(c.ConceptLine.Args(), &c.ConceptName)
}

type ConceptLine struct {
	PlateId string    `json:"plateId"`
	Date    time.Time `json:"date"`
	Open    float64   `json:"open"`
	High    float64   `json:"high"`
	Low     float64   `json:"low"`
	Close   float64   `json:"close"`
	PctChg  float64   `json:"pctChg"` // 涨跌幅%
	Volume  int       `json:"volume"` // 成交量
	Amount  float64   `json:"amount"` // 成交额
}

func (c *ConceptLine) Args() []any {
	return []any{&c.PlateId, &c.Date, &c.Open, &c.High, &c.Low, &c.Close, &c.PctChg, &c.Volume, &c.Amount}
}
func (c *ConceptLine) Columns() []string {
	return []string{"plate_id", "date", "open", "high", "low", "close", "pct_chg", "volume", "amount"}
}
func (c *ConceptLine) TableName() string {
	return "concept_line"
}
func (c *ConceptLine) Config() *dbh.Config {
	return dbh.DefaultConfig
}

type ConceptLineSortedByDate []*ConceptLine

func (c ConceptLineSortedByDate) Len() int {
	return len(c)
}
func (c ConceptLineSortedByDate) Less(i int, j int) bool {
	return c[i].Date.Before(c[j].Date)
}
func (c ConceptLineSortedByDate) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}
