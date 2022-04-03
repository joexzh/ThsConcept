package model

import (
	"time"

	"github.com/joexzh/dbh"
)

type ConceptStockView struct {
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

func (s *ConceptStockView) Args() []any {
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

type ConceptStock struct {
	StockCode   string    `json:"stockCode"`
	StockName   string    `json:"stockName"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description"`
	ConceptId   string    `json:"conceptId"`
}

func (s *ConceptStock) Cmp(old *ConceptStock) bool {
	return s.StockName == old.StockName &&
		s.Description == old.Description
}

func (s *ConceptStock) Args() []any {
	return []any{&s.StockCode, &s.StockName, &s.UpdatedAt, &s.Description, &s.ConceptId}
}
func (s *ConceptStock) TableName() string {
	return "concept_stock"
}
func (s *ConceptStock) Columns() []string {
	return []string{"stock_code", "stock_name", "updated_at", "description", "concept_id"}
}
func (s *ConceptStock) Config() *dbh.Config {
	return dbh.DefaultConfig
}

type Concept struct {
	Id        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	PlateId   int       `json:"plateId" db:"plate_id"`
	Define    string    `json:"define" db:"define"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

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
