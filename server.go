package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joexzh/ThsConcept/config"
	"log"
	"path/filepath"
)

func startServer() {
	r := gin.Default()
	r.Use(cors.Default())

	r.LoadHTMLGlob(filepath.Join(ExeDir(), "tmpl/*"))

	r.GET("/query/:name", ginQuery)
	r.GET("/queryrex/:name", ginQueryRex)
	r.GET("/:conceptId", ginConceptId)
	r.GET("/sc", ginQuerySc)
	r.GET("/page/sc", ginPageSc)

	port := config.GetEnv().ServerPort
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(fmt.Sprintf(":%v", port)))
}
