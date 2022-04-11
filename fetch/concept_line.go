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

const ConceptIndexUrl = "http://d.10jqka.com.cn/v4/line/bk_%s/01/last.js"

func FetchConceptLine(ctx context.Context, plateId string) ([]*model.ConceptLine, error) {
	resp, err := util.HttpGet(ctx, fmt.Sprintf(ConceptIndexUrl, plateId), nil, nil)
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
	return cLine.ConverTo(plateId)
}
