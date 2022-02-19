package coapserver

import (
	"errors"
	"fmt"
	"net/http"

	"koudai-box/iot/gateway/model"
	"koudai-box/iot/gateway/utils"

	"github.com/sirupsen/logrus"
)

var cache_servers map[int]*http.Server = make(map[int]*http.Server)

type CoapServerDriver struct {
	Gateway *model.GatewayConfig
}

func (coapDriver *CoapServerDriver) GetGatewayConfig() *model.GatewayConfig {
	return coapDriver.Gateway
}

func (coapServerDrvier *CoapServerDriver) Start() error {
	gateway := coapServerDrvier.GetGatewayConfig()
	if _, ok := cache_servers[gateway.Id]; ok {
		return fmt.Errorf("http server[%s] has been started", gateway.Key)
	}
	err := RegisterURL()
	if err != nil {
		return err
	}
	logrus.Infof("++++++++++++++开启coap server[%s]应用[%s:%d]++++++++++++++", gateway.Key, gateway.Ip, gateway.Port)
	return nil
}

func (coapServerDrvier *CoapServerDriver) Stop() error {
	return nil
}

func (coapServerDrvier *CoapServerDriver) StartDevicde(device *model.Device) error {
	return nil
}

func (coapServerDrvier *CoapServerDriver) StopDevicde(device *model.Device) error {
	return nil
}

//数据抽取接口
func (coapServerDrvier *CoapServerDriver) ExtracterProp(data interface{}, product *model.Product) (interface{}, error) {
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
func (CoapServerDriver *CoapServerDriver) TransformerProp(data interface{}, device *model.Device) (interface{}, error) {
	return utils.Transformer2DeviceProp(data, device)
}

//数据抽取接口,单条
func (CoapServerDriver *CoapServerDriver) FetchEvent(device *model.Device, ts int64) (interface{}, error) {
	return nil, nil
}

//网关数据发送设备
func (CoapServerDriver *CoapServerDriver) PostOperation(api model.ApiConfig, data interface{}, device *model.Device) (interface{}, error) {
	// messageData, err := json.Marshal(data)
	// if err != nil {
	// 	return nil, err
	// }
	// SendWsMessage(string(messageData), CoapServerDriver.Gateway.Key)
	return nil, nil
}

//事件抽取接口
func (CoapServerDriver *CoapServerDriver) ExtracterEvent(data interface{}, product *model.Product) (interface{}, error) {
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
