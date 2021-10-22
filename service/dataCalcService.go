package service

import (
	"errors"
	"fmt"
	"main/model"
	"main/utils"
	"time"

	"github.com/sirupsen/logrus"
)

const Message_queen_size = 10

type DataGatewayApi interface { //数据抽取接口
	//数据计算接口
	Calculater(data interface{}) (interface{}, error)
	//告警过滤接口
	FilterAlarm(data model.PropertyMessage) (interface{}, error)
	//数据存储接口
	LoaderMessage(data model.PropertyMessage) (interface{}, error)
	//数据存储接口
	LoaderAlarm(data model.AlarmMessage) (interface{}, error)
	//数据存储接口
	Push(interface{}, string) bool
}

type DataGateway struct {
	Device model.Device
}

/**
 *数据计算
 */
func (gateway *DataGateway) Calculater(data interface{}) (interface{}, error) {
	device := gateway.Device
	logrus.Debug(device)
	if len(device.Product.FunctionConfigs) == 0 {
		logrus.Errorf("Function_Calc is null,device id=%s", device.Key)
		return data, errors.New("function is null")
	}
	function, ok := device.Product.FunctionConfigs[model.Function_Calc]
	if !ok {
		logrus.Errorf("Function_Calc is null,device id=%s", device.Key)
		return data, errors.New("Function_Calc is null")
	}
	logrus.Debugf("Function_Calc funtion name=%s", model.Function_Calc)
	return utils.ExecJS(function.Function, function.Key, data)
}

/**
 *告警过滤
 */
func (gateway *DataGateway) FilterAlarm(data model.PropertyMessage) ([]interface{}, error) {
	logrus.Debugf("FilterAlarm,message id = %s", data.MessageId)
	device := gateway.Device
	if len(device.Product.AlarmConfigs) == 0 {
		logrus.Infof("AlarmConfigs is null,device id=%s", device.Key)
		return nil, errors.New("AlarmConfigs is null")
	}
	alarms := make([]interface{}, 0)
	for _, config := range device.Product.AlarmConfigs {
		if len(config.Conditions) == 0 {
			logrus.Infof("Conditions is null,device id=%s,alarm key=%s", device.Key, config.Key)
			continue
		}
		hasAlarm := false
		poprs := make([]interface{}, 0)
		for _, condition := range config.Conditions {
			match := utils.MatchContidion(data, condition)
			hasAlarm = hasAlarm || match
			if match {
				poprs = append(poprs, fmt.Sprintf("%s(%v %s %v)", condition.Name, data.Properties[condition.Key].Value, condition.Compare, condition.Value))
			}
		}
		if hasAlarm {
			alarms = append(alarms, model.AlarmMessage{SN: "", DeviceId: device.Key, MessageId: utils.GetUUID(),
				Timestamp: time.Now().Unix(), Type: config.Type, Title: config.Name, Message: config.Message, Properties: poprs})
		}
	}
	return alarms, nil

}

/**
 *消息存储
 */
func (gateway *DataGateway) LoaderProperty(data model.PropertyMessage) (bool, error) {
	logrus.Debugf("save prop,message id = %s", data.MessageId)
	return true, nil

}

/**
 *告警存储
 */
func (gateway *DataGateway) LoaderAlarm(data model.AlarmMessage) (bool, error) {
	logrus.Debugf("save alarm,message id = %s", data.MessageId)
	return true, nil
}

/**
 *数据推送
 */
func Push(data interface{}, router string) bool {
	return Public(data, router)
}
