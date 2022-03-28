package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joexzh/ThsConcept/joexzherror"
	"github.com/joexzh/ThsConcept/repos"
	"net/http"
	"strconv"
)

func ginQuerySc(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 32)
	stock := c.Query("stock")
	concept := c.Query("concept")

	ctx := context.Background()

	repo, err := repos.NewStockMarketRepo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	scs, err := repo.QueryConceptStock(ctx, stock, concept, int(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, wrapResult(0, "", scs))
}

func ginQueryConcept(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 32)
	concept := c.Query("concept")

	ctx := context.Background()

	repo, err := repos.NewStockMarketRepo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	concepts, err := repo.QueryConcepts(ctx, concept, int(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, wrapResult(0, "", concepts))
}

func ginQueryStockByConceptId(c *gin.Context) {
	conceptId := c.Param("conceptId")

	ctx := context.Background()

	repo, err := repos.NewStockMarketRepo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	stocks, err := repo.QueryStockByConceptId(ctx, conceptId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, wrapResult(0, "", stocks))
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
