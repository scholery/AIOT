package dto

import "koudai-box/iot/gateway/model"

type SaveGatewayRequest struct {
	GatewayId     int                `json:"gatewayId"`
	GatewayName   string             `json:"gatewayName"`
	Sign          string             `json:"sign"`
	Type          string             `json:"type"`
	Status        int                `json:"status"`
	Protocol      string             `json:"protocol"`      //'协议书'
	Ip            string             `json:"ip"`            //'IP'
	Port          int                `json:"port"`          //'端口'
	AuthInfo      []AuthInfoItem     `json:"authInfo"`      //'认证信息'
	Routers       []RouterItem       `json:"routers"`       //'路由定义'
	CollectType   string             `json:"collectType"`   //modbus\opcua\bacnet采集方式：定时、轮询
	CollectPeriod int                `json:"collectPeriod"` //modbus\opcua\bacnet采集周期:秒
	Cron          string             `json:"cron"`          //modbus\opcua\bacnet时间表达式
	ModbusConfig  model.ModbusConfig `json:"modbusConfig"`  //modbus
	Describe      string             `json:"describe"`      //'描述 '
}

type UpdateGatewayStatusRequest struct {
	GatewayId int `json:"gatewayId"`
	Status    int `json:"status"`
}

type DeleteGatewayRequest struct {
	GatewayId  int   `json:"gatewayId"`
	GatewayIds []int `json:"gatewayIds"`
}

type QueryGatewayDataRequest struct {
	GatewayName     string `json:"gatewayName"`
	GatewayType     string `json:"type"`
	GatewayProtocol string `json:"Protocol"`
	GatewayStatus   string `json:"status"`
	PageNo          int    `json:"pageNo"`
	PageSize        int    `json:"PageSize"`
}
