package main

import (
	"context"
	"github.com/joexzh/ThsConcept/model"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joexzh/ThsConcept/dto"
	"github.com/joexzh/ThsConcept/joexzherror"
	"github.com/joexzh/ThsConcept/repos"
)

func ginQuerySc(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 32)
	stock := c.Query("stock")
	concept := c.Query("concept")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repo, err := repos.NewStockMarketRepo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	concepts, err := repo.QueryStockConcept(ctx, stock, concept, int(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	var stocks []*model.ConceptStock
	for _, c := range concepts {
		stocks = append(stocks, c.Stocks...)
	}
	sort.Sort(model.ConceptStockByUpdateAtDesc(stocks))

	c.JSON(http.StatusOK, wrapResult(0, "", &dto.ConceptsDto{
		Concepts: concepts,
		Stocks:   stocks,
	}))
}

func wrapResult(code int, msg string, result interface{}) gin.H {
	if code == 0 && msg == "" {
		msg = "success"
	}
	return gin.H{"code": code, "msg": msg, "result": result}
}

func errCode(err error) int {
	if _, ok := err.(joexzherror.BizError); ok {
		return 1
	}
	return 2
}
