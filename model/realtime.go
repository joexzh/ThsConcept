package model

type RealtimeData struct {
	List   []RealtimeMessage `json:"list"`
	Filter []RealtimeTag     `json:"filter"`
	Total  string            `json:"total"`
}

type RealtimeMessage struct {
	Id       string                 `json:"id"`
	Seq      string                 `json:"seq"`
	Title    string                 `json:"title"`
	Digest   string                 `json:"digest"`
	Url      string                 `json:"url"`
	AppUrl   string                 `json:"appUrl"`
	ShareUrl string                 `json:"shareUrl"`
	Stock    []RealtimeMessageStock `json:"stock"`
	Field    []RealtimeMessageStock `json:"field"`
	Color    string                 `json:"color"`
	Tag      string                 `json:"tag"`
	Tags     []RealtimeTag          `json:"tags"`
	Ctime    string                 `json:"ctime"`
	Rtime    string                 `json:"rtime"`
	Source   string                 `json:"source"`
	Short    string                 `json:"short"`
	Nature   string                 `json:"nature"`
	Import   string                 `json:"import"`
	TagInfo  []RealtimeTagInfo      `json:"tagInfo"`
}

type RealtimeMessageStock struct {
	Name        string `json:"name"`
	StockCode   string `json:"stockCode"`
	StockMarket string `json:"stockMarket"`
}

type RealtimeTagInfo struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Score string `json:"score"`
	Type  string `json:"type"`
}

type RealtimeTag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Bury string `json:"bury"`
}
