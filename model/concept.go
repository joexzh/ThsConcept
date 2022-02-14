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
	StockCode   string `bson:"stockCode,omitempty" json:"stockCode"`
	StockName   string `bson:"stockName,omitempty" json:"stockName"`
	Description string `bson:"description,omitempty" json:"description"`
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

// ConceptShortMap key for short, value for model name
var ConceptShortMap = map[string]string{
	"EDR":     "EDR概念",
	"三胎概念":    "三胎概念",
	"仿制药":     "仿制药一致性评价",
	"分拆上市":    "分拆上市意愿",
	"氟化工":     "氟化工概念",
	"富士康":     "富士康概念",
	"共同富裕":    "共同富裕示范区",
	"光伏":      "光伏概念",
	"白酒":      "白酒概念",
	"阿里巴巴":    "阿里巴巴概念",
	"北交所":     "北交所概念",
	"百度":      "百度概念",
	"冰雪":      "冰雪产业",
	"宁德时代概念":  "宁德时代概念",
	"创业板重组":   "创业板重组松绑",
	"代糖":      "代糖概念",
	"国家大基金":   "国家大基金持股",
	"大基金":     "国家大基金持股",
	"航运":      "航运概念",
	"鸿蒙":      "鸿蒙概念",
	"换电":      "换电概念",
	"黄金":      "黄金概念",
	"海思":      "华为海思",
	"集成电路":    "集成电路概念",
	"华为":      "华为概念",
	"今日头条":    "今日头条概念",
	"机器人":     "机器人概念",
	"快手":      "快手概念",
	"蚂蚁金服":    "蚂蚁金服概念",
	"煤炭":      "煤炭概念",
	"NFT":     "NFT概念",
	"独角兽":     "独角兽概念",
	"赛马":      "赛马概念",
	"啤酒":      "啤酒概念",
	"拼多多":     "拼多多概念",
	"MSCI":    "MSCI概念",
	"苹果":      "苹果概念",
	"汽车拆解":    "汽车拆解概念",
	"PPP":     "PPP概念",
	"期货":      "期货概念",
	"柔性直流":    "柔性直流输电",
	"上海国资":    "上海国资改革",
	"水泥":      "水泥概念",
	"钛白粉":     "钛白粉概念",
	"深圳国资":    "深圳国资改革",
	"特钢":      "特钢概念",
	"腾讯":      "腾讯概念",
	"消费电子":    "消费电子概念",
	"小米":      "小米概念",
	"小金属":     "小金属概念",
	"新材料":     "新材料概念",
	"新三板精选层":  "新三板精选层概念",
	"信托":      "信托概念",
	"养老":      "养老概念",
	"央企国资":    "央企国资改革",
	"医疗废物":    "医疗废物处理",
	"医美":      "医美概念",
	"有机硅":     "有机硅概念",
	"医疗器械":    "医疗器械概念",
	"中芯国际":    "中芯国际概念",
	"知识产权":    "知识产权保护",
	"中字头":     "中字头股票",
	"注册制次新":   "注册制次新股",
	"证金":      "证金持股",
	"足球":      "足球概念",
	"DRG/DIP": "DRG/DIP概念",
}
