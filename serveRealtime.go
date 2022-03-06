package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/joexzh/ThsConcept/realtime"
	"github.com/joexzh/ThsConcept/repos"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
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

type ConceptShortsMapOnce struct {
	ConceptShortsMap map[string][]string
	Once             *sync.Once
}

var conceptShortsMapOnce = ConceptShortsMapOnce{
	Once: new(sync.Once),
}

func getConceptShortsMap() map[string][]string {
	conceptShortsMapOnce.Once.Do(func() {
		ctx := context.Background()
		repo, err := repos.NewRealtimeRepo()
		if err != nil {
			log.Println(err)
			conceptShortsMapOnce.Once = new(sync.Once)
			return
		}

		conceptNames, err := repo.GetAllConceptNames(ctx)
		if err != nil {
			log.Println(err)
			conceptShortsMapOnce.Once = new(sync.Once)
			return
		}
		conceptShortsMapOnce.ConceptShortsMap = realtime.MergeConceptShortsMap(conceptNames)
	})

	return conceptShortsMapOnce.ConceptShortsMap
}

func keywordsCounts(messages []realtime.Message) []realtime.KeywordCount {
	start := time.Now()
	defer func() {
		log.Println("keywordsCounts time: ", time.Since(start))
	}()

	conceptShortsMap := getConceptShortsMap()
	totalKeywordCounts := make([]realtime.KeywordCount, 0)
	if conceptShortsMap == nil {
		return totalKeywordCounts
	}

	for _, msg := range messages {
		kwc := realtime.KeywordCounts(msg.Digest, conceptShortsMap)
		msg.KeywordCounts = kwc
		totalKeywordCounts = append(totalKeywordCounts, kwc...)
	}
	realtime.SortKeywordCounts(totalKeywordCounts)
	return totalKeywordCounts
}

func ginRealtimeApi(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 10, 32)
	pagesize, _ := strconv.ParseInt(c.Query("pagesize"), 10, 32)
	tag := c.Query("tag")
	ctime, _ := strconv.ParseInt(c.Query("ctime"), 10, 64)

	resp, err := HttpGetRealTime(int(page), int(pagesize), tag, int(ctime))
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var apiResp = realtime.Response{}
	if err = decoder.Decode(&apiResp); err != nil {
		log.Println(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// kwc := keywordsCounts(apiResp.Data.List)
	dto := realtime.Dto{
		List:   apiResp.Data.List,
		Filter: apiResp.Data.Filter,
		Total:  apiResp.Data.Total,
		// KeywordCounts: kwc,
	}

	c.JSON(http.StatusOK, &dto)
}

func ginRealtimeApiRaw(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Query("page"), 10, 32)
	pagesize, _ := strconv.ParseInt(c.Query("pagesize"), 10, 32)
	tag := c.Query("tag")
	ctime, _ := strconv.ParseInt(c.Query("ctime"), 10, 64)

	resp, err := HttpGetRealTime(int(page), int(pagesize), tag, int(ctime))
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
func ginRealtimeSaveMsg(c *gin.Context) {
	userId := c.Param("userId")
	var msg realtime.Message
	if err := c.BindJSON(&msg); err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx := context.Background()
	repo, err := repos.NewRealtimeRepo()
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err = repo.SaveMessage(ctx, userId, &msg); err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusOK)
}

func ginRealtimeDelMsg(c *gin.Context) {
	userId := c.Param("userId")
	objId := c.Query("objId")

	ctx := context.Background()
	repo, err := repos.NewRealtimeRepo()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err = repo.DelSaveMessage(ctx, userId, objId); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusOK)
}

// retrieve saved message list
func ginRealtimeGetSavedMsgList(c *gin.Context) {
	userId := c.Param("userId")

	ctx := context.Background()
	repo, err := repos.NewRealtimeRepo()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	list, err := repo.QuerySaveMessageDesc(ctx, userId)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, list)
}
