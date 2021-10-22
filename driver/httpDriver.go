package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"main/model"
	utils "main/utils"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

type HttpDriver struct {
	Gateway model.GatewayConfig
}

func (httpDriver *HttpDriver) FetchPropBatch() (interface{}, error) {
	return nil, errors.New("no implatement")
}

func (httpDriver *HttpDriver) FetchProp(device model.Device) (interface{}, error) {
	api, ok := httpDriver.Gateway.ApiConfigs[model.API_GetProp]
	if !ok {
		return nil, errors.New("FetchData api is null")
	}
	logrus.Debug("httpDriver:", *httpDriver)
	address, err := url.Parse(fmt.Sprintf("http://%s:%d%s", httpDriver.Gateway.Ip, httpDriver.Gateway.Port, api.Path))
	logrus.Debugf("request url=%s", address.String())
	if err != nil {
		logrus.Error("url err:", err)
		return nil, err
	}
	client := &http.Client{}
	req, _ := http.NewRequest(strings.ToUpper(api.Method), address.String(), nil)
	for _, param := range httpDriver.Gateway.Parameters {
		req.Header.Add(param.Key, param.Value)
	}
	//logrus.Debug("request:", req)
	resp, err := client.Do(req)
	//resp, err := http.Get(url)
	if err != nil {
		logrus.Error("http err:", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		logrus.Error("resp.StatusCode:", resp.StatusCode)
		return nil, fmt.Errorf("resp.StatusCode=%d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	//logrus.Debug("body:", string(body))
	if err != nil {
		logrus.Error("body read err:", err)
		return nil, err
	}
	//logrus.Debug("response body:", string(body))
	defer resp.Body.Close()
	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logrus.Error("err:", err)
		return nil, err
	}
	return data, err
}

//数据抽取接口
func (httpDriver *HttpDriver) ExtracterProp(data interface{}, device model.Device) (interface{}, error) {
	if len(device.Product.FunctionConfigs) == 0 {
		return nil, errors.New("function is null")
	}
	function, ok := device.Product.FunctionConfigs[model.Function_Extract]
	if !ok {
		return nil, errors.New("function is null")
	}
	logrus.Debugf("Extracter funtion name=%s", model.Function_Extract)
	return utils.ExecJS(function.Function, function.Key, data)
}

//数据转换接口
func (httpDriver *HttpDriver) TransformerProp(data interface{}, device model.Device) (interface{}, error) {
	return utils.Transformer2DeviceProp(data, device)
}

func (httpDriver *HttpDriver) GetCollectPeriod(key string) int {
	api, ok := httpDriver.Gateway.ApiConfigs[model.API_GetProp]
	if !ok {
		return -1
	}
	return api.CollectPeriod
}
