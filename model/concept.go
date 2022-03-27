package model

import (
	"time"
)

type ConceptStock struct {
	StockCode         string    `json:"stockCode"`
	StockName         string    `json:"stockName"`
	PinyinFirstLetter string    `json:"pinyinFirstLetter"`
	PinyinNormal      string    `json:"pinyinNormal"`
	ConceptId         string    `json:"conceptId"`
	ConceptName       string    `json:"conceptName"`
	Description       string    `json:"description"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func (s *ConceptStock) CmpConcept(o *ConceptStock) bool {
	return s.ConceptId == o.ConceptId &&
		s.ConceptName == o.ConceptName &&
		s.Description == o.Description
}

func (s *ConceptStock) CmpStock(o *ConceptStock) bool {
	return s.StockCode == o.StockCode &&
		s.StockName == o.StockName &&
		s.PinyinFirstLetter == o.PinyinFirstLetter &&
		s.PinyinNormal == o.PinyinNormal
}

type ConceptStockByUpdateAtDesc []*ConceptStock

func (b ConceptStockByUpdateAtDesc) Len() int { return len(b) }
func (b ConceptStockByUpdateAtDesc) Less(i, j int) bool {
	return b[i].UpdatedAt.After(b[j].UpdatedAt)
}
func (b ConceptStockByUpdateAtDesc) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

type Concept struct {
	Id                string    `json:"id"`
	Name              string    `json:"name"`
	PinyinFirstLetter string    `json:"pinyinFirstLetter"`
	PinyinNormal      string    `json:"pinyinNormal"`
	PlateId           int       `json:"plateId"`
	Define            string    `json:"define"`
	UpdatedAt         time.Time `json:"updatedAt"`

	Stocks []*ConceptStock `json:"stocks"`
}

func (c *Concept) Cmp(o *Concept) bool {
	return c.Id == o.Id &&
		c.Name == o.Name &&
		c.PinyinFirstLetter == o.PinyinFirstLetter &&
		c.PinyinNormal == o.PinyinNormal &&
		c.PlateId == o.PlateId &&
		c.Define == o.Define
}

type ConceptByUpdateAtDesc []*Concept

func (b ConceptByUpdateAtDesc) Len() int { return len(b) }
func (b ConceptByUpdateAtDesc) Less(i, j int) bool {
	return b[i].UpdatedAt.After(b[j].UpdatedAt)
}
func (b ConceptByUpdateAtDesc) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
