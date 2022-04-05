package dto

import (
	"log"
	"time"

	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/model"
)

type ConceptsDto struct {
	Concepts []*model.Concept        `json:"concepts"`
	Stocks   []*model.ConceptStockFt `json:"stocks"`
}

// ConceptListApiReturn is a mapping to ths api result
type ConceptListApiReturn struct {
	ConceptId string
	Errorcode int                  `json:"errorode"`
	Errormsg  string               `json:"errormsg"`
	Result    ConceptListApiResult `json:"result"`
}

type ConceptListApiResult struct {
	Report   string `json:"report"`
	Name     string `json:"name"`
	Plateid  int    `json:"plateid"`
	Define   string
	Listdata map[string][][]interface{} `json:"listdata"`
}

// ConvertToConcept 原始格式参考 http://rap2.taobao.org/repository/editor?id=284321&mod=459202&itf=1980737 或 概念对比.json
func (ret *ConceptListApiReturn) ConvertToConcept() (*model.Concept, error) {
	t, err := time.ParseInLocation(config.TimeLayoutDate, ret.Result.Report, config.ChinaLoc())
	now := time.Now()
	if err != nil {
		return nil, err
	}

	concept := &model.Concept{
		Id:        ret.ConceptId,
		Name:      ret.Result.Name,
		PlateId:   ret.Result.Plateid,
		Define:    ret.Result.Define,
		UpdatedAt: t,
		Stocks:    make([]*model.ConceptStock, 0),
	}
	if concept.Name == "" {
		log.Printf("Concept name is empty, condept_id: %s\n", ret.ConceptId)
	}
	if concept.Define == "" {
		log.Printf("Concept define is empty, condept_id: %s\n", ret.ConceptId)
	}

	for _, v := range ret.Result.Listdata {
		for _, arr := range v {
			stock := &model.ConceptStock{
				StockCode:   arr[0].(string),
				StockName:   arr[1].(string),
				ConceptId:   ret.ConceptId,
				Description: arr[8].(string),
				UpdatedAt:   now,
			}
			concept.Stocks = append(concept.Stocks, stock)
			if stock.StockName == "" {
				log.Printf("Concept stock name is empty, concept_id: %s, code: %s\n", ret.ConceptId, stock.StockCode)
			}
			if stock.Description == "" {
				log.Printf("Concept stock description is empty, concept_id: %s, code: %s\n", ret.ConceptId, stock.StockCode)
			}
		}
	}
	return concept, nil
}
