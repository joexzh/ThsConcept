package main

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/model"
)

func allStockSymbol() ([]string, error) {
	res, err := HttpGet(config.StockSymbolUrl, nil, nil)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	reg := regexp.MustCompile(config.RexStockSymbol)
	matches := reg.FindAllStringSubmatch(string(data), -1)
	if len(matches) < 1 {
		err = fmt.Errorf("表达式`%v`无法匹配`%v`中的内容", config.RexStockSymbol, config.StockSymbolUrl)
		return nil, err
	}
	slice := make([]string, 0, len(matches))
	for _, match := range matches {
		slice = append(slice, match[1])
	}
	slice = RemoveDuplicate(slice)
	slice = filterValidSymbols(slice)
	return slice, nil
}

func filterValidSymbols(slice []string) []string {
	newSlice := make([]string, 0, len(slice))
	for _, symbol := range slice {
		if isValidSymbol(symbol) {
			newSlice = append(newSlice, symbol)
		}
	}
	return newSlice
}

func isValidSymbol(symbol string) bool {
	if strings.HasPrefix(symbol, "30") ||
		strings.HasPrefix(symbol, "60") ||
		strings.HasPrefix(symbol, "68") ||
		strings.HasPrefix(symbol, "00") {
		return true
	}
	return false
}

func getCidsInOnePage(symbol string) ([]string, error) {
	url := fmt.Sprintf(config.ConceptPageUrl, symbol)
	res, err := HttpGet(url, nil, nil)
	if err != nil {
		return nil, err
	}
	page, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	re := regexp.MustCompile(config.RexValidConceptPage)
	if !re.Match(page) {
		return nil, fmt.Errorf("regex %v didn't match the page %v", config.RexValidConceptPage, url)
	}

	re = regexp.MustCompile(config.RexConceptId)
	matches := re.FindAllStringSubmatch(string(page), -1)
	cids := make([]string, 0, len(matches))
	for _, match := range matches {
		cids = append(cids, match[1])
	}
	return cids, nil
}

func ConceptCodesFromPage() ([]string, error) {
	res, err := HttpGet(config.ConceptAllUrl, nil, nil)
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
	cids = RemoveDuplicate(cids)
	return cids, nil
}

func ConceptFromApi(conceptId string) (*model.Return, error) {
	url := fmt.Sprintf(config.ConceptApiUrl, conceptId)
	res, err := HttpGet(url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ret model.Return

	if err := json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}
	ret.ConceptId = conceptId
	return &ret, nil
}

func ConceptDefineFromPage(conceptId string) (string, error) {
	url := fmt.Sprintf(config.ConceptDetailPageUrl, conceptId)
	res, err := HttpGet(url, nil, nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	data, err = GbkToUtf8(data)
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
