package driver

import (
	httpClient "koudai-box/iot/gateway/driver/httpClient"
	httpServer "koudai-box/iot/gateway/driver/httpServer"
	modbusClient "koudai-box/iot/gateway/driver/modbusClient"
	mqttClient "koudai-box/iot/gateway/driver/mqttClient"
	opcClient "koudai-box/iot/gateway/driver/opcClient"
	webSocketClient "koudai-box/iot/gateway/driver/webSocketClient"
	model "koudai-box/iot/gateway/model"
)

type Driver interface {
	//启动
	Start() error
	//停止
	Stop() error
	//启动
	StartDevicde(device *model.Device) error
	//停止
	StopDevicde(device *model.Device) error
	//数据抽取接口,批量
	FetchPropBatch(ts int64) (interface{}, error)
	//数据抽取接口,单条
	FetchProp(device *model.Device, ts int64) (interface{}, error)
	//数据抽取接口
	ExtracterProp(data interface{}, product *model.Product) (interface{}, error)
	//数据转换接口
	TransformerProp(data interface{}, device *model.Device) (interface{}, error)
	//数据抽取接口,单条
	FetchEvent(device *model.Device, ts int64) (interface{}, error)
	//事件抽取接口
	ExtracterEvent(data interface{}, product *model.Product) (interface{}, error)
	//网关数据发送设备
	PostOperation(api model.ApiConfig, data interface{}, device *model.Device) (interface{}, error)
	//获取网关配置
	GetGatewayConfig() *model.GatewayConfig
}

func GetDriver(gateway *model.GatewayConfig) (Driver, bool) {
	var driver Driver
	ok := true
	switch gateway.Protocol {
	case model.Geteway_Protocol_HTTP_Client:
		driver = &httpClient.HttpDriver{Gateway: gateway}
	case model.Geteway_Protocol_HTTP_Server:
		driver = &httpServer.HttpServerDriver{Gateway: gateway}
	case model.Geteway_Protocol_WebSocket_Client:
		driver = &webSocketClient.WebSocketClientDriver{Gateway: gateway}
	case model.Geteway_Protocol_WebSocket_Server:
	case model.Geteway_Protocol_MQTT:
		driver = &mqttClient.MqttClientDriver{Gateway: gateway}
	case model.Geteway_Protocol_MQTTSN:
	case model.Geteway_Protocol_ModbusTCP, model.Geteway_Protocol_ModbusRTU:
		driver = &modbusClient.ModbusClientDriver{Gateway: gateway}
	case model.Geteway_Protocol_OPCUA:
		driver = &opcClient.OPCClientDriver{Gateway: gateway}
	case model.Geteway_Protocol_CoAP:
	case model.Geteway_Protocol_LwM2M:
	case model.Geteway_Protocol_BACnet_IP:
	default:
	}
	return driver, ok
}
