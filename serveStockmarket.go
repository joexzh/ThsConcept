package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joexzh/ThsConcept/repos"
)

func ginLongShort(c *gin.Context) {
	repo, err := repos.InitStockMarketRepo()
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	zdts, err := repo.ZdtListDesc(ctx, time.Now().AddDate(-2, 0, 0), 0)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, zdts)
}
