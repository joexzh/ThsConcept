package realtime

import (
	"sort"
	"strings"
)

type Dto struct {
	List          []Message      `json:"list"`
	Filter        []Tag          `json:"filter"`
	Total         string         `json:"total"`
	KeywordCounts []KeywordCount `json:"keywordCounts"`
}

var ConceptShortsMap = map[string][]string{
	"EDR概念":     {"EDR"},
	"三胎概念":      {"三胎"},
	"仿制药一致性评价":  {"仿制药", "仿制药一致性"},
	"分拆上市意愿":    {"分拆上市"},
	"氟化工概念":     {"氟化工"},
	"富士康概念":     {"富士康"},
	"共同富裕示范区":   {"共同富裕"},
	"光伏概念":      {"光伏"},
	"白酒概念":      {"白酒"},
	"阿里巴巴概念":    {"阿里巴巴"},
	"北交所概念":     {"北交所"},
	"百度概念":      {"百度"},
	"冰雪产业":      {"冰雪"},
	"宁德时代概念":    {"宁德时代"},
	"创业板重组松绑":   {"创业板重组"},
	"代糖概念":      {"代糖"},
	"国家大基金持股":   {"大基金"},
	"航运概念":      {"航运"},
	"鸿蒙概念":      {"鸿蒙"},
	"换电概念":      {"换电"},
	"黄金概念":      {"黄金"},
	"华为海思":      {"海思"},
	"集成电路概念":    {"集成电路"},
	"快手概念":      {"快手"},
	"蚂蚁金服概念":    {"蚂蚁金服"},
	"NFT概念":     {"NFT"},
	"独角兽概念":     {"独角兽"},
	"赛马概念":      {"赛马"},
	"啤酒概念":      {"啤酒"},
	"拼多多概念":     {"拼多多"},
	"MSCI概念":    {"MSCI"},
	"苹果概念":      {"苹果"},
	"汽车拆解概念":    {"汽车拆解"},
	"PPP概念":     {"PPP"},
	"期货概念":      {"期货"},
	"柔性直流输电":    {"柔性直流"},
	"上海国资改革":    {"上海国资"},
	"水泥概念":      {"水泥"},
	"钛白粉概念":     {"钛白粉"},
	"深圳国资改革":    {"深圳国资"},
	"特钢概念":      {"特钢"},
	"腾讯概念":      {"腾讯"},
	"消费电子概念":    {"消费电子"},
	"小米概念":      {"小米"},
	"小金属概念":     {"小金属"},
	"新材料概念":     {"新材料"},
	"新三板精选层概念":  {"新三板精选层"},
	"信托概念":      {"信托"},
	"养老概念":      {"养老"},
	"央企国资改革":    {"央企国资"},
	"医疗废物处理":    {"医疗废物"},
	"医美概念":      {"医美"},
	"有机硅概念":     {"有机硅"},
	"医疗器械概念":    {"医疗器械"},
	"中芯国际概念":    {"中芯国际"},
	"知识产权保护":    {"知识产权"},
	"中字头股票":     {"中字头"},
	"注册制次新股":    {"注册制次新"},
	"证金持股":      {"证金"},
	"足球概念":      {"足球"},
	"DRG/DIP概念": {"DRG/DIP"},
}

type KeywordCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func KeywordCounts(str string, conceptShortsMap map[string][]string) []KeywordCount {
	matchDict := make(map[string]int)

	for k, v := range conceptShortsMap {
		if matchKeywords(str, v) {
			if _, ok := matchDict[k]; ok {
				matchDict[k] += 1
			} else {
				matchDict[k] = 1
			}
			// match one concept is enough
			break
		}
	}
	kwCounts := make([]KeywordCount, 0, len(matchDict))
	for k, v := range matchDict {
		kwCounts = append(kwCounts, KeywordCount{Name: k, Count: v})
	}
	return kwCounts
}

func MergeConceptShortsMap(conceptNames []string) map[string][]string {
	dict := make(map[string][]string, len(conceptNames))

	for _, name := range conceptNames {
		var value []string
		if v, ok := ConceptShortsMap[name]; ok {
			value = append(value, v...)
		}
		dict[name] = value
	}
	return dict
}

func matchKeywords(str string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(str, kw) {
			return true
		}
	}
	return false
}

func SortKeywordCounts(kwc []KeywordCount) {
	sort.Slice(kwc, func(i, j int) bool {
		return kwc[i].Count > kwc[j].Count
	})
}
