package main

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"path/filepath"
)

func RemoveDuplicate(slice []string) []string {
	if len(slice) < 1 {
		return slice
	}

	checkMap := make(map[string]struct{})
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if _, ok := checkMap[item]; !ok {
			checkMap[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func ExeDir() string {
	file, _ := os.Executable()
	return filepath.Dir(file)
}
