package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/joexzh/ThsConcept/dto"
	"github.com/joexzh/ThsConcept/model"
	"github.com/joexzh/ThsConcept/util"
	"github.com/pkg/errors"
)

const (
	ConceptLineUrl      = "http://d.10jqka.com.cn/v4/line/bk_%s/01/last.js"
	ConceptLineTodayUrl = "http://d.10jqka.com.cn/v4/line/bk_%s/01/today.js"
)

func FetchConceptLine(ctx context.Context, plateId string) ([]*model.ConceptLine, error) {
	headers := map[string]string{
		"Referer": "http://q.10jqka.com.cn/",
	}
	resp, err := util.HttpGet(ctx, fmt.Sprintf(ConceptLineUrl, plateId), headers, nil)
	if err != nil {
		return nil, errors.Wrap(err, "FetchConceptLine, plateId="+plateId)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "FetchConceptLine, plateId="+plateId)
	}
	if len(data) == 0 {
		log.Printf("concept_line: plate_id %s return 0 byte", plateId)
		return make([]*model.ConceptLine, 0), nil
	}
	// remove jsonp wrap
	data = data[len(fmt.Sprintf("quotebridge_v4_line_bk_%s_01_last(", plateId)) : len(data)-1]
	var cLine dto.ConceptLine
	if err = json.Unmarshal(data, &cLine); err != nil {
		return nil, errors.Wrap(err, "FetchConceptLine, plateId="+plateId)
	}

	lines, latestIncluded, err := cLine.ConverTo(plateId)
	if err != nil {
		return nil, errors.Wrap(err, "FetchConceptLine, plateId="+plateId)
	}
	if !latestIncluded {
		log.Println("FetchConceptLine: latestIncluded=fales, plateId=" + plateId)
		line, err := fetchConceptLineToday(ctx, plateId)
		if err != nil {
			return nil, errors.Wrap(err, "FetchConceptLine, plateId="+plateId)
		}
		lines = append(lines, line)
	}
	return lines, err
}

func fetchConceptLineToday(ctx context.Context, plateId string) (*model.ConceptLine, error) {
	headers := map[string]string{
		"Referer": "http://q.10jqka.com.cn/",
	}
	resp, err := util.HttpGet(ctx, fmt.Sprintf(ConceptLineTodayUrl, plateId), headers, nil)
	if err != nil {
		return nil, errors.Wrap(err, "fetchConceptLineToday, plateId="+plateId)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "fetchConceptLineToday, plateId="+plateId)
	}
	if len(data) == 0 {
		log.Printf("fetchConceptLineToday: plate_id %s return 0 byte", plateId)
		return nil, errors.Wrap(err, "fetchConceptLineToday: return 0 byte, plateId="+plateId)
	}
	// remove jsonp wrap
	data = data[len(fmt.Sprintf("quotebridge_v4_line_bk_%s_01_today(", plateId)) : len(data)-1]
	dto := make(dto.ConceptLineToday)
	if err = json.Unmarshal(data, &dto); err != nil {
		return nil, errors.Wrap(err, "fetchConceptLineToday, plateId="+plateId)
	}
	return dto.ConverTo(plateId)
}
