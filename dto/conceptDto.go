package dto

import (
	"github.com/joexzh/ThsConcept/util"
	"time"

	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/model"
)

type ConceptsDto struct {
	Concepts []*model.Concept      `json:"concepts"`
	Stocks   []*model.ConceptStock `json:"stocks"`
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
		Id:                ret.ConceptId,
		Name:              ret.Result.Name,
		PinyinFirstLetter: util.Pinyin(ret.Result.Name, util.PinyinFirstLetterArgs),
		PinyinNormal:      util.Pinyin(ret.Result.Name, util.PinyinNormalArgs),
		PlateId:           ret.Result.Plateid,
		Define:            ret.Result.Define,
		UpdatedAt:         t,
		Stocks:            make([]*model.ConceptStock, 0),
	}

	for _, v := range ret.Result.Listdata {
		for _, arr := range v {
			stock := &model.ConceptStock{
				StockCode:         arr[0].(string),
				StockName:         arr[1].(string),
				PinyinFirstLetter: util.Pinyin(arr[1].(string), util.PinyinFirstLetterArgs),
				PinyinNormal:      util.Pinyin(arr[1].(string), util.PinyinNormalArgs),
				ConceptId:         ret.ConceptId,
				ConceptName:       ret.Result.Name,
				Description:       arr[8].(string),
				UpdatedAt:         now,
			}
			concept.Stocks = append(concept.Stocks, stock)
		}
	}
	return concept, nil
}
