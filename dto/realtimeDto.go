package dto

import (
	"github.com/joexzh/ThsConcept/model"
)

type RealtimeDto struct {
	List   []model.RealtimeMessage `json:"list"`
	Filter []model.RealtimeTag     `json:"filter"`
	Total  string                  `json:"total"`
}

type RealtimeResponse struct {
	Code string             `json:"code"`
	Msg  string             `json:"msg"`
	Time string             `json:"time"`
	Data model.RealtimeData `json:"data"`
}
