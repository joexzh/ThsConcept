package main

import (
	"context"
	"encoding/json"
	"github.com/joexzh/ThsConcept/repos"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joexzh/ThsConcept/dto"
	"github.com/joexzh/ThsConcept/model"
	"github.com/joexzh/ThsConcept/util"
)

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first proxy retain prior
	// X-Forwarded-For information as a comma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

func ginRealtimeApi(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 10, 32)
	pagesize, _ := strconv.ParseInt(c.Query("pagesize"), 10, 32)
	tag := c.Query("tag")
	ctime, _ := strconv.ParseInt(c.Query("ctime"), 10, 64)

	ctx := context.Background()
	resp, err := util.HttpGetRealTime(ctx, int(page), int(pagesize), tag, int(ctime))
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	decoder := json.NewDecoder(resp.Body)
	var apiResp = dto.RealtimeResponse{}
	if err = decoder.Decode(&apiResp); err != nil {
		log.Println(err.Error())
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	dto := dto.RealtimeDto{
		List:   apiResp.Data.List,
		Filter: apiResp.Data.Filter,
		Total:  apiResp.Data.Total,
	}

	c.JSON(http.StatusOK, &dto)
}

func ginRealtimeApiRaw(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 10, 32)
	pagesize, _ := strconv.ParseInt(c.Query("pagesize"), 10, 32)
	tag := c.Query("tag")
	ctime, _ := strconv.ParseInt(c.Query("ctime"), 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := util.HttpGetRealTime(ctx, int(page), int(pagesize), tag, int(ctime))
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	defer resp.Body.Close()

	delHopHeaders(resp.Header)
	copyHeader(c.Writer.Header(), resp.Header)
	c.Writer.WriteHeader(http.StatusOK)
	io.Copy(c.Writer, resp.Body)
}

// save message list
func ginSaveRealtimeArchive(c *gin.Context) {
	var msg model.RealtimeMessage
	if err := c.BindJSON(&msg); err != nil {
		log.Println(err)
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// todo fake user
	msg.UserId = 1

	ctx := context.Background()

	repo, err := repos.NewStockMarketRepo()
	if err != nil {
		log.Println(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	_, err = repo.SaveRealtimeArchive(ctx, &msg)
	if err != nil {
		log.Println(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func ginDeleteRealtimeArchive(c *gin.Context) {
	seq := c.Param("seq")
	userId := 1 // todo fake user
	ctx := context.Background()

	repo, err := repos.NewStockMarketRepo()
	if err != nil {
		log.Println(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	_, err = repo.DeleteRealtimeArchive(ctx, userId, seq)
	if err != nil {
		log.Println(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func ginRealtimeArchive(c *gin.Context) {
	ctx := context.Background()

	userId := 1 // todo fake user

	repo, err := repos.NewStockMarketRepo()
	if err != nil {
		log.Println(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	list, err := repo.QueryRealtimeArchive(ctx, userId, 1000)
	if err != nil {
		log.Println(err)
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, list)
}
