package main

import (
	"context"
	"fmt"
	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/fetch"
	"log"
	"math/rand"
	"time"

	"github.com/joexzh/ThsConcept/model"
	"github.com/joexzh/ThsConcept/repos"
)

func retrieveData() {
	fmt.Println("Starting to retrieve data")

	retrieveConcept() // 从同花顺获取概念
	retrieveSohuZdt() // 从搜狐网获取涨跌停数据
}

func retrieveConcept() {
	throttleChan := make(chan struct{}, config.Throttle) // throttle goroutine, prevent system or network crash
	rand.Seed(time.Now().UnixNano())                     // 用于每个goroutine随机睡眠

	// get concept ids from page
	cids, err := ConceptCodesFromPage()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Match concept codes: ", cids)

	// get concepts from api
	conceptChan := make(chan *model.Concept)
	for _, cid := range cids {
		go func(cid string) {
			throttleChan <- struct{}{}
			errStr := "goroutine 2 error: conceptId: %v, err: %v\n"

			defer func() {
				if r := recover(); r != nil {
					switch err := r.(type) {
					case error:
						fmt.Printf(errStr, cid, err.Error())
					default:
						fmt.Printf(errStr, cid, err)
					}

					conceptChan <- nil
				}

				time.Sleep(time.Duration(rand.Intn(config.SleepRandUpTo)) * time.Millisecond) // 随机睡眠
				<-throttleChan
			}()

			ret, err := ConceptFromApi(cid)
			if err != nil {
				panic(err)
			}
			define, err := ConceptDefineFromPage(cid)
			if err != nil {
				panic(err)
			}
			fmt.Printf("从api获取了concept: %v\n", ret.Result.Name)
			ret.Result.Define = define
			concept, err := ret.ConvertToConcept()
			if err != nil {
				panic(err)
			}
			conceptChan <- concept
		}(cid)
	}
	conceptSlice := make([]model.Concept, 0, len(cids))
	for i := 0; i < len(cids); i++ {
		concept := <-conceptChan
		if concept != nil {
			conceptSlice = append(conceptSlice, *concept)
		}
	}

	if len(conceptSlice) < 1 {
		fmt.Println("获得的概念列表为空, 请检查网络或控制goroutine的并发数量")
		return
	}

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
	fmt.Printf("UpdateConceptColl: deleted %v, updated %v\n", deleted, updated)

	deleted, updated, err = repo.UpdateStockConcept(ctx, conceptSlice...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("UpdateStockConcept: deleted %v, updated %v\n", deleted, updated)
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
	dbList, err := repo.QueryLongShort(context.Background(), date, repos.DateDesc, 1)
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

	rows, err := repo.InsertLongShort(context.Background(), newList)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("mysql: inserted %v rows to stock_market/long_short", rows)
}
