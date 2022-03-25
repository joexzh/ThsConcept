package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/dto"
	"github.com/joexzh/ThsConcept/util"
)

// ConceptCodesFromPage 从 html 页面 http://q.10jqka.com.cn/gn/ 获取概念列表
func ConceptCodesFromPage(ctx context.Context) ([]string, error) {
	res, err := util.HttpGet(ctx, config.ConceptAllUrl, nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	page, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(config.RexConceptCode)
	matches := re.FindAllStringSubmatch(string(page), -1)
	cids := make([]string, 0, len(matches))
	for _, match := range matches {
		cids = append(cids, match[1])
	}
	cids = util.RemoveDuplicate(cids)
	return cids, nil
}

// ConceptFromConceptListApi 根据 conceptId 从 http://basic.10jqka.com.cn/ajax/stock/conceptlist.php?cid=%v 获取 concept list
func ConceptFromConceptListApi(ctx context.Context, conceptId string) (*dto.ConceptListApiReturn, error) {
	url := fmt.Sprintf(config.ConceptApiUrl, conceptId)
	res, err := util.HttpGet(ctx, url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ret dto.ConceptListApiReturn

	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}
	ret.ConceptId = conceptId
	return &ret, nil
}

// ConceptDefineFromPage 从 html 页面 http://q.10jqka.com.cn/gn/detail/code/%v/ 获取概念定义
func ConceptDefineFromPage(ctx context.Context, conceptId string) (string, error) {
	url := fmt.Sprintf(config.ConceptDetailPageUrl, conceptId)
	res, err := util.HttpGet(ctx, url, nil, nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	data, err = util.GbkToUtf8(data)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(config.RexConceptDefine)
	matches := re.FindAllStringSubmatch(string(data), -1)
	for _, match := range matches {
		return match[1], nil
	}
	return "", nil
}
