package main

import (
	"flag"
	"log"
)

func main() {
	const _Usage = `This is mode this program works. 
1. Mode "server" will start a server for our concept query in Mongodb.
2. Mode "retrieve" will retrieve concept data from ths api, and store the data to our Mongodb.
3. Mode "test" run the test() function.`
	var mode string
	flag.StringVar(&mode, "mode", "server", _Usage)
	flag.Parse()

	switch mode {
	case "server":
		startServer()
	case "retrieve":
		retrieveData()
	case "test":
		test()
	default:
		log.Fatal("wrong mode")
	}
}
