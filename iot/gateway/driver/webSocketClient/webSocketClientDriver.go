package websocket

import (
	"encoding/json"
	"errors"

	"koudai-box/iot/gateway/model"
	"koudai-box/iot/gateway/utils"
)

type WebSocketClientDriver struct {
	Gateway *model.GatewayConfig
}

func (webSocketClientDriver *WebSocketClientDriver) GetGatewayConfig() *model.GatewayConfig {
	return webSocketClientDriver.Gateway
}

func (webSocketClientDriver *WebSocketClientDriver) Start() error {
	go StartWsClientConnect(webSocketClientDriver.Gateway)
	return nil
}

func (webSocketClientDriver *WebSocketClientDriver) Stop() error {
	return StopWsClientConnect(webSocketClientDriver.Gateway)
}
func (webSocketClientDriver *WebSocketClientDriver) StartDevicde(device *model.Device) error {
	return nil
}

func (webSocketClientDriver *WebSocketClientDriver) StopDevicde(device *model.Device) error {
	return nil
}

func (webSocketClientDriver *WebSocketClientDriver) FetchPropBatch(ts int64) (interface{}, error) {
	return nil, nil
}

func (webSocketClientDriver *WebSocketClientDriver) FetchProp(device *model.Device, ts int64) (interface{}, error) {
	return nil, nil
}

//数据抽取接口
func (webSocketClientDriver *WebSocketClientDriver) ExtracterProp(data interface{}, product *model.Product) (interface{}, error) {
	if len(product.FunctionConfigs) == 0 {
		return nil, errors.New("function is null")
	}
	function, ok := product.FunctionConfigs[model.Function_Extract_Prop]
	if !ok || len(function.Function) == 0 {
		return nil, errors.New("function is null")
	}
	// logrus.Debugf("Extracter funtion name=%s,function[%s]", model.Function_Extract_Prop, function.Function)
	return utils.ExecJS(function.Function, function.Key, data)
}

//数据转换接口
func (webSocketClientDriver *WebSocketClientDriver) TransformerProp(data interface{}, device *model.Device) (interface{}, error) {
	return utils.Transformer2DeviceProp(data, device)
}

//数据抽取接口,单条
func (webSocketClientDriver *WebSocketClientDriver) FetchEvent(device *model.Device, ts int64) (interface{}, error) {
	return nil, nil
}

//网关数据发送设备
func (webSocketClientDriver *WebSocketClientDriver) PostOperation(api model.ApiConfig, data interface{}, device *model.Device) (interface{}, error) {
	messageData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	SendWsMessage(api, string(messageData), webSocketClientDriver.Gateway.Key)
	return nil, nil
}

//事件抽取接口
func (webSocketClientDriver *WebSocketClientDriver) ExtracterEvent(data interface{}, product *model.Product) (interface{}, error) {
	if len(product.FunctionConfigs) == 0 {
		return nil, errors.New("function is null")
	}
	function, ok := product.FunctionConfigs[model.Function_Extract_Event]
	if !ok || len(function.Function) == 0 {
		return nil, errors.New("function is null")
	}
	// logrus.Debugf("Extracter funtion name=%s,function[%s]", model.Function_Extract_Event, function.Function)
	return utils.ExecJS(function.Function, function.Key, data)
}
