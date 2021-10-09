package model

type PropertyMessage struct {
	SN         string                  `json:"sn"`
	DeviceId   string                  `json:"deviceId"`
	MessageId  string                  `json:"messageId"`
	Timestamp  int64                   `json:"timestamp"`
	Properties map[string]PropertyItem `json:"properties"` //物模型属性列表
}
type PropertyItem struct {
	Key      string       `json:"key"`
	Name     string       `json:"name"`
	Value    interface{}  `json:"value"`
	DataType ItemDataType `json:"dataType"` //数据类型
}

type AlarmMessage struct {
	SN         string        `json:"sn"`
	DeviceId   string        `json:"deviceId"`
	MessageId  string        `json:"messageId"`
	Timestamp  int64         `json:"timestamp"`
	Type       string        `json:"type"` //event、alarm
	Title      string        `json:"title"`
	Message    string        `json:"message"`
	Properties []interface{} `json:"properties"` //告警属性列表
}
