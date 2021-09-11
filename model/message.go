package model

type PropertyMessage struct {
	SN         string                  `json:"sn"`
	DeviceId   string                  `json:"deviceId"`
	MessageId  string                  `json:"messageId"`
	Timestamp  int64                   `json:"timestamp"`
	Properties map[string]PropertyItem `json:"properties"`
}
type PropertyItem struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Unit  string      `json:"unit"`
	Name  string      `json:"name"`
}
type EventMessage struct {
	SN         string                 `json:"sn"`
	DeviceId   string                 `json:"deviceId"`
	MessageId  string                 `json:"messageId"`
	Timestamp  int64                  `json:"timestamp"`
	Properties map[string]interface{} `json:"properties"`
}

type AlarmMessage struct {
	SN         string        `json:"sn"`
	DeviceId   string        `json:"deviceId"`
	MessageId  string        `json:"messageId"`
	Timestamp  int64         `json:"timestamp"`
	Type       string        `json:"type"`
	Title      string        `json:"title"`
	Message    string        `json:"message"`
	Properties []interface{} `json:"properties"`
}
