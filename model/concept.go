package model

import (
	"fmt"
	"github.com/joexzh/ThsConcept/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

type Concept struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ConceptId    string             `bson:"conceptId,omitempty" json:"conceptId"`
	ConceptName  string             `bson:"conceptName,omitempty" json:"conceptName"`
	PlateId      int                `bson:"plateId,omitempty" json:"plateId"` // 软件中的板块id
	Define       string             `bson:"define,omitempty" json:"define"`   // 概念定义
	ReportDate   int64              `bson:"reportDate,omitempty" json:"reportDate"`
	Stocks       []Stock            `bson:"stocks,omitempty" json:"stocks"`
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

type Stock struct {
	StockCode   string  `bson:"stockCode,omitempty" json:"stockCode"`
	StockName   string  `bson:"stockName,omitempty" json:"stockName"`
	Description string  `bson:"description,omitempty" json:"description"`
}

func (s *Stock) Compare(newStocks []Stock) bool {
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

func NewStockConcept(stock Stock, conceptId string, conceptName string) *StockConcept {
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

// Return is a mapping to ths api result
type Return struct {
	ConceptId string
	Errorcode int          `json:"errorode"`
	Errormsg  string       `json:"errormsg"`
	Result    ReturnResult `json:"result"`
}

type ReturnResult struct {
	Report   string `json:"report"`
	Name     string `json:"name"`
	Plateid  int    `json:"plateid"`
	Define   string
	Listdata map[string][][]interface{} `json:"listdata"`
}

// ConvertToConcept 原始格式参考 http://rap2.taobao.org/repository/editor?id=284321&mod=459202&itf=1980737 或 概念对比.json
func (ret *Return) ConvertToConcept() (*Concept, error) {
	local, err := time.LoadLocation(config.TimeLocal)
	if err != nil {
		return nil, err
	}
	t, err := time.ParseInLocation(config.TimeLayoutDate, ret.Result.Report, local)
	if err != nil {
		return nil, err
	}

	concept := &Concept{
		ConceptId:   ret.ConceptId,
		ConceptName: ret.Result.Name,
		PlateId:     ret.Result.Plateid,
		Define:      ret.Result.Define,
		ReportDate:  t.Unix(),
	}

	for _, v := range ret.Result.Listdata {
		for _, arr := range v {
			stock := Stock{
				StockCode:   arr[0].(string),
				StockName:   arr[1].(string),
				Description: arr[8].(string),
			}

			//// 停牌有可能是"--", 需要跳过处理
			//
			//// 流通市值
			//s2 := strings.TrimRight(arr[2].(string), "亿")
			//stock.FlowCap = convert(s2)
			//
			//// 市值
			//s3 := strings.TrimRight(arr[3].(string), "亿")
			//stock.Cap = convert(s3)
			//
			//// 净利润
			//s4 := strings.TrimRight(arr[4].(string), "亿")
			//stock.NetProfit = convert(s4)
			//
			//// 涨跌幅%
			//stock.Change = convert(arr[5])
			//
			//// 现价
			//stock.Price = convert(arr[6])
			//
			//// 市盈率
			//stock.PE = convert(arr[12])

			concept.Stocks = append(concept.Stocks, stock)
		}
	}

	return concept, nil
}

func convert(value interface{}) float64 {
	switch v := value.(type) {
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			fmt.Printf("Error setting field, expected a number but get (%v) of type %t, set value to 0\n", value, value)
			return 0
		}
		return f
	case float64:
		return v
	default:
		fmt.Printf("Error setting field, expected a number but get (%v) of type %t, set value to 0\n", value, value)
		return 0
	}
}
