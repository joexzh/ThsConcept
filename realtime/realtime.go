package realtime

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SavedMessage struct {
	Id      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId  string             `json:"userId" bson:"userId"`
	Message Message            `json:"message" bson:"message"`
}

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Time string `json:"time"`
	Data Data   `json:"data"`
}

type Data struct {
	List   []Message `json:"list"`
	Filter []Tag     `json:"filter"`
	Total  string    `json:"total"`
}

type Message struct {
	Id            string         `json:"id" bson:"id"`
	Seq           string         `json:"seq" bson:"seq"`
	Title         string         `json:"title" bson:"title"`
	Digest        string         `json:"digest" bson:"digest"`
	Url           string         `json:"url" bson:"url"`
	AppUrl        string         `json:"appUrl" bson:"appUrl"`
	ShareUrl      string         `json:"shareUrl" bson:"shareUrl"`
	Stock         []Stock        `json:"stock" bson:"stock"`
	Field         []Stock        `json:"field" bson:"field"`
	Color         string         `json:"color" bson:"color"`
	Tag           string         `json:"tag" bson:"tag"`
	Tags          []Tag          `json:"tags" bson:"tags"`
	Ctime         string         `json:"ctime" bson:"ctime"`
	Rtime         string         `json:"rtime" bson:"rtime"`
	Source        string         `json:"source" bson:"source"`
	Short         string         `json:"short" bson:"short"`
	Nature        string         `json:"nature" bson:"nature"`
	Import        string         `json:"import" bson:"import"`
	TagInfo       []TagInfo      `json:"tagInfo" bson:"tagInfo"`
	KeywordCounts []KeywordCount `json:"KeywordCounts" bson:"KeywordCounts"`
}

type Stock struct {
	Name        string `json:"name" bson:"name"`
	StockCode   string `json:"stockCode" bson:"stockCode"`
	StockMarket string `json:"stockMarket" bson:"stockMarket"`
}

type TagInfo struct {
	Id    string `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Score string `json:"score" bson:"score"`
	Type  string `json:"type" bson:"type"`
}

type Tag struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Bury string `json:"bury" bson:"bury,omitempty"`
}
