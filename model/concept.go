package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Concept struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ConceptId    string             `bson:"conceptId,omitempty" json:"conceptId"`
	ConceptName  string             `bson:"conceptName,omitempty" json:"conceptName"`
	PlateId      int                `bson:"plateId,omitempty" json:"plateId"` // 软件中的板块id
	Define       string             `bson:"define,omitempty" json:"define"`   // 概念定义
	ReportDate   int64              `bson:"reportDate,omitempty" json:"reportDate"`
	Stocks       []ConceptStock     `bson:"stocks,omitempty" json:"stocks"`
	LastModified int64              `bson:"lastModified,omitempty" json:"lastModified"`
}

func (c *Concept) SetLastModifiedNow() {
	c.LastModified = time.Now().Unix()
}

func (c *Concept) Compare(new *Concept) bool {

	if c.ConceptId != new.ConceptId {
		return false
	}
	if c.ConceptName != new.ConceptName {
		return false
	}
	if c.PlateId != new.PlateId {
		return false
	}
	if c.Define != new.Define {
		return false
	}
	if c.ReportDate != new.ReportDate {
		return false
	}
	if len(c.Stocks) != len(new.Stocks) {
		return false
	}
	for _, stock := range c.Stocks {
		if !stock.Compare(new.Stocks) {
			return false
		}
	}
	return true
}

type ConceptStock struct {
	StockCode   string `bson:"stockCode,omitempty" json:"stockCode"`
	StockName   string `bson:"stockName,omitempty" json:"stockName"`
	Description string `bson:"description,omitempty" json:"description"`
}

func (s *ConceptStock) Compare(newStocks []ConceptStock) bool {
	for _, newSt := range newStocks {
		if newSt.StockCode == s.StockCode {
			if newSt.Description == s.Description {
				return true
			}
		}
	}
	return false
}

type StockConcept struct {
	Id           string `bson:"_id,omitempty" json:"id"` // stockCode + conceptId combine
	StockCode    string `bson:"stockCode,omitempty" json:"stockCode"`
	StockName    string `bson:"stockName,omitempty" json:"stockName"`
	ConceptId    string `bson:"conceptId,omitempty" json:"conceptId"`
	ConceptName  string `bson:"conceptName,omitempty" json:"conceptName"`
	Description  string `bson:"description,omitempty" json:"description"`
	LastModified int64  `bson:"lastModified,omitempty" json:"lastModified"`
}

func NewStockConcept(stock ConceptStock, conceptId string, conceptName string) *StockConcept {
	return &StockConcept{
		Id:           stock.StockCode + conceptId,
		StockCode:    stock.StockCode,
		StockName:    stock.StockName,
		ConceptId:    conceptId,
		ConceptName:  conceptName,
		Description:  stock.Description,
		LastModified: time.Now().Unix(),
	}
}

func (sc *StockConcept) Compare(old *StockConcept) bool {
	if sc.Description != old.Description {
		return false
	}
	return true
}
