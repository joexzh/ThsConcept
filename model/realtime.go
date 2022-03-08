package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RealtimeSavedMessage struct {
	Id      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId  string             `json:"userId" bson:"userId"`
	Message RealtimeMessage    `json:"message" bson:"message"`
}

type RealtimeData struct {
	List   []RealtimeMessage `json:"list"`
	Filter []RealtimeTag     `json:"filter"`
	Total  string            `json:"total"`
}

type RealtimeMessage struct {
	Id       string                 `json:"id" bson:"id"`
	Seq      string                 `json:"seq" bson:"seq"`
	Title    string                 `json:"title" bson:"title"`
	Digest   string                 `json:"digest" bson:"digest"`
	Url      string                 `json:"url" bson:"url"`
	AppUrl   string                 `json:"appUrl" bson:"appUrl"`
	ShareUrl string                 `json:"shareUrl" bson:"shareUrl"`
	Stock    []RealtimeMessageStock `json:"stock" bson:"stock"`
	Field    []RealtimeMessageStock `json:"field" bson:"field"`
	Color    string                 `json:"color" bson:"color"`
	Tag      string                 `json:"tag" bson:"tag"`
	Tags     []RealtimeTag          `json:"tags" bson:"tags"`
	Ctime    string                 `json:"ctime" bson:"ctime"`
	Rtime    string                 `json:"rtime" bson:"rtime"`
	Source   string                 `json:"source" bson:"source"`
	Short    string                 `json:"short" bson:"short"`
	Nature   string                 `json:"nature" bson:"nature"`
	Import   string                 `json:"import" bson:"import"`
	TagInfo  []RealtimeTagInfo      `json:"tagInfo" bson:"tagInfo"`
}

type RealtimeMessageStock struct {
	Name        string `json:"name" bson:"name"`
	StockCode   string `json:"stockCode" bson:"stockCode"`
	StockMarket string `json:"stockMarket" bson:"stockMarket"`
}

type RealtimeTagInfo struct {
	Id    string `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Score string `json:"score" bson:"score"`
	Type  string `json:"type" bson:"type"`
}

type RealtimeTag struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Bury string `json:"bury" bson:"bury,omitempty"`
}
