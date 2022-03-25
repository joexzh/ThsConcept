package util

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/joexzh/ThsConcept/config"
)

// HttpGet Default header User-Agent will be auto set.
func HttpGet(ctx context.Context, url string, headers map[string]string, query map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", config.UserAgent)
	for k, v := range headers {
		req.Header.Set(k, headers[v])
	}
	req.URL.RawQuery = appendRawQuery(req.URL, query)
	client := http.Client{}
	return client.Do(req)
}

func appendRawQuery(url *url.URL, query map[string]string) string {
	q := url.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	return q.Encode()
}

// HttpGetRealTime
//
// all params are optional
func HttpGetRealTime(ctx context.Context, page int, pagesize int, tag string, ctime int) (*http.Response, error) {
	query := make(map[string]string, 4)
	if page > 0 {
		query["page"] = strconv.FormatInt(int64(page), 10)
		query["pagesize"] = strconv.FormatInt(int64(pagesize), 10)
	}
	query["tag"] = tag
	if ctime > 0 {
		query["ctime"] = strconv.FormatInt(int64(ctime), 10)
	}

	return HttpGet(ctx, config.RealTimeUrl, nil, query)
}
