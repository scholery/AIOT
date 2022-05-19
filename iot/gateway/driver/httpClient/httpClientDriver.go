package httpClient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"koudai-box/iot/gateway/model"
	utils "koudai-box/iot/gateway/utils"

	"github.com/sirupsen/logrus"
)

type HttpDriver struct {
	Gateway *model.GatewayConfig
}

func (httpDriver *HttpDriver) Start() error {
	return nil
}

func (httpDriver *HttpDriver) Stop() error {
	return nil
}

func (httpDriver *HttpDriver) StartDevicde(device *model.Device) error {
	return nil
}

func (httpDriver *HttpDriver) StopDevicde(device *model.Device) error {
	return nil
}

func (httpDriver *HttpDriver) GetGatewayConfig() *model.GatewayConfig {
	return httpDriver.Gateway
}

func fetchData(address string, method string, params []model.Parameter) (interface{}, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(strings.ToUpper(method), address, nil)
	for _, param := range params {
		if len(param.Key) == 0 {
			continue
		}
		req.Header.Add(param.Key, param.Value)
	}
	logrus.Debug("request:", req)
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("http err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logrus.Errorf("addr:%s,resp.StatusCode:%d", address, resp.StatusCode)
		return nil, fmt.Errorf("resp.StatusCode=%d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	//logrus.Debug("body:", string(body))
	if err != nil {
		logrus.Error("body read err:", err)
		return nil, err
	}
	logrus.Debug("response body:", string(body))
	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logrus.Error("err:", err)
		return nil, err
	}
	return data, err
}

func postData(address string, method string, data interface{}, params []model.Parameter) (interface{}, error) {
	client := &http.Client{}
	var buf *bytes.Buffer = nil
	if nil != data {
		body, err := json.Marshal(data)
		if err != nil {
			logrus.Errorf("api[%s]'s body parse err:", address, err)
			return nil, err
		}
		buf = bytes.NewBuffer(body)
	}

	req, _ := http.NewRequest(strings.ToUpper(method), address, buf)
	for _, param := range params {
		if len(param.Key) == 0 {
			continue
		}
		req.Header.Add(param.Key, param.Value)
	}
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	logrus.Debug("request:", req)
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("http err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logrus.Errorf("api[%s]'s resp.StatusCode[%d]:", address, resp.StatusCode)
		return nil, fmt.Errorf("api[%s]'s resp.StatusCode[%d]", address, resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("body read err:", err)
		return nil, err
	}
	logrus.Debug("response body:", string(respBody))
	var retData interface{}
	err = json.Unmarshal(respBody, &retData)
	if err != nil {
		logrus.Error("err:", err)
		return nil, err
	}
	return retData, nil
}

func (httpDriver *HttpDriver) FetchPropBatch(ts int64) (interface{}, error) {
	api, ok := httpDriver.Gateway.ApiConfigs[model.API_GetProp]
	if !ok {
		return nil, fmt.Errorf("%s api is null", model.API_GetProp)
	}
	var urlStr string
	if strings.Index(api.Path, "/") != 0 {
		urlStr = fmt.Sprintf("http://%s:%d%s/", httpDriver.Gateway.Ip, httpDriver.Gateway.Port, api.Path)
	} else {
		urlStr = fmt.Sprintf("http://%s:%d%s", httpDriver.Gateway.Ip, httpDriver.Gateway.Port, api.Path)
	}
	dataMap := make(map[string]interface{})
	dataMap["ts"] = ts
	path := utils.ParseTpl(urlStr, dataMap)
	address, err := url.Parse(path)
	if err != nil {
		logrus.Error("url err:", err)
		return nil, err
	}
	return fetchData(address.String(), api.Method, httpDriver.Gateway.Parameters)
}

func (httpDriver *HttpDriver) FetchProp(device *model.Device, ts int64) (interface{}, error) {
	api, ok := httpDriver.Gateway.ApiConfigs[model.API_GetProp]
	if !ok {
		return nil, fmt.Errorf("%s api is null", model.API_GetProp)
	}
	dataMap := utils.ToMap(device)
	dataMap["ts"] = ts
	path := utils.ParseTpl(api.Path, dataMap)
	address, err := url.Parse(fmt.Sprintf("http://%s:%d%s", httpDriver.Gateway.Ip, httpDriver.Gateway.Port, path))
	if err != nil {
		logrus.Error("url err:", err)
		return nil, err
	}
	return fetchData(address.String(), api.Method, httpDriver.Gateway.Parameters)
}

//数据抽取接口
func (httpDriver *HttpDriver) ExtracterProp(data interface{}, product *model.Product) (interface{}, error) {
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
func (httpDriver *HttpDriver) TransformerProp(data interface{}, device *model.Device) (interface{}, error) {
	return utils.Transformer2DeviceProp(data, device)
}

//数据抽取接口,单条
func (httpDriver *HttpDriver) FetchEvent(device *model.Device, ts int64) (interface{}, error) {
	api, ok := httpDriver.Gateway.ApiConfigs[model.API_GetEvent]
	if !ok {
		return nil, fmt.Errorf("%s api is null", model.API_GetEvent)
	}
	var urlStr string
	if strings.Index(api.Path, "/") != 0 {
		urlStr = fmt.Sprintf("http://%s:%d%s/", httpDriver.Gateway.Ip, httpDriver.Gateway.Port, api.Path)
	} else {
		urlStr = fmt.Sprintf("http://%s:%d%s", httpDriver.Gateway.Ip, httpDriver.Gateway.Port, api.Path)
	}
	var dataMap map[string]interface{}
	if nil != device {
		dataMap = utils.ToMap(device)
	} else {
		dataMap = make(map[string]interface{})
	}
	dataMap["ts"] = ts
	path := utils.ParseTpl(urlStr, dataMap)
	address, err := url.Parse(path)
	if err != nil {
		logrus.Errorf("api[%s]'s url err:", model.API_GetEvent, err)
		return nil, err
	}
	return fetchData(address.String(), api.Method, httpDriver.Gateway.Parameters)
}

//网关数据发送设备
func (httpDriver *HttpDriver) PostOperation(api model.ApiConfig, data interface{}, device *model.Device) (interface{}, error) {
	if len(api.Path) == 0 {
		return nil, fmt.Errorf("api[%s]'s path is null", api.Name)
	}
	var urlStr string
	if strings.Index(api.Path, "/") != 0 {
		urlStr = fmt.Sprintf("http://%s:%d%s/", httpDriver.Gateway.Ip, httpDriver.Gateway.Port, api.Path)
	} else {
		urlStr = fmt.Sprintf("http://%s:%d%s", httpDriver.Gateway.Ip, httpDriver.Gateway.Port, api.Path)
	}
	var dataMap map[string]interface{}
	if nil != device {
		dataMap = utils.ToMap(device)
	} else {
		dataMap = make(map[string]interface{})
	}
	path := utils.ParseTpl(urlStr, dataMap)
	address, err := url.Parse(path)
	if err != nil {
		logrus.Errorf("api[%s]'s url err:", api.Name, err)
		return nil, err
	}
	return postData(address.String(), api.Method, data, httpDriver.Gateway.Parameters)
}

//事件抽取接口
func (httpDriver *HttpDriver) ExtracterEvent(data interface{}, product *model.Product) (interface{}, error) {
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
