package model

import (
	"database/sql/driver"
	"time"

	"github.com/joexzh/ThsConcept/db"
	"github.com/joexzh/dbh"
)

type ComparableConcept[T any] interface {
	GetId() string
	Cmp(old T) bool
}

type ConceptStockFt struct {
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

func (t *ConceptStockFt) Args() []any {
	return []any{
		&t.StockCode,
		&t.StockName,
		&t.UpdatedAt,
		&t.Description,
		&t.ConceptId,
		&t.ConceptName,
		&t.ConceptPlateId,
		&t.ConceptDefine,
		&t.ConceptUpdatedAt,
	}
}
func (t *ConceptStockFt) TableName() string {
	return "concept_stock_ft"
}
func (t *ConceptStockFt) Columns() []string {
	return []string{"stock_code", "stock_name", "updated_at", "description",
		"concept_id", "concept_name", "concept_plate_id", "concept_define", "concept_updated_at"}
}
func (t *ConceptStockFt) Config() *dbh.Config {
	return dbh.DefaultConfig
}

func (t *ConceptStockFt) Scan(src interface{}) error {
	return db.JsonScan(t, src)
}
func (t ConceptStockFt) Value() (driver.Value, error) {
	return db.JsonValue(t)
}

type ConceptStock struct {
	StockCode   string    `json:"stockCode"`
	StockName   string    `json:"stockName"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	ConceptId   string    `json:"conceptId"`
}

func (s *ConceptStock) GetId() string {
	return s.StockCode
}

func (s *ConceptStock) Cmp(old *ConceptStock) bool {
	return s.StockName == old.StockName &&
		s.Description == old.Description
}

func (s *ConceptStock) Args() []any {
	return []any{&s.StockCode, &s.StockName, &s.Description, &s.UpdatedAt, &s.ConceptId}
}
func (s *ConceptStock) TableName() string {
	return "concept_stock"
}
func (s *ConceptStock) Columns() []string {
	return []string{"stock_code", "stock_name", "description", "updated_at", "concept_id"}
}
func (s *ConceptStock) Config() *dbh.Config {
	return dbh.DefaultConfig
}

type ConceptStockSortByCode []*ConceptStock

func (s ConceptStockSortByCode) Len() int {
	return len(s)
}
func (s ConceptStockSortByCode) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ConceptStockSortByCode) Less(i, j int) bool {
	return s[i].StockCode < s[j].StockCode
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

func (c *Concept) GetId() string {
	return c.Id
}

func (c *Concept) Cmp(o *Concept) bool {
	return c.Name == o.Name &&
		c.PlateId == o.PlateId &&
		c.Define == o.Define &&
		c.UpdatedAt == o.UpdatedAt
}

type ConceptSortById []*Concept

func (c ConceptSortById) Len() int {
	return len(c)
}
func (c ConceptSortById) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c ConceptSortById) Less(i, j int) bool {
	return c[i].Id < c[j].Id
}

const (
	InsertConcept = iota + 1
	UpdateConcept
	DeleteConcept
	InsertConceptStock
	UpdateConceptStock
	DeleteConceptStock
)

type ConceptFtCommand struct {
	Id      int
	Command int
	Obj     ConceptStockFt
}

func (c *ConceptFtCommand) Args() []any {
	return []any{&c.Id, &c.Command, &c.Obj}
}
func (c *ConceptFtCommand) TableName() string {
	return "concept_stock_ft_sync"
}
func (c *ConceptFtCommand) Columns() []string {
	return []string{"id", "command", "obj"}
}
func (c *ConceptFtCommand) Config() *dbh.Config {
	return dbh.DefaultConfig
}
