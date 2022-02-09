package mqttClient

import (
	"errors"

	"koudai-box/iot/gateway/model"
	"koudai-box/iot/gateway/utils"

	"github.com/sirupsen/logrus"
)

type MqttClientDriver struct {
	Gateway *model.GatewayConfig
}

func (mqttClientDriver *MqttClientDriver) GetGatewayConfig() *model.GatewayConfig {
	return mqttClientDriver.Gateway
}
func (mqttClientDriver *MqttClientDriver) Start() error {
	return StartMQTT(mqttClientDriver.Gateway)
}

func (mqttClientDriver *MqttClientDriver) Stop() error {
	return StopMQTT(mqttClientDriver.Gateway)
}
func (mqttClientDriver *MqttClientDriver) StartDevicde(device *model.Device) error {
	return nil
}

func (mqttClientDriver *MqttClientDriver) StopDevicde(device *model.Device) error {
	return nil
}

func (mqttClientDriver *MqttClientDriver) FetchPropBatch(ts int64) (interface{}, error) {
	return nil, nil
}

func (mqttClientDriver *MqttClientDriver) FetchProp(device *model.Device, ts int64) (interface{}, error) {
	return nil, nil
}

//数据抽取接口
func (mqttClientDriver *MqttClientDriver) ExtracterProp(data interface{}, product *model.Product) (interface{}, error) {
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
func (mqttClientDriver *MqttClientDriver) TransformerProp(data interface{}, device *model.Device) (interface{}, error) {
	return utils.Transformer2DeviceProp(data, device)
}

//数据抽取接口,单条
func (mqttClientDriver *MqttClientDriver) FetchEvent(device *model.Device, ts int64) (interface{}, error) {
	return nil, nil
}

//网关数据发送设备
func (mqttClientDriver *MqttClientDriver) PostOperation(api model.ApiConfig, data interface{}, device *model.Device) (interface{}, error) {
	return nil, nil
}

//事件抽取接口
func (mqttClientDriver *MqttClientDriver) ExtracterEvent(data interface{}, product *model.Product) (interface{}, error) {
	return nil, nil
}
