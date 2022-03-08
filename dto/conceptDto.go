package dto

import (
	"fmt"
	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/model"
	"strconv"
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
	local, err := time.LoadLocation(config.TimeLocal)
	if err != nil {
		return nil, err
	}
	t, err := time.ParseInLocation(config.TimeLayoutDate, ret.Result.Report, local)
	if err != nil {
		return nil, err
	}

	concept := &model.Concept{
		ConceptId:   ret.ConceptId,
		ConceptName: ret.Result.Name,
		PlateId:     ret.Result.Plateid,
		Define:      ret.Result.Define,
		ReportDate:  t.Unix(),
	}

	for _, v := range ret.Result.Listdata {
		for _, arr := range v {
			stock := model.ConceptStock{
				StockCode:   arr[0].(string),
				StockName:   arr[1].(string),
				Description: arr[8].(string),
			}

			// // 停牌有可能是"--", 需要跳过处理
			//
			// // 流通市值
			// s2 := strings.TrimRight(arr[2].(string), "亿")
			// stock.FlowCap = convert(s2)
			//
			// // 市值
			// s3 := strings.TrimRight(arr[3].(string), "亿")
			// stock.Cap = convert(s3)
			//
			// // 净利润
			// s4 := strings.TrimRight(arr[4].(string), "亿")
			// stock.NetProfit = convert(s4)
			//
			// // 涨跌幅%
			// stock.Change = convert(arr[5])
			//
			// // 现价
			// stock.Price = convert(arr[6])
			//
			// // 市盈率
			// stock.PE = convert(arr[12])

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
