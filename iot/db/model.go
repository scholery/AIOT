package db

import (
	"time"
)

// "github.com/astaxie/beego/orm"

type Gateway struct {
	Id            int       `orm:"auto"`                        //'网关id'
	Name          string    `json:"name"`                       //'网关名称'
	Sign          string    `json:"sign"`                       //'标识'
	Status        int       `json:"status"`                     //'状态'
	Protocol      string    `json:"protocol"`                   //'协议书'
	Ip            string    `json:"ip"`                         //'IP'
	Port          int       `json:"port"`                       //'端口'
	AuthInfo      string    `json:"authInfo" orm:"type(text)"`  //'认证信息'
	Routers       string    `json:"routers" orm:"type(text)"`   //'路由定义'
	CollectType   string    `json:"collectType"`                //modbus\opcua\bacnet采集方式：定时、轮询
	CollectPeriod int       `json:"collectPeriod"`              //modbus\opcua\bacnet采集周期:秒
	Cron          string    `json:"cron"`                       //modbus\opcua\bacnet时间表达式
	ModbusConfig  string    `json:"modbusConfig"`               //modbus
	Describe      string    `json:"describe" orm:"type(text)"`  //'描述 '
	CreateTime    time.Time `orm:"auto_now_add;type(datetime)"` //'创建时间'
	UpdateTime    time.Time `orm:"auto_now;type(datetime)"`
	CreateBy      int       //'创建用户id'
	UpdateBy      int
	DelFlag       int //0:默认1:删除'
}

type DeviceProperty struct {
	Id         int       `json:"id"` //'记录id'
	DeviceId   string    `json:"deviceId"`
	MessageId  string    `json:"messageId"`
	Timestamp  int64     `json:"timestamp"`
	Properties string    `json:"properties" orm:"type(text)"` //物模型属性列表
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`  //'创建时间'
	PushFlag   int       //0:默认1:已推送'
}

type Alarm struct {
	Id          int       `json:"id"` //'告警id'
	MessageId   string    `json:"messageId"`
	ProductId   int       `json:"productId"` //'产品id'
	ProductName string    `json:"productName"`
	DeviceId    string    `json:"deviceId"`   //'设备id'
	DeviceName  string    `json:"deviceName"` //'设备name'
	DeviceSign  string    `json:"deviceSign"` //'设备标识'
	Timestamp   int64     `json:"timestamp"`
	Code        string    `json:"code"`                       //'告警标识'
	Title       string    `json:"title"`                      //'告警名称'
	Type        string    `json:"type"`                       //'告警类型'
	Level       string    `json:"level"`                      //'告警等级'
	CreateTime  time.Time `orm:"auto_now_add;type(datetime)"` //'创建时间'
	Message     string    `json:"message" orm:"type(text)"`
	Properties  string    `json:"properties" orm:"type(text)"` //属性列表
	Conditions  string    `json:"conditions" orm:"type(text)"` //条件列表
	PushFlag    int       //0:默认1:已推送'
}

type Event struct {
	Id          int       `orm:"auto"`                        //'事件id'
	DeviceId    string    `json:"deviceId"`                   //'设备id'
	DeviceName  string    `json:"deviceName"`                 //'设备名'
	DeviceSign  string    `json:"deviceSign"`                 //'设备标识'
	ProductId   int       `json:"productId"`                  //'产品id'
	ProductName string    `json:"productName"`                //'产品名'
	CreateTime  time.Time `orm:"auto_now_add;type(datetime)"` //'创建时间'
	Title       string    `json:"title"`                      //'事件名称'
	Sign        string    `json:"sign"`                       //'事件标识'
	Type        string    `json:"type"`                       //'事件类型'
	Level       string    `json:"level"`                      //'事件等级'
	Timestamp   int64     `json:"timestamp"`
	MessageId   string    `json:"messageId"`
	Message     string    `json:"message" orm:"type(text)"`
	Properties  string    `json:"properties" orm:"type(text)"`
	PushFlag    int       //0:默认1:已推送'
}

type Dict struct {
	Id    int    `orm:"auto" json:"id"`
	Pid   int    `json:"pid"`
	Type  string `orm:"size(20)" json:"type"`
	Sort  int    `json:"sort"`
	Name  string `orm:"size(50)" json:"name"`
	Value string `orm:"size(50)" json:"value"`
	Tips  string `orm:"size(50)" json:"tips"`
	Extra string `orm:"size(20)" json:"extra"` //预留字段
}

func (d *Dict) TableUnique() [][]string {
	return [][]string{
		{"Type", "Value"},
	}
}

//产品信息表
type Product struct {
	Id               int        `orm:"auto" json:"id"`
	Code             string     `json:"code"`                                                   //产品标识
	Name             string     `json:"name"`                                                   //产品名称
	Category         string     `json:"category"`                                               //分类
	Desc             string     `json:"desc"`                                                   //产品说明
	Image            string     `json:"image"`                                                  //产品logo
	Items            string     `json:"items" orm:"auto_now_add;default([])"`                   //物模型
	EventConfigs     string     `json:"eventConfigs" orm:"auto_now_add;default([]);type(text)"` //函数配置
	OperationConfigs string     `json:"operationConfigs" orm:"type(text)"`                      //操作定义
	AlarmConfigs     string     `json:"alarmSettings" orm:"type(text)"`                         //告警配置
	FunctionConfigs  string     `json:"functionConfigs" orm:"type(text)"`                       //函数配置
	GatewayId        int        `json:"gatewayId"`                                              //对应网关
	CreateTime       time.Time  `json:"createTime" orm:"auto_now_add;type(datetime)"`           //'创建时间'
	PublishTime      *time.Time `json:"publishTime" orm:"column(publish_time);type(datetime)"`  //'发布时间'
	State            int        `json:"state"`                                                  //状态 0-未发布,1-发布
	DelFlag          int        `json:"delFlag" orm:"default(0)"`                               //0:默认1:删除'
}

//设备信息表
type Device struct {
	Id             int       `orm:"auto" json:"id"`
	Code           string    `json:"code"`                                         //设备标识
	Name           string    `json:"name"`                                         //设备名称
	SourceId       string    `json:"sourceId"`                                     //设备源标识，modbus：slaveid，opc ua：namespace。json：property key
	Geo            string    `json:"geo"`                                          //经纬度
	Locale         string    `json:"locale"`                                       //位置描述
	ProductId      int       `json:"productId"`                                    //产品id
	ActivateStatus int       `json:"activateStatus"`                               //状态 0 未激活   1 激活
	RunningStatus  int       `json:"runningStatus" orm:"default(-1)"`              //是否在线 0 不在线 1 在线
	Desc           string    `json:"desc" orm:"type(text)"`                        //排序
	ExtProps       string    `json:"extProps" orm:"type(text)"`                    //扩展属性
	CreateTime     time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"` //创建时间
	UpdateTime     time.Time `json:"updateTime" orm:"auto_now;type(datetime)"`     //更新时间
	DelFlag        int       `json:"delFlag" orm:"default(0)"`                     //删除标识 0 未删除  1 删除
}

//设备信息表
type GatewayApiHistory struct {
	Id         int       `orm:"auto" json:"id"`
	Name       string    `json:"name"`
	GatewayId  int       `json:"gatewayId"`
	DeviceId   int       `json:"deviceId"`
	Timestamp  int64     `json:"timestamp"`
	UpdateTime time.Time `json:"updateTime"` //时间
	Status     int       `json:"status"`     //状态：0 失败，1:成功
}

//设备信息表
type OperationRecord struct {
	Id         int       `orm:"auto" json:"id"`
	GatewayId  int       `json:"gatewayId"`
	ProductId  int       `json:"productId"` //产品ID
	DeviceId   int       `json:"deviceId"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	CreateTime time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"` //创建时间
	Type       string    `json:"type"`                                         //异步同步
	Desc       string    `json:"desc"`                                         //描述
	Inputs     string    `json:"inputs" orm:"type(text)"`                      //入口参数
	Outputs    string    `json:"outputs" orm:"type(text)"`                     //出口参数
}

type DeviceStatus struct {
	DeviceId     int    `json:"deviceId" orm:"pk"`
	Status       int    `json:"status"` //状态
	Timestamp    int64  `json:"timestamp" orm:"type(text)"`
	LastStatus   string `json:"lastStatus" orm:"type(text)"`   //最新
	ZeroStatus   string `json:"zeroStatus" orm:"type(text)"`   //初始值
	PreDayStatus string `json:"preDayStatus" orm:"type(text)"` //前一天平均
}
