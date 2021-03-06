package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joexzh/ThsConcept/config"
	"github.com/joexzh/ThsConcept/joexzherror"
	"github.com/joexzh/ThsConcept/repos"
)

func ginQuerySc(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 32)
	stock := c.Query("stock")
	concept := c.Query("concept")

	ctx := context.Background()

	repo, err := repos.InitStockMarketRepo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	scs, err := repo.QueryConceptFtStockByKw(ctx, stock, concept, int(limit))
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

	repo, err := repos.InitStockMarketRepo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	concepts, err := repo.QueryConcepts(ctx, concept, int(limit), false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, wrapResult(0, "", concepts))
}

func ginQueryStockByConceptId(c *gin.Context) {
	conceptId := c.Param("conceptId")

	ctx := context.Background()

	repo, err := repos.InitStockMarketRepo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	stocks, err := repo.QueryConceptStockByConceptId(ctx, conceptId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, wrapResult(0, "", stocks))
}

func ginConceptLineCmp(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")

	var startDate time.Time
	var endDate time.Time
	var err error
	if start != "" {
		startDate, err = time.ParseInLocation(config.TimeLayoutDate, start, config.ChinaLoc())
		if err != nil {
			startDate = time.Time{}
		}
	}
	if end != "" {
		endDate, err = time.ParseInLocation(config.TimeLayoutDate, end, config.ChinaLoc())
		if err != nil {
			endDate = time.Time{}
		}
	}

	ctx := context.Background()
	repo, err := repos.InitStockMarketRepo()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	view, err := repo.ViewConceptLineByDateRange(ctx, startDate, endDate)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, view)
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
