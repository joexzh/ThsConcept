package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/fetch"

	"github.com/joexzh/ThsConcept/model"
	"github.com/joexzh/ThsConcept/repos"
)

func retrieveData() {
	var wg sync.WaitGroup
	log.Println("Starting to retrieve data")

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		retrieveConcept()      // 从同花顺获取概念
		syncConceptStockFt()   // 同步 fulltext 表
		retrieveConceptLines() // 从同花顺获取概念日k
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
	throttleChan := make(chan struct{}, config.Throttle) // throttle goroutine, prevent system or network crash
	rand.Seed(time.Now().UnixNano())                     // 用于每个goroutine随机睡眠

	// get concept ids from page
	ctx := context.Background()
	cids, err := fetch.ConceptCodesFromPage(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("concept: start retrieve, total:", len(cids))

	// get concepts from api
	conceptChan := make(chan *model.Concept)
	conceptSlice := make([]*model.Concept, 0, len(cids))
	for _, cid := range cids {
		go func(cid string) {
			throttleChan <- struct{}{}
			errStr := "concept: goroutine 2 error: conceptId: %v, err: %v"

			defer func() {
				if r := recover(); r != nil {
					switch err := r.(type) {
					case error:
						log.Println(fmt.Sprintf(errStr, cid, err.Error()))
					default:
						log.Println(fmt.Sprintf(errStr, cid, err))
					}

					conceptChan <- nil
				}

				time.Sleep(time.Duration(rand.Intn(config.SleepRandUpTo)) * time.Millisecond) // 随机睡眠
				<-throttleChan
			}()

			ret, err := fetch.ConceptFromConceptListApi(ctx, cid)
			if err != nil {
				panic(err)
			}
			define, err := fetch.ConceptDefineFromPage(ctx, cid)
			if err != nil {
				panic(err)
			}
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
			conceptSlice = append(conceptSlice, concept)
		}
	}

	if len(conceptSlice) < 1 {
		log.Println("concept: 获得的概念列表为空, 请检查网络或控制goroutine的并发数量")
		return
	}
	log.Println("concept: retrieve done")

	repo, err := repos.InitStockMarketRepo()
	if err != nil {
		log.Fatal(err)
	}
	updateResult, err := repo.UpdateConcept(ctx, conceptSlice...)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("concept: concept_concept")
	log.Printf("concept: inserted: %d\n", updateResult.ConceptConceptInserted)
	log.Printf("concept: updated: %d\n", updateResult.ConceptConceptInserted)
	log.Printf("concept: deleted: %d\n", updateResult.ConceptConceptInserted)
	log.Println("concept: concept_stock")
	log.Printf("concept: inserted: %d\n", updateResult.ConceptStockInserted)
	log.Printf("concept: updated: %d\n", updateResult.ConceptStockUpdated)
	log.Printf("concept: deleted: %d\n", updateResult.ConceptStockDeleted)
}

func syncConceptStockFt() {
	repo, err := repos.InitStockMarketRepo()
	if err != nil {
		log.Fatal(err)
	}
	err = repo.ConceptStockFtSync(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("concept: concept_stock_ft done")
}

func retrieveConceptLines() {
	now := time.Now().In(config.ChinaLoc())
	h, _, _ := now.Clock()
	if h < 17 {
		now = now.AddDate(0, 0, -1)
	}
	weekDay := now.Weekday()
	for weekDay == time.Sunday || weekDay == time.Saturday {
		now = now.AddDate(0, 0, -1)
		weekDay = now.Weekday()
	}

	ctx := context.Background()
	repo, err := repos.InitStockMarketRepo()
	if err != nil {
		log.Fatal(err)
	}
	pIds, err := repo.QueryAllPlateIds(ctx)
	if err != nil {
		log.Fatal(err)
	}
	t, _ := time.ParseInLocation(config.TimeLayoutDate, now.Format(config.TimeLayoutDate), config.ChinaLoc())
	latestLinesDb, err := repo.QueryConceptLineByDate(ctx, t)
	if err != nil {
		log.Fatal(err)
	}
	exclude := make(map[string]struct{})
	exclude["885582"] = struct{}{} // 跳过`新股发行`概念
	for i := range latestLinesDb {
		exclude[latestLinesDb[i].PlateId] = struct{}{}
	}

	lineMap := make(map[string][]*model.ConceptLine)
	log.Println("concept_line: start")
	log.Println("concept_line: exclude len", len(exclude))
	for i := range pIds {
		if _, ok := exclude[pIds[i]]; ok {
			continue
		}
		lines, err := fetch.FetchConceptLine(ctx, pIds[i])
		if err != nil {
			log.Println("concept_line:", err)
			continue // 忽略被server限制的请求, try next time
		}
		lineMap[pIds[i]] = lines
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}
	lines, err := repo.FilterConceptLineMap(ctx, lineMap)
	if err != nil {
		log.Fatal(err)
	}
	r, err := repo.InsertConceptLines(ctx, lines)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("concept_line: inserted: %d\n", r)
}

func retrieveSohuZdt() {
	list, err := fetch.SohuZDT(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(list) < 1 {
		return
	}

	repo, err := repos.InitStockMarketRepo()
	if err != nil {
		log.Fatal(err)
	}
	date := time.Now().AddDate(0, -3, 0)
	dbList, err := repo.ZdtListDesc(context.Background(), date, 1)
	if err != nil {
		log.Fatal(err)
	}

	var newList []*model.ZDTHistory
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
	log.Printf("zdt: inserted %v rows to stock_market/long_short\n", rows)
}
