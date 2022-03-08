package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joexzh/ThsConcept/repos"
	"log"
	"net/http"
	"time"
)

func ginLongShort(c *gin.Context) {
	repo, err := repos.NewStockMarketRepo()
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	zdts, err := repo.ZdtListDesc(context.Background(), time.Now().AddDate(-2, 0, 0), 0)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, zdts)
}
