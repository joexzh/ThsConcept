package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joexzh/ThsConcept/dto"
	"github.com/joexzh/ThsConcept/joexzherror"
	"github.com/joexzh/ThsConcept/repos"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ginQuery(c *gin.Context) {
	conceptName := c.Param("name")

	ctx := context.Background()
	repo, err := repos.NewConceptRepo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}

	concept, err := repo.QueryByConceptNameRex(ctx, conceptName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, wrapResult(0, "", concept))
}

func ginQueryRex(c *gin.Context) {
	conceptName := c.Param("name")

	ctx := context.Background()
	repo, err := repos.NewConceptRepo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}

	concept, err := repo.QueryByConceptNameRex(ctx, conceptName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, wrapResult(0, "", concept))
}

func ginConceptId(c *gin.Context) {
	conceptId := c.Param("conceptId")

	ctx := context.Background()
	repo, err := repos.NewConceptRepo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}

	concept, err := repo.QueryByConceptId(ctx, conceptId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, wrapResult(0, "", concept))
}

func ginQuerySc(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 32)
	stockName := c.Query("stock")
	conceptNameRegex := c.Query("concept")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dtos, err := scDtos(ctx, conceptNameRegex, stockName, int(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapResult(errCode(err), err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, wrapResult(0, "", dtos))
}

func ginPageSc(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 32)
	stockName := c.Query("stockname")
	conceptRegex := c.Query("concept")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dtos, err := scDtos(ctx, conceptRegex, stockName, int(limit))
	scPageDto := dto.ScPageDto{
		Concept:   conceptRegex,
		StockName: stockName,
		Scs:       dtos,
	}
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", scPageDto)
}

func scDtos(ctx context.Context, concept string, stockName string, limit int) ([]dto.StockConceptDto, error) {
	repo, err := repos.NewConceptRepo()
	if err != nil {
		return nil, errors.Wrap(err, "failed to NewConceptRepo")
	}

	scs, err := repo.QueryScDesc(ctx, stockName, concept, limit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to QueryScDesc")
	}

	dtos := dto.ScToScDto(scs...)
	return dtos, nil
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
