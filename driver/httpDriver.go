package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	. "main/model"
	utils "main/utils"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type HttpDriver struct {
	Gateway GatewayConfig
	Device  Device
}

func (httpDriver *HttpDriver) FetchData() (interface{}, error) {
	api, ok := httpDriver.Gateway.ApiConfigs[API_GetProp]
	if !ok {
		return nil, errors.New("FetchData api is null")
	}
	logrus.Debug("httpDriver:", *httpDriver)
	address, err := url.Parse(fmt.Sprintf("http://%s:%d%s", httpDriver.Gateway.Ip, httpDriver.Gateway.Port, api.Value))
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
func (httpDriver *HttpDriver) Extracter(data interface{}) (interface{}, error) {
	if len(httpDriver.Device.Product.DataCombination) == 0 {
		return nil, errors.New("product is null")
	}
	if len(httpDriver.Device.Product.FunctionConfigs) == 0 {
		return nil, errors.New("function is null")
	}
	function, ok := httpDriver.Device.Product.FunctionConfigs[Function_Extract]
	if !ok {
		return nil, errors.New("function is null")
	}
	logrus.Debugf("Extracter funtion name=%s", Function_Extract)
	return utils.ExecJS(function.Function, function.Key, data)
}

//数据转换接口
func (httpDriver *HttpDriver) Transformer(data interface{}) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("transformer:data format error")
	}
	if len(httpDriver.Device.Product.Items) == 0 {
		return nil, errors.New("product model item is empty")
	}
	dataTmp := make(map[string]PropertyItem)
	for _, item := range httpDriver.Device.Product.Items {
		dataTmp[item.Key] = utils.GetPropertyItem(item, utils.GetMapValue(dataMap, item.Source))
	}

	return PropertyMessage{DeviceId: httpDriver.Device.Key, MessageId: utils.GetUUID(),
		Timestamp: time.Now().Unix(), Properties: dataTmp}, nil
}
