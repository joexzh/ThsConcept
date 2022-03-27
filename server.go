package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joexzh/ThsConcept/config"
	"log"
)

func startServer() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/api/stockconcept", ginQuerySc)

	r.GET("/api/realtime", ginRealtimeApi)
	r.GET("/api/realtime/archive", ginRealtimeArchive)
	r.POST("/api/realtime/archive", ginSaveRealtimeArchive)
	r.DELETE("/api/realtime/archive", ginDeleteRealtimeArchive)

	r.GET("/api/stock/zdt", ginLongShort)

	port := config.GetEnv().ServerPort
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(fmt.Sprintf(":%v", port)))
}
