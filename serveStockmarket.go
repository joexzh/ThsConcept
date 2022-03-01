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
	repo, err := repos.GetStockMarketRepo()
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	zdts, err := repo.QueryLongShort(context.Background(), time.Now().AddDate(-2, 0, 0), repos.DateAsc, 100)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, zdts)
}
