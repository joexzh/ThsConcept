package main

import (
	"fmt"
	"github.com/joexzh/ThsConcept/fetch"
	"log"
)

func test() {
	list, err := fetch.SohuZDT()
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%v\n", list)
}
