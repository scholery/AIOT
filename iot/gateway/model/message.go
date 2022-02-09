package model

//properties：设备属性，对应PropertyMessage
//alarm：AIOT计算事件/告警，对应AlarmMessage
//event：子系统上报事件/告警，对应EventMessage
const (
	Message_Type_Prop      = "properties"
	Message_Type_Iot_Event = "iot_event"
	Message_Type_Event     = "event"
	Message_Type_Device    = "device"
)

//消息统一封装
type Message struct {
	SN    string      `json:"sn"`
	Type  string      `json:"type"`
	MsgId string      `json:"msgId"`
	Msg   interface{} `json:"msg"`
}

//设备基本数据
type DeviceMessage struct {
	DeviceId       int                    `json:"deviceId"`
	Key            string                 `json:"key"`
	Name           string                 `json:"name"`
	SourceId       string                 `json:"sourceId"` //物理设备标识
	Category       string                 `json:"category"` //分类
	Geo            string                 `json:"geo"`
	Locale         string                 `json:"locale"`         //位置描述
	ActivateStatus int                    `json:"activateStatus"` //状态 0 未激活   1 激活
	RunningStatus  int                    `json:"runningStatus"`  //是否在线 0 不在线 1 在线
	ExtProps       map[string]interface{} `json:"extProps"`       //扩展属性
	Desc           string                 `json:"desc"`           //排序
	UpdateTime     string                 `json:"updateTime"`     //更新时间
	DelFlag        int                    `json:"delFlag"`        //删除标识 0 未删除  1 删除
}

//设备属性
type PropertyMessage struct {
	DeviceId   string                  `json:"deviceId"`
	DeviceSign string                  `json:"deviceSign"`
	MessageId  string                  `json:"messageId"`
	Timestamp  int64                   `json:"timestamp"`
	Properties map[string]PropertyItem `json:"properties"` //物模型属性列表
}
type PropertyItem struct {
	//Key      string       `json:"key"`
	Code     string       `json:"code"`
	Name     string       `json:"name"`
	Value    interface{}  `json:"value"`
	DataType ItemDataType `json:"dataType"` //数据类型
}

//AIOT计算的事件/告警
type IotEventMessage struct {
	Code       string         `json:"code"`
	DeviceId   string         `json:"deviceId"`
	DeviceSign string         `json:"deviceSign"`
	MessageId  string         `json:"messageId"`
	Timestamp  int64          `json:"timestamp"`
	Type       string         `json:"type"` //event、alarm
	Level      string         `json:"level"`
	Title      string         `json:"title"`
	Message    string         `json:"message"`
	Conditions []Condition    `json:"conditions"` //条件
	Properties []PropertyItem `json:"properties"` //物模型属性列表
}

//接入设备上报的事件/告警
type EventMessage struct {
	DeviceId   string      `json:"deviceId"`
	DeviceSign string      `json:"deviceSign"`
	MessageId  string      `json:"messageId"`
	Timestamp  int64       `json:"timestamp"`
	Type       string      `json:"type"` //event、alarm
	Title      string      `json:"title"`
	Message    string      `json:"message"`
	Properties interface{} `json:"properties"` //物模型属性列表
}
