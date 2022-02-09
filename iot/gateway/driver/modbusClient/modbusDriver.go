package modbusClient

import (
	"errors"
	"fmt"

	"koudai-box/iot/gateway/model"
	"koudai-box/iot/gateway/utils"

	"github.com/sirupsen/logrus"
)

type ModbusClientDriver struct {
	Gateway *model.GatewayConfig
}

func (modbusClientDriver *ModbusClientDriver) GetGatewayConfig() *model.GatewayConfig {
	return modbusClientDriver.Gateway
}
func (modbusClientDriver *ModbusClientDriver) Start() error {
	return ConnectModbus(modbusClientDriver.Gateway)
}

func (modbusClientDriver *ModbusClientDriver) Stop() error {
	return CloseModbus(modbusClientDriver.Gateway)
}
func (modbusClientDriver *ModbusClientDriver) StartDevicde(device *model.Device) error {
	return nil
}

func (modbusClientDriver *ModbusClientDriver) StopDevicde(device *model.Device) error {
	return nil
}

func (modbusClientDriver *ModbusClientDriver) FetchPropBatch(ts int64) (interface{}, error) {
	return nil, nil
}

func (modbusClientDriver *ModbusClientDriver) FetchProp(device *model.Device, ts int64) (interface{}, error) {
	if len(device.Product.Items) == 0 {
		logrus.Errorf("device[%s]'s item is null", device.Key)
		return nil, fmt.Errorf("device[%s]'s item is null", device.Key)
	}
	data := make(map[string]interface{})
	hasData := false
	for _, item := range device.Product.Items {
		val, err := QueryValue(modbusClientDriver.Gateway, device, item)
		if err != nil {
			continue
		}
		data[item.Source] = val
		hasData = true
	}
	if !hasData {
		return nil, fmt.Errorf("device[%s]'s data is null", device.Key)
	}
	return data, nil
}

//数据抽取接口
func (modbusClientDriver *ModbusClientDriver) ExtracterProp(data interface{}, product *model.Product) (interface{}, error) {
	if len(product.FunctionConfigs) == 0 {
		return nil, errors.New("function is null")
	}
	function, ok := product.FunctionConfigs[model.Function_Extract_Prop]
	if !ok || len(function.Function) == 0 {
		return nil, errors.New("function is null")
	}
	// logrus.Debugf("Extracter funtion name=%s,function[%s]", model.Function_Extract_Prop, function.Function)
	logrus.Debugf("Extracter funtion name=%s", model.Function_Extract_Prop)
	return utils.ExecJS(function.Function, function.Key, data)
}

//数据转换接口
func (modbusClientDriver *ModbusClientDriver) TransformerProp(data interface{}, device *model.Device) (interface{}, error) {
	return utils.Transformer2DeviceProp(data, device)
}

//数据抽取接口,单条
func (modbusClientDriver *ModbusClientDriver) FetchEvent(device *model.Device, ts int64) (interface{}, error) {
	return nil, nil
}

//网关数据发送设备
func (modbusClientDriver *ModbusClientDriver) PostOperation(api model.ApiConfig, data interface{}, device *model.Device) (interface{}, error) {
	return nil, nil
}

//事件抽取接口
func (modbusClientDriver *ModbusClientDriver) ExtracterEvent(data interface{}, product *model.Product) (interface{}, error) {
	return nil, nil
}
