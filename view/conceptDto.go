package view

import (
	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/model"
	"time"
)

type StockConceptDto struct {
	Id           string `bson:"_id,omitempty" json:"id"` // stockCode + conceptId combine
	StockCode    string `bson:"stockCode,omitempty" json:"stockCode"`
	StockName    string `bson:"stockName,omitempty" json:"stockName"`
	ConceptId    string `bson:"conceptId,omitempty" json:"conceptId"`
	ConceptName  string `bson:"conceptName,omitempty" json:"conceptName"`
	Description  string `bson:"description,omitempty" json:"description"`
	LastModified string `bson:"lastModified,omitempty" json:"lastModified"`
}

func NewStockConceptDto(sc *model.StockConcept, loc *time.Location) *StockConceptDto {
	t := time.Unix(sc.LastModified, 0)
	if loc != nil {
		t = t.In(loc)
	}
	return &StockConceptDto{
		Id:           sc.Id,
		StockCode:    sc.StockCode,
		StockName:    sc.StockName,
		ConceptId:    sc.ConceptId,
		ConceptName:  sc.ConceptName,
		Description:  sc.Description,
		LastModified: t.Format(config.TimeLayoutDate),
	}
}

func ScToScDto(scs ...model.StockConcept) []StockConceptDto {
	dtos := make([]StockConceptDto, 0, len(scs))
	loc, _ := time.LoadLocation(config.TimeLocal)
	for _, sc := range scs {
		dtos = append(dtos, *NewStockConceptDto(&sc, loc))
	}
	return dtos
}

type ScPageDto struct {
	Concept   string            `json:"concept"`
	StockName string            `json:"stockName"`
	Scs       []StockConceptDto `json:"scs"`
}
