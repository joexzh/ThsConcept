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

func (c *ConceptLine) ConverTo(plateId string) ([]*model.ConceptLine, bool, error) {
	lines := make([]*model.ConceptLine, 0, c.Num)
	latestIncluded := true
	days := strings.Split(c.Data, ";")

	for i, s := range days {
		line, err := parseToConceptLine(plateId, strings.Split(s, ","))
		if err != nil {
			if i == len(days)-1 {
				latestIncluded = false
			}
			log.Println("ConceptLine.ConvertTo: ", err)
			continue
		}

		lines = append(lines, line)
	}

	return lines, latestIncluded, nil
}

func parseToConceptLine(plateId string, ss []string) (*model.ConceptLine, error) {
	if len(ss) < 7 {
		return nil, errors.New("dto.parseToConceptLine: len(ss) < 7")
	}

	date, err := time.ParseInLocation("20060102", ss[0], config.ChinaLoc())
	if err != nil {
		return nil, fmt.Errorf("dto.parseToConceptLine, plateId=%s, date=%s, err=%s\n", plateId, ss[0], err.Error())

	}
	open, err := strconv.ParseFloat(ss[1], 10)
	if err != nil {
		return nil, fmt.Errorf("dto.parseToConceptLine, plateId=%s, open=%s, err=%s\n", plateId, ss[1], err.Error())
	}
	high, err := strconv.ParseFloat(ss[2], 10)
	if err != nil {
		return nil, fmt.Errorf("dto.parseToConceptLine, plateId=%s, high=%s, err=%s\n", plateId, ss[2], err.Error())
	}
	low, err := strconv.ParseFloat(ss[3], 10)
	if err != nil {
		return nil, fmt.Errorf("dto.parseToConceptLine, plateId=%s, low=%s, err=%s\n", plateId, ss[3], err.Error())
	}
	close, err := strconv.ParseFloat(ss[4], 10)
	if err != nil {
		return nil, fmt.Errorf("dto.parseToConceptLine, plateId=%s, close=%s, err=%s\n", plateId, ss[4], err.Error())
	}
	volume, err := strconv.Atoi(ss[5])
	if err != nil {
		return nil, fmt.Errorf("dto.parseToConceptLine, plateId=%s, volume=%s, err=%s\n", plateId, ss[5], err.Error())
	}
	amount, err := strconv.ParseFloat(ss[6], 10)
	if err != nil {
		return nil, fmt.Errorf("dto.parseToConceptLine, plateId=%s, amount=%s, err=%s\n", plateId, ss[6], err.Error())
	}

	return &model.ConceptLine{
		PlateId: plateId,
		Date:    date,
		Open:    open,
		High:    high,
		Low:     low,
		Close:   close,
		PctChg:  (close - open) / open,
		Volume:  volume,
		Amount:  amount,
	}, nil
}

/*
example:

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
type ConceptLineToday map[string]map[string]string

func (c ConceptLineToday) ConverTo(plateId string) (*model.ConceptLine, error) {
	today, ok := c["bk_"+plateId]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no today data")
	}
	var ss []string
	date, ok := today["1"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no date")
	}
	ss = append(ss, date)

	open, ok := today["7"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no open")
	}
	ss = append(ss, open)

	high, ok := today["8"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no high")
	}
	ss = append(ss, high)

	low, ok := today["9"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no low")
	}
	ss = append(ss, low)

	close, ok := today["11"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no close")
	}
	ss = append(ss, close)

	volume, ok := today["13"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no volume")
	}
	ss = append(ss, volume)

	amount, ok := today["19"]
	if !ok {
		return nil, errors.New("dto.ConceptLineToday.ConvertTo: no amount")
	}
	ss = append(ss, amount)

	return parseToConceptLine(plateId, ss)
}
