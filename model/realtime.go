package model

import (
	"database/sql/driver"
	"github.com/joexzh/ThsConcept/db"
)

type RealtimeData struct {
	List   []RealtimeMessage `json:"list"`
	Filter []RealtimeTag     `json:"filter"`
	Total  string            `json:"total"`
}

type RealtimeMessage struct {
	UserId   int                   `json:"userId" db:"user_id"`
	Id       string                `json:"id" db:"id"`
	Seq      string                `json:"seq" db:"seq"`
	Title    string                `json:"title" db:"title"`
	Digest   string                `json:"digest" db:"digest"`
	Url      string                `json:"url" db:"url"`
	AppUrl   string                `json:"appUrl" db:"app_url"`
	ShareUrl string                `json:"shareUrl" db:"share_url"`
	Stock    RealtimeMessageStocks `json:"stock" db:"stock"`
	Field    RealtimeMessageStocks `json:"field" db:"field"`
	Color    string                `json:"color" db:"color"`
	Tag      string                `json:"tag" db:"tag"`
	Tags     RealtimeTags          `json:"tags" db:"tags"`
	Ctime    string                `json:"ctime" db:"ctime"`
	Rtime    string                `json:"rtime" db:"rtime"`
	Source   string                `json:"source" db:"source"`
	Short    string                `json:"short" db:"short"`
	Nature   string                `json:"nature" db:"nature"`
	Import   string                `json:"import" db:"import"`
	TagInfo  RealtimeTagInfos      `json:"tagInfo" db:"tag_info"`
}

type RealtimeMessageStock struct {
	Name        string `json:"name"`
	StockCode   string `json:"stockCode"`
	StockMarket string `json:"stockMarket"`
}

type RealtimeMessageStocks []RealtimeMessageStock

func (s *RealtimeMessageStocks) Scan(src interface{}) error {
	return db.Scan(s, src)
}
func (s RealtimeMessageStocks) Value() (driver.Value, error) {
	return db.Value(s)
}

type RealtimeTagInfo struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Score string `json:"score"`
	Type  string `json:"type"`
}

type RealtimeTagInfos []RealtimeTagInfo

func (i *RealtimeTagInfos) Scan(src interface{}) error {
	return db.Scan(i, src)
}
func (i RealtimeTagInfos) Value() (driver.Value, error) {
	return db.Value(i)
}

type RealtimeTag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Bury string `json:"bury"`
}

type RealtimeTags []RealtimeTag

func (t *RealtimeTags) Scan(src interface{}) error {
	return db.Scan(t, src)
}
func (t RealtimeTags) Value() (driver.Value, error) {
	return db.Value(t)
}
