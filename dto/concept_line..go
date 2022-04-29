package dto

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/model"
)

type ConceptLine struct {
	Rt  string `json:"rt"`
	Num int    `json:"num"`
	// Total      string         `json:"total"` // can be int or string
	Start      string         `json:"start"`
	Year       map[string]int `json:"year"`
	Name       string         `json:"name"`
	Data       string         `json:"data"`
	MarkType   string         `json:"markType"`
	IssuePrice string         `json:"issuePrice"`
	Today      string         `json:"today"`
}

// Convert2 api result to []*model.ConceptLine.
// If ConceptLine.Today is after last day in ConceptLine.Data, the last day is unreliable, we should fetch it from today's api
func (c *ConceptLine) Convert2(plateId string) ([]*model.ConceptLine, float64, bool, error) {
	lines := make([]*model.ConceptLine, 0, c.Num)
	days := strings.Split(c.Data, ";")
	issuePrice, _ := strconv.ParseFloat(c.IssuePrice, 64)
	startDate, err := time.ParseInLocation("20060102", c.Start, config.ChinaLoc())
	if err != nil {
		return nil, 0, false, fmt.Errorf("dto.ConceptLine.Convert2: parse start \"%s\" err: %v", c.Today, err)
	}

	prevClose := issuePrice
	var shouldFetchToday bool
	for i := range days {
		line, valid, err := parseToConceptLine(plateId, strings.Split(days[i], ","), prevClose)
		if err != nil {
			if i < len(days)-1 {
				return nil, 0, false, fmt.Errorf("dto.ConceptLine.Convert2: plateId %s, failed to parse ConceptLine \"%s\", %v", plateId, days[i], err)
			}
			// if the last day is not complete yet, exclude it to the result list
			log.Printf(`dto.ConceptLine.Convert2: plateId %s, the last day "%s" failed to parse, %v`, plateId, days[i], err)
			shouldFetchToday = true
			break
		}
		if i == len(days)-1 && !valid {
			shouldFetchToday = true
			break
		}
		prevClose = line.Close
		if i == 0 {
			if issuePrice > 0 && line.Date.Equal(startDate) {
				lines = append(lines, line)
			}
		} else {
			lines = append(lines, line)
		}
	}

	return lines, prevClose, shouldFetchToday, nil
}

func parseToConceptLine(plateId string, dailyDataArr []string, prevClose float64) (*model.ConceptLine, bool, error) {
	var valid bool
	if len(dailyDataArr) < 7 {
		return nil, valid, errors.New("len(dailyDataArr) < 7")
	}
	if len(dailyDataArr) > 7 {
		if len(dailyDataArr) == 12 {
			valid = true
		}
	}

	date, err := time.ParseInLocation("20060102", dailyDataArr[0], config.ChinaLoc())
	if err != nil {
		return nil, valid, fmt.Errorf(`failed to parse date "%s", %s`, dailyDataArr[0], err.Error())
	}
	open, err := strconv.ParseFloat(dailyDataArr[1], 10)
	if err != nil {
		return nil, valid, fmt.Errorf(`failed to parse open "%s", %s`, dailyDataArr[1], err.Error())
	}
	high, err := strconv.ParseFloat(dailyDataArr[2], 10)
	if err != nil {
		return nil, valid, fmt.Errorf(`failed to parse high "%s", %s`, dailyDataArr[2], err.Error())
	}
	low, err := strconv.ParseFloat(dailyDataArr[3], 10)
	if err != nil {
		return nil, valid, fmt.Errorf(`failed to parse low "%s", %s`, dailyDataArr[3], err.Error())
	}
	klose, err := strconv.ParseFloat(dailyDataArr[4], 10)
	if err != nil {
		return nil, valid, fmt.Errorf(`failed to parse close "%s", %s`, dailyDataArr[4], err.Error())
	}
	volume, err := strconv.Atoi(dailyDataArr[5])
	if err != nil {
		return nil, valid, fmt.Errorf(`failed to parse volume "%s", %s`, dailyDataArr[5], err.Error())
	}
	amount, err := strconv.ParseFloat(dailyDataArr[6], 10)
	if err != nil {
		return nil, valid, fmt.Errorf(`failed to parse amount "%s", %s`, dailyDataArr[6], err.Error())
	}

	return &model.ConceptLine{
		PlateId: plateId,
		Date:    date,
		Open:    open,
		High:    high,
		Low:     low,
		Close:   klose,
		PctChg:  (klose - prevClose) / prevClose,
		Volume:  volume,
		Amount:  amount,
	}, valid, nil
}

/*
ConceptLineToday example:

{
    "bk_885978": {
        "1": "20220413",
        "7": "1025.235",
        "8": "1025.235",
        "9": "992.358",
        "11": "992.438",
        "13": 651449050,
        "19": "4205627800.000",
        "74": "",
        "1968584": "",
        "66": "",
        "open": 1,
        "dt": "2348",
        "name": "低辐射玻璃（Low-E）",
        "marketType": ""
    }
}
*/
type ConceptLineToday map[string]map[string]interface{}

func (c ConceptLineToday) Convert2(plateId string, prevClose float64) (*model.ConceptLine, error) {
	today, ok := c["bk_"+plateId]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no today data")
	}
	var ss []string
	date, ok := today["1"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no date")
	}
	ss = append(ss, interfaceToString(date))

	open, ok := today["7"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no open")
	}
	ss = append(ss, interfaceToString(open))

	high, ok := today["8"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no high")
	}
	ss = append(ss, interfaceToString(high))

	low, ok := today["9"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no low")
	}
	ss = append(ss, interfaceToString(low))

	klose, ok := today["11"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no close")
	}
	ss = append(ss, interfaceToString(klose))

	volume, ok := today["13"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no volume")
	}
	ss = append(ss, interfaceToString(volume))

	amount, ok := today["19"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no amount")
	}
	ss = append(ss, interfaceToString(amount))

	line, _, err := parseToConceptLine(plateId, ss, prevClose)
	return line, err
}

func interfaceToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}
