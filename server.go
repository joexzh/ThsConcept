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
	r.GET("/sc", ginQuerySc)
	r.GET("/page/sc", ginPageSc)

	r.GET("/api/realtime", ginRealtimeApi)
	r.GET("/list/:userId", ginRealtimeGetSavedMsgList)
	r.POST("/list/:userId", ginRealtimeSaveMsg)
	r.DELETE("/list/:userId", ginRealtimeDelMsg)

	port := config.GetEnv().ServerPort
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(fmt.Sprintf(":%v", port)))
}
