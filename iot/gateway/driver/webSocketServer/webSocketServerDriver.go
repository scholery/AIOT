package websocket

import (
	"errors"

	"koudai-box/iot/gateway/model"
	"koudai-box/iot/gateway/utils"

	"github.com/sirupsen/logrus"
)

type WebSocketServerDriver struct {
	Gateway *model.GatewayConfig
}

func (webSocketServerDriver *WebSocketServerDriver) GetGatewayConfig() *model.GatewayConfig {
	return webSocketServerDriver.Gateway
}

func (webSocketServerDriver *WebSocketServerDriver) Start() error {
	go StartWsServerConnect(webSocketServerDriver.Gateway)
	return nil
}

func (webSocketServerDriver *WebSocketServerDriver) Stop() error {
	go StopWsServerConnect(webSocketServerDriver.Gateway)
	return nil
}

func (webSocketServerDriver *WebSocketServerDriver) StartDevicde(device *model.Device) error {
	return nil
}

func (webSocketServerDriver *WebSocketServerDriver) StopDevicde(device *model.Device) error {
	return nil
}

func (webSocketServerDriver *WebSocketServerDriver) FetchPropBatch(ts int64) (interface{}, error) {
	return nil, nil
}

func (webSocketServerDriver *WebSocketServerDriver) FetchProp(device *model.Device, ts int64) (interface{}, error) {
	return nil, nil
}

//数据抽取接口
func (webSocketServerDriver *WebSocketServerDriver) ExtracterProp(data interface{}, product *model.Product) (interface{}, error) {
	if len(product.FunctionConfigs) == 0 {
		return nil, errors.New("function is null")
	}
	function, ok := product.FunctionConfigs[model.Function_Extract_Prop]
	if !ok || len(function.Function) == 0 {
		return nil, errors.New("function is null")
	}
	logrus.Debugf("Extracter funtion name=%s,function[%s]", model.Function_Extract_Prop, function.Function)
	return utils.ExecJS(function.Function, function.Key, data)
}

//数据转换接口
func (webSocketServerDriver *WebSocketServerDriver) TransformerProp(data interface{}, device *model.Device) (interface{}, error) {
	return utils.Transformer2DeviceProp(data, device)
}

//数据抽取接口,单条
func (webSocketServerDriver *WebSocketServerDriver) FetchEvent(device *model.Device, ts int64) (interface{}, error) {
	return nil, nil
}

//网关数据发送设备
func (webSocketServerDriver *WebSocketServerDriver) PostOperation(api model.ApiConfig, data interface{}, device *model.Device) (interface{}, error) {
	// messageData, err := json.Marshal(data)
	// if err != nil {
	// 	return nil, err
	// }
	// SendWsMessage(api, string(messageData), webSocketClientDriver.Gateway.Key)
	return nil, nil
}

//事件抽取接口
func (webSocketServerDriver *WebSocketServerDriver) ExtracterEvent(data interface{}, product *model.Product) (interface{}, error) {
	if len(product.FunctionConfigs) == 0 {
		return nil, errors.New("function is null")
	}
	function, ok := product.FunctionConfigs[model.Function_Extract_Event]
	if !ok || len(function.Function) == 0 {
		return nil, errors.New("function is null")
	}
	logrus.Debugf("Extracter funtion name=%s,function[%s]", model.Function_Extract_Event, function.Function)
	return utils.ExecJS(function.Function, function.Key, data)
}
