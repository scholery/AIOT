package dto

type AlarmItem struct {
	AlarmId     string        `json:"alarmId"`
	AlarmSign   string        `json:"alarmSign"`
	AlarmTitle  string        `json:"alarmTitle"`
	AlarmType   string        `json:"alarmType"` //event、alarm
	AlarmLevel  string        `json:"alarmLevel"`
	DeviceId    string        `json:"deviceId"`
	DeviceName  string        `json:"deviceName"`
	DeviceSign  string        `json:"deviceSign"`
	ProductId   int           `json:"productId"`
	ProductName string        `json:"productName"`
	CreateTime  string        `json:"createTime"`
	MessageId   string        `json:"messageId"`
	Message     string        `json:"message"`
	Timestamp   int64         `json:"timestamp"`
	Props       []interface{} `json:"props"`      //告警属性
	Conditions  []interface{} `json:"conditions"` //条件
}
