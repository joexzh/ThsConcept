package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	log.Fatal(r.Run(":8088"))
}
