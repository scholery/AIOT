package dto

type EventItem struct {
	EventId     int    `json:"eventId"`
	EventSign   string `json:"eventSign"`
	EventTitle  string `json:"eventTitle"`
	EventType   string `json:"eventType"` //event、event
	EventLevel  string `json:"eventLevel"`
	DeviceId    string `json:"deviceId"`
	DeviceName  string `json:"deviceName"`
	DeviceSign  string `json:"deviceSign"`
	ProductId   int    `json:"productId"`
	ProductName string `json:"productName"`
	CreateTime  string `json:"createTime"`
	MessageId   string `json:"messageId"`
	Message     string `json:"message"`
	Timestamp   int64  `json:"timestamp"`
}
