package main

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joexzh/ThsConcept/config"
	"html/template"
	"log"
)

//go:embed tmpl
var fs embed.FS

func startServer() {
	r := gin.Default()
	r.Use(cors.Default())

	tmpl := template.Must(template.ParseFS(fs, "tmpl/*.tmpl"))
	r.SetHTMLTemplate(tmpl)

	r.GET("/query/:name", ginQuery)
	r.GET("/queryrex/:name", ginQueryRex)
	r.GET("/concept/:conceptId", ginConceptId)
	r.GET("/api/stockconcept", ginQuerySc)
	r.GET("/page/sc", ginPageSc)

	r.GET("/api/realtime", ginRealtimeApi)
	r.GET("/api/realtime/save/:userId", ginRealtimeGetSavedMsgList)
	r.POST("/api/realtime/save/:userId", ginRealtimeSaveMsg)
	r.DELETE("/api/realtime/save/:userId", ginRealtimeDelMsg)

	r.GET("/api/stock/zdt", ginLongShort)

	port := config.GetEnv().ServerPort
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(fmt.Sprintf(":%v", port)))
}
