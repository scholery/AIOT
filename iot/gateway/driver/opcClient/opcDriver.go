package opcClient

import (
	"errors"
	"fmt"

	"koudai-box/iot/gateway/model"
	"koudai-box/iot/gateway/utils"

	"github.com/sirupsen/logrus"
)

type OPCClientDriver struct {
	Gateway *model.GatewayConfig
}

func (OPCClientDriver *OPCClientDriver) GetGatewayConfig() *model.GatewayConfig {
	return OPCClientDriver.Gateway
}
func (OPCClientDriver *OPCClientDriver) Start() error {
	return ConnectOPC(OPCClientDriver.Gateway)
}

func (OPCClientDriver *OPCClientDriver) Stop() error {
	return CloseOPC(OPCClientDriver.Gateway)
}
func (OPCClientDriver *OPCClientDriver) StartDevicde(device *model.Device) error {
	return nil
}

func (OPCClientDriver *OPCClientDriver) StopDevicde(device *model.Device) error {
	return nil
}

func (OPCClientDriver *OPCClientDriver) FetchPropBatch(ts int64) (interface{}, error) {
	return nil, nil
}

func (OPCClientDriver *OPCClientDriver) FetchProp(device *model.Device, ts int64) (interface{}, error) {
	if len(device.Product.Items) == 0 {
		logrus.Errorf("device[%s]'s item is null", device.Key)
		return nil, fmt.Errorf("device[%s]'s item is null", device.Key)
	}
	data := make(map[string]interface{})
	hasData := false
	for _, item := range device.Product.Items {
		val, err := QueryValue(OPCClientDriver.Gateway, device, item)
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
func (OPCClientDriver *OPCClientDriver) ExtracterProp(data interface{}, product *model.Product) (interface{}, error) {
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
func (OPCClientDriver *OPCClientDriver) TransformerProp(data interface{}, device *model.Device) (interface{}, error) {
	return utils.Transformer2DeviceProp(data, device)
}

//数据抽取接口,单条
func (OPCClientDriver *OPCClientDriver) FetchEvent(device *model.Device, ts int64) (interface{}, error) {
	return nil, nil
}

//网关数据发送设备
func (OPCClientDriver *OPCClientDriver) PostOperation(api model.ApiConfig, data interface{}, device *model.Device) (interface{}, error) {
	return nil, nil
}

//事件抽取接口
func (OPCClientDriver *OPCClientDriver) ExtracterEvent(data interface{}, product *model.Product) (interface{}, error) {
	return nil, nil
}
