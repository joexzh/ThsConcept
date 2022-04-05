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
	err = repo.ConceptStockFtSync(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("concept: concept_stock_ft done")
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
