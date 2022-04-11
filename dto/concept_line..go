package dto

import (
	"strconv"
	"strings"
	"time"

	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/model"
	"github.com/pkg/errors"
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

func (c *ConceptLine) ConverTo(plateId string) ([]*model.ConceptLine, error) {
	lines := make([]*model.ConceptLine, 0, c.Num)

	for _, s := range strings.Split(c.Data, ";") {
		ss := strings.Split(s, ",")
		if len(ss) < 7 {
			continue
		}

		date, err := time.ParseInLocation("20060102", ss[0], config.ChinaLoc())
		if err != nil {
			return nil, errors.Wrap(err, "ConceptLine.ConverTo, date="+ss[0])
		}
		open, err := strconv.ParseFloat(ss[1], 10)
		if err != nil {
			return nil, errors.Wrap(err, "ConceptLine.ConverTo, open="+ss[1])
		}
		high, err := strconv.ParseFloat(ss[2], 10)
		if err != nil {
			return nil, errors.Wrap(err, "ConceptLine.ConverTo, high="+ss[2])
		}
		low, err := strconv.ParseFloat(ss[3], 10)
		if err != nil {
			return nil, errors.Wrap(err, "ConceptLine.ConverTo, low="+ss[3])
		}
		close, err := strconv.ParseFloat(ss[4], 10)
		if err != nil {
			return nil, errors.Wrap(err, "ConceptLine.ConverTo, close="+ss[4])
		}
		volume, err := strconv.Atoi(ss[5])
		if err != nil {
			return nil, errors.Wrap(err, "ConceptLine.ConverTo, volume="+ss[5])
		}
		amount, err := strconv.ParseFloat(ss[6], 10)
		if err != nil {
			return nil, errors.Wrap(err, "ConceptLine.ConverTo, amount="+ss[6])
		}

		lines = append(lines, &model.ConceptLine{
			PlateId: plateId,
			Date:    date,
			Open:    open,
			High:    high,
			Low:     low,
			Close:   close,
			PctChg:  (close - open) / open,
			Volume:  volume,
			Amount:  amount,
		})
	}

	return lines, nil
}
