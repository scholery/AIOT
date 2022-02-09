package dto

import "koudai-box/iot/gateway/model"

type GatewayItem struct {
	GatewayId     int                `json:"gatewayId"`     //'网关id'
	GatewayName   string             `json:"gatewayName"`   //'网关名称'
	Sign          string             `json:"sign"`          //'标识'
	Type          string             `json:"type"`          //'类型'
	Status        int                `json:"status"`        //'状态'
	Protocol      string             `json:"protocol"`      //'协议书'
	ProtocolName  string             `json:"protocolName"`  //'协议名称'
	Ip            string             `json:"ip"`            //'IP'
	Port          int                `json:"port"`          //'端口'
	AuthInfo      []AuthInfoItem     `json:"authInfo"`      //'认证信息'
	Routers       []RouterItem       `json:"routers"`       //'路由定义'
	CollectType   string             `json:"collectType"`   //modbus\opcua\bacnet采集方式：定时、轮询
	CollectPeriod int                `json:"collectPeriod"` //modbus\opcua\bacnet采集周期:秒
	Cron          string             `json:"cron"`          //modbus\opcua\bacnet时间表达式
	ModbusConfig  model.ModbusConfig `json:"modbusConfig"`  //modbus
	Describe      string             `json:"describe"`      //'描述 '
	// CreateTime  time.Time `json:"createTime"` //'创建时间'
	// UpdateTime  time.Time `json:"updateTime"`
	CreateBy int `json:"createBy"`
}

type AuthInfoItem struct {
	Key   string `json:"key"`   //'key'
	Value string `json:"value"` //'value'
}

type RouterItem struct {
	RouterName      string `json:"routerName"`      //'路由名称'
	RouterUrl       string `json:"routerUrl"`       //'路由路径'
	RequestType     string `json:"requestType"`     //'请求方式'
	DataCombination string `json:"dataCombination"` //'数据格式'
	CollectType     string `json:"collectType"`     //'数据获取方式'
	CornExpression  string `json:"cornExpression"`  //'corn表达式'
	Period          int    `json:"period"`          //'周期'
}
