package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/joexzh/ThsConcept/dto"
	"github.com/joexzh/ThsConcept/model"
	"github.com/joexzh/ThsConcept/util"
	"github.com/pkg/errors"
)

const (
	ConceptLineUrl      = "http://d.10jqka.com.cn/v4/line/bk_%s/01/last.js"
	ConceptLineTodayUrl = "http://d.10jqka.com.cn/v4/line/bk_%s/01/today.js"
)

func ConceptLine(ctx context.Context, plateId string) (*dto.ConceptLine, error) {
	headers := map[string]string{
		"Referer": "http://q.10jqka.com.cn/",
	}
	resp, err := util.HttpGet(ctx, fmt.Sprintf(ConceptLineUrl, plateId), headers, nil)
	if err != nil {
		return nil, errors.Wrap(err, "fetch.ConceptLine: plateId="+plateId)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "fetch.ConceptLine: plateId="+plateId)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("fetch.ConceptLine: plate_id %s return 0 byte", plateId)
	}
	// remove jsonp wrap
	data = data[len(fmt.Sprintf("quotebridge_v4_line_bk_%s_01_last(", plateId)) : len(data)-1]
	var conceptLineDto dto.ConceptLine
	if err = json.Unmarshal(data, &conceptLineDto); err != nil {
		return nil, errors.Wrap(err, "fetch.ConceptLine: plateId="+plateId)
	}

	return &conceptLineDto, err
}

// ConceptLineToday fetch today's conceptLine
func ConceptLineToday(ctx context.Context, plateId string, prevClose float64) (*model.ConceptLine, error) {
	headers := map[string]string{
		"Referer": "http://q.10jqka.com.cn/",
	}
	resp, err := util.HttpGet(ctx, fmt.Sprintf(ConceptLineTodayUrl, plateId), headers, nil)
	if err != nil {
		return nil, errors.New("ConceptLineToday: " + err.Error())
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("ConceptLineToday: " + err.Error())
	}
	if len(data) == 0 {
		return nil, errors.New("ConceptLineToday: len(data) == 0")
	}
	// remove jsonp wrap
	data = data[len(fmt.Sprintf("quotebridge_v4_line_bk_%s_01_today(", plateId)) : len(data)-1]
	lineDto := make(dto.ConceptLineToday)
	if err = json.Unmarshal(data, &lineDto); err != nil {
		return nil, errors.New("ConceptLineToday: " + err.Error())
	}
	return lineDto.Convert2(plateId, prevClose)
}
