package fetch

import (
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joexzh/ThsConcept/model"
)

const errPrefix = "sohu api"

// SohuZDT 获取搜狐的涨跌停历史数据
// url: https://q.stock.sohu.com/cn/zdt.shtml
func SohuZDT() ([]model.ZDTHistory, error) {
	resp, err := http.Get("https://q.stock.sohu.com/cn/zdt.shtml")
	if err != nil {
		return nil, errors.Wrap(err, errPrefix)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var errList []error

	// 跳过表头, 占2个tr
	trs := doc.Find(".data-main .data-table table tbody tr").Slice(2, goquery.ToEnd)
	length := trs.Length()
	list := make([]model.ZDTHistory, length, length)
	trs.Each(func(tri int, selection *goquery.Selection) {
		selection.Find("td").Each(func(tdi int, td *goquery.Selection) {
			err = parseTd(tdi, td, &list[tri], now)
			if err != nil {
				errList = append(errList, errors.Wrap(err, errPrefix))
			}
		})
	})

	if len(errList) > 0 {
		return nil, errList[0]
	}

	return list, nil

}

// parseTd parse <td>xxx</td> content to ZDTHistory
func parseTd(i int, td *goquery.Selection, zdt *model.ZDTHistory, now time.Time) error {
	text := td.Text()
	text = strings.Trim(text, " \n\r")

	switch i {
	case 0: // 日期
		monthDay := strings.Split(text, "/")
		if len(monthDay) < 2 {
			return errors.Errorf("wrong date format: %s", text)
		}
		month, err := strconv.ParseInt(monthDay[0], 10, 8)
		if err != nil {
			return err
		}
		day, err := strconv.ParseInt(monthDay[1], 10, 8)
		if err != nil {
			return err
		}

		year := now.Year()
		if int64(now.Month()) < month {
			year -= 1
		}
		loc, err := time.LoadLocation("UTC")
		if err != nil {
			return err
		}
		zdt.Date = time.Date(year, time.Month(month), int(day), 0, 0, 0, 0, loc)

	case 1, 2, 3, 5, 6, 7, 8, 9, 10:
		num, err := strconv.ParseInt(text, 10, 16)
		if err != nil {
			return err
		}
		count := uint16(num)
		switch i {
		case 1: // 涨停只数
			zdt.LongLimitCount = count
		case 2: // 跌停只数
			zdt.ShortLimitCount = count
		case 3: // 停牌只数
			zdt.StopTradeCount = count
		case 5: // 沪市上涨只数
			zdt.SHLongCount = count
		case 6: // 沪市平盘只数
			zdt.SHEvenCount = count
		case 7: // 沪市下跌只数
			zdt.SHShortCount = count
		case 8: // 深市上涨只数
			zdt.SZLongCount = count
		case 9: // 深市平盘只数
			zdt.SZEvenCount = count
		case 10: // 深市下跌只数
			zdt.SZShortCount = count
		}

	case 4: // 两市交易(亿)
		amount, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return err
		}
		zdt.Amount = amount
	}

	return nil
}
