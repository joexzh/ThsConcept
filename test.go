package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joexzh/ThsConcept/fetch"
)

func test() {
	list, err := fetch.SohuZDT(context.Background())
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%v\n", list)
}
