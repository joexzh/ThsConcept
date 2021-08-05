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

	mode := flag.String("mode", "server", _Usage)
	flag.Parse()

	switch *mode {
	case "server":
		startServer()
		break
	case "retrieve":
		retrieveData()
		break
	case "test":
		test()
		break
	default:
		log.Fatal("wrong mode")
	}

	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
}