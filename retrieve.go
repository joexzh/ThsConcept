package main

import (
	"context"
	"fmt"
	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/fetch"
	"github.com/joexzh/writeline"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/joexzh/ThsConcept/model"
	"github.com/joexzh/ThsConcept/repos"
)

var wl *writeline.LineWriter

func init() {
	_wl, err := writeline.NewWithStdout(3)
	if err != nil {
		log.Fatal(err)
	}
	wl = _wl
}

func retrieveData() {
	defer wl.Close()
	var wg sync.WaitGroup
	log.Println("Starting to retrieve data")

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		retrieveConcept() // 从同花顺获取概念
		wg.Done()
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		retrieveSohuZdt() // 从搜狐网获取涨跌停数据
		wg.Done()
	}(&wg)

	wg.Wait()
}

func retrieveConcept() {
	wl.WriteLine(0, "开始获取concept...")
	throttleChan := make(chan struct{}, config.Throttle) // throttle goroutine, prevent system or network crash
	rand.Seed(time.Now().UnixNano())                     // 用于每个goroutine随机睡眠

	// get concept ids from page
	cids, err := fetch.ConceptCodesFromPage()
	if err != nil {
		log.Fatal(err)
	}
	wl.WriteLine(0, fmt.Sprint("开始获取concept...concept len: ", len(cids)))

	// get concepts from api
	conceptChan := make(chan *model.Concept)
	conceptSlice := make([]model.Concept, 0, len(cids))
	for _, cid := range cids {
		go func(cid string) {
			throttleChan <- struct{}{}
			errStr := "goroutine 2 error: conceptId: %v, err: %v"

			defer func() {
				if r := recover(); r != nil {
					switch err := r.(type) {
					case error:
						wl.WriteNewLine(fmt.Sprintf(errStr, cid, err.Error()))
					default:
						wl.WriteNewLine(fmt.Sprintf(errStr, cid, err))
					}

					conceptChan <- nil
				}

				time.Sleep(time.Duration(rand.Intn(config.SleepRandUpTo)) * time.Millisecond) // 随机睡眠
				<-throttleChan
			}()

			ret, err := fetch.ConceptFromConceptListApi(cid)
			if err != nil {
				panic(err)
			}
			define, err := fetch.ConceptDefineFromPage(cid)
			if err != nil {
				panic(err)
			}
			wl.WriteLine(1, fmt.Sprintf("从api获取了concept: %v...%v/%v", ret.Result.Name, len(conceptSlice)+1, len(cids)))
			ret.Result.Define = define
			concept, err := ret.ConvertToConcept()
			if err != nil {
				panic(err)
			}
			conceptChan <- concept
		}(cid)
	}

	for i := 0; i < len(cids); i++ {
		concept := <-conceptChan
		if concept != nil {
			conceptSlice = append(conceptSlice, *concept)
		}
	}

	if len(conceptSlice) < 1 {
		wl.WriteLine(0, "获得的概念列表为空, 请检查网络或控制goroutine的并发数量")
		return
	}
	wl.WriteLine(0, "开始获取concept...done")

	ctx := context.Background()
	repo, err := repos.NewConceptRepo()
	if err != nil {
		log.Fatal(err)
	}

	deleted, updated, err := repo.UpdateConceptColl(ctx, conceptSlice...)
	if err != nil {
		log.Default()
		log.Fatal(err)
	}
	wl.WriteNewLine(fmt.Sprintf("UpdateConceptColl: deleted %v, updated %v", deleted, updated))

	deleted, updated, err = repo.UpdateStockConcept(ctx, conceptSlice...)
	if err != nil {
		log.Fatal(err)
	}
	wl.WriteNewLine(fmt.Sprintf("UpdateStockConcept: deleted %v, updated %v", deleted, updated))
}

func retrieveSohuZdt() {
	list, err := fetch.SohuZDT()
	if err != nil {
		log.Fatal(err)
	}
	if len(list) < 1 {
		return
	}

	repo, err := repos.NewStockMarketRepo()
	if err != nil {
		log.Fatal(err)
	}
	date := time.Now().AddDate(0, -3, 0)
	dbList, err := repo.ZdtListDesc(context.Background(), date, 1)
	if err != nil {
		log.Fatal(err)
	}

	var newList []model.ZDTHistory
	if len(dbList) > 0 {
		lastDate := dbList[0].Date
		for _, zdt := range list {
			if zdt.Date.After(lastDate) {
				newList = append(newList, zdt)
			}
		}
	} else {
		newList = list
	}

	rows, err := repo.InsertZdtList(context.Background(), newList)
	if err != nil {
		log.Fatal(err)
	}
	wl.WriteLine(2, fmt.Sprintf("zdt: inserted %v rows to stock_market/long_short", rows))
}
