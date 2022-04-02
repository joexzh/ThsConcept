package model

import (
	"time"

	"github.com/joexzh/dbh"
)

type ConceptStock struct {
	StockCode        string    `json:"stockCode" db:"stock_code"`
	StockName        string    `json:"stockName" db:"stock_name"`
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`
	Description      string    `json:"description" db:"description"`
	ConceptId        string    `json:"conceptId" db:"concept_id"`
	ConceptName      string    `json:"conceptName" db:"concept_name"`
	ConceptPlateId   int       `json:"conceptPlateId" db:"concept_plate_id"`
	ConceptDefine    string    `json:"conceptDefine" db:"concept_define"`
	ConceptUpdatedAt time.Time `json:"conceptUpdatedAt" db:"concept_updated_at"`
}

func (s *ConceptStock) Args() []any {
	return []any{
		&s.StockCode,
		&s.StockName,
		&s.UpdatedAt,
		&s.Description,
		&s.ConceptId,
		&s.ConceptName,
		&s.ConceptPlateId,
		&s.ConceptDefine,
		&s.ConceptUpdatedAt,
	}
}

func (s *ConceptStock) CmpConcept(o *ConceptStock) bool {
	return s.ConceptId == o.ConceptId &&
		s.ConceptName == o.ConceptName &&
		s.Description == o.Description
}

func (s *ConceptStock) CmpStock(o *ConceptStock) bool {
	return s.StockCode == o.StockCode &&
		s.StockName == o.StockName
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
	Id        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	PlateId   int       `json:"plateId" db:"plate_id"`
	Define    string    `json:"define" db:"define"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`

	Stocks []*ConceptStock `json:"stocks"`
}

func (c *Concept) Args() []any {
	return []any{&c.Id, &c.Name, &c.PlateId, &c.Define, &c.UpdatedAt}
}
func (c *Concept) Columns() []string {
	return []string{"id", "name", "plate_id", "define", "updated_at"}
}
func (c *Concept) TableName() string {
	return "concept_concept"
}
func (c *Concept) Config() *dbh.Config {
	return dbh.DefaultConfig
}

func (c *Concept) Cmp(o *Concept) bool {
	return c.Id == o.Id &&
		c.Name == o.Name &&
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
