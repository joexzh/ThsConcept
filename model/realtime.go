package model

import (
	"database/sql/driver"

	"github.com/joexzh/ThsConcept/db"
	"github.com/joexzh/dbh"
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

func (r *RealtimeMessage) Args() []any {
	return []any{
		&r.UserId,
		&r.Id,
		&r.Seq,
		&r.Title,
		&r.Digest,
		&r.Url,
		&r.AppUrl,
		&r.ShareUrl,
		&r.Stock,
		&r.Field,
		&r.Color,
		&r.Tag,
		&r.Tags,
		&r.Ctime,
		&r.Rtime,
		&r.Source,
		&r.Short,
		&r.Nature,
		&r.Import,
		&r.TagInfo,
	}
}
func (r *RealtimeMessage) Columns() []string {
	return []string{
		"user_id",
		"id",
		"seq",
		"title",
		"digest",
		"url",
		"app_url",
		"share_url",
		"stock",
		"field",
		"color",
		"tag",
		"tags",
		"ctime",
		"rtime",
		"source",
		"short",
		"nature",
		"import",
		"tag_info",
	}
}
func (r *RealtimeMessage) TableName() string {
	return "realtime_archive"
}
func (r *RealtimeMessage) Config() *dbh.Config {
	return dbh.DefaultConfig
}

type RealtimeMessageStock struct {
	Name        string `json:"name"`
	StockCode   string `json:"stockCode"`
	StockMarket string `json:"stockMarket"`
}

type RealtimeMessageStocks []RealtimeMessageStock

func (s *RealtimeMessageStocks) Scan(src interface{}) error {
	return db.JsonScan(s, src)
}
func (s RealtimeMessageStocks) Value() (driver.Value, error) {
	return db.JsonValue(s)
}

type RealtimeTagInfo struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Score string `json:"score"`
	Type  string `json:"type"`
}

type RealtimeTagInfos []RealtimeTagInfo

func (i *RealtimeTagInfos) Scan(src interface{}) error {
	return db.JsonScan(i, src)
}
func (i RealtimeTagInfos) Value() (driver.Value, error) {
	return db.JsonValue(i)
}

type RealtimeTag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Bury string `json:"bury"`
}

type RealtimeTags []RealtimeTag

func (t *RealtimeTags) Scan(src interface{}) error {
	return db.JsonScan(t, src)
}
func (t RealtimeTags) Value() (driver.Value, error) {
	return db.JsonValue(t)
}
