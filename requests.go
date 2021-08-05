package main

import (
	"github.com/joexzh/ThsConcept/config"
	"net/http"
)

// HttpGet Default header User-Agent will be auto set.
func HttpGet(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", config.UserAgent)
	for k, v := range headers {
		req.Header.Set(k, headers[v])
	}
	client := http.Client{}
	return client.Do(req)
}
