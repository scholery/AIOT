package model

type Device struct {
	Key      string     `json:"key"`
	Name     string     `json:"name"`
	SourceId string     `json:"sourceId"`
	Geo      [2]float32 `json:"geo"`
	Product  Product    `json:"product"`
	Desc     string     `json:"desc"`
}
