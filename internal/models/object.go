package models

type DataObject struct {
	ObjectId string      `json:"object_id"`
	Id       string      `json:"id,omitempty"`
	Data     interface{} `json:"data"`
	Token    string      `json:"token,omitempty"`
}

type MessageObject struct {
	Text string `json:"text"`
}

type RouteObject struct {
	Routes interface{} `json:"routes"`
}
type TextObject struct {
	Text string `json:"text"`
}
