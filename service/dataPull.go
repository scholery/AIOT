package service

import (
	"errors"
	"fmt"
	"main/driver"
	"main/model"
	"main/utils"
	"time"

	"github.com/sirupsen/logrus"
)

const Message_queen_size = 10

//设备属性消息通道
var propMessChan = make(chan PropertyChan, Message_queen_size)

//运行标记
var run bool

type DataGatewayApi interface { //数据抽取接口
	//数据计算接口
	Calculater(data interface{}) (interface{}, error)
	//告警过滤接口
	Filter(data model.PropertyMessage) (interface{}, error)
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

type PropertyChan struct {
	PropertyMessage model.PropertyMessage
	Device          model.Device
}

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

func (gateway *DataGateway) Filter(data model.PropertyMessage) ([]interface{}, error) {
	logrus.Debugf("Filter,message id = %s", data.MessageId)
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
				poprs = append(poprs, fmt.Sprintf("%s[%d] %s %s", condition.Name, data.Properties[condition.Key].Value, condition.Compare, condition.Value))
			}
		}
		if hasAlarm {
			alarms = append(alarms, model.AlarmMessage{SN: "", DeviceId: device.Key, MessageId: utils.GetUUID(),
				Timestamp: time.Now().Unix(), Title: config.Name, Message: config.Message, Properties: poprs})
		}
	}
	return alarms, nil

}

func (gateway *DataGateway) LoaderMessage(data model.PropertyMessage) (bool, error) {
	logrus.Debugf("save prop,message id = %s", data.MessageId)
	return true, nil

}

func (gateway *DataGateway) LoaderAlarm(data model.AlarmMessage) (bool, error) {
	logrus.Debugf("save alarm,message id = %s", data.MessageId)
	return true, nil

}
func Push(data interface{}, router string) bool {
	return Public(data, router)
}

//设备运行状态控制
var deviceThreads map[string]bool = make(map[string]bool)

//设备状态
var deviceProps map[string]model.PropertyMessage = make(map[string]model.PropertyMessage)

/**
 *开启所有属性获取及事件监听
**/
func StartPull() {
	logrus.Info("StartPull")
	gatewayConfigs, ok := GetGatewayConfigs()
	if !ok {
		logrus.Error("GetGatewayConfigs is null")
		return
	}
	for _, gateway := range gatewayConfigs {
		StartGatewayPull(gateway)
	}
	run = true
	for run {
		select {
		case tmp := <-propMessChan:
			go ExecCalc(tmp.PropertyMessage, tmp.Device)
		default:
		}
	}
}

/**
 *停止所有属性获取及事件监听
**/
func StopPull() {
	logrus.Info("StopPull")
	run = false
	for k := range deviceThreads {
		delete(deviceThreads, k)
	}
}

func StartGatewayPull(gateway model.GatewayConfig) {
	logrus.Infof("StartGatewayPull GetGateway %s", gateway.Key)
	devices, ok := GetDevices(gateway.Key)
	if !ok {
		logrus.Error("GetGateway %s's related device is null", gateway.Key)
		return
	}
	for _, device := range devices {
		StartDevicePull(gateway, device)
	}

}

func StopGatewayPull(gateway model.GatewayConfig) bool {
	logrus.Infof("StopGatewayPull GetGateway %s", gateway.Key)
	devices, ok := GetDevices(gateway.Key)
	if !ok {
		logrus.Error("GetGateway %s's related device is null", gateway.Key)
		return false
	}
	for _, device := range devices {
		StopDevicePull(gateway, device)
	}
	return true
}

func StartDevicePull(gateway model.GatewayConfig, device model.Device) {
	logrus.Infof("StartDevicePull device %s", device.Key)
	deviceThreads[device.Key] = true
	driver, ok := driver.GetDriver(gateway, device)
	if !ok {
		logrus.Errorf("driver init failed,type is %s,device is %s", gateway.Protocol, device.Key)
		return
	}
	ExecPull(driver, device)
}

func StopDevicePull(gateway model.GatewayConfig, device model.Device) bool {
	logrus.Infof("StopDevicePull device %s", device.Key)
	delete(deviceThreads, device.Key)
	return true
}

func ExecPull(driver driver.Driver, device model.Device) {
	if st, ok := deviceThreads[device.Key]; !ok || !st {
		logrus.Infof("Stop ExecPull device %s", device.Key)
		return
	}
	logrus.Infof("ExecPull device %s", device.Key)
	start := time.Now() // 获取当前时间
	data, err := driver.FetchData()
	if err != nil {
		logrus.Errorf("FetchData error,device is %s,err:%s ", device.Key, err)
		return
	}
	data, err = driver.Extracter(data)
	if err != nil {
		logrus.Errorf("Extracter data error,device is %s,data:\r\n%s ", device.Key, data)
		logrus.Error(err)
	}
	data, err = driver.Transformer(data)
	if err != nil {
		logrus.Errorf("Transformer data error,device is %s,data:\r\n%s ", device.Key, data)
		logrus.Error(err)
		return
	}
	elapsed := time.Since(start)
	logrus.Debug("ExecPull执行完成耗时：", elapsed)
	tmp, ok := data.(model.PropertyMessage)
	if ok {
		propMessChan <- PropertyChan{tmp, device}
	}
	time.Sleep(time.Duration(device.Product.CollectPeriod) * time.Second)
	go ExecPull(driver, device)
	logrus.Debug("下一次执行")
}

func ExecCalc(data model.PropertyMessage, device model.Device) {
	start := time.Now() // 获取当前时间
	dataGateway := &DataGateway{Device: device}
	res, err := dataGateway.Calculater(data)
	if err != nil {
		logrus.Error("Calculater error.", err)
	}
	tmpP, ok := res.(model.PropertyMessage)
	if !ok {
		logrus.Error("calc error")
		return
	}
	dataGateway.LoaderMessage(tmpP)
	//变化上报
	old, ok := deviceProps[device.Key]
	if !ok || HasChange(device, old, tmpP) {
		Public(tmpP, Router_prop)
	} else {
		logrus.Info("no change")
	}
	//更新最新状态
	deviceProps[device.Key] = tmpP
	alarms, err := dataGateway.Filter(tmpP)
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, alarm := range alarms {
		tmpA, ok := alarm.(model.AlarmMessage)
		if !ok {
			logrus.Error("alarm is null")
		}
		Public(tmpA, Router_alarm)
		dataGateway.LoaderAlarm(tmpA)
	}
	elapsed := time.Since(start)
	logrus.Info("ExecCalc执行完成耗时：", elapsed)
}

func HasChange(device model.Device, old model.PropertyMessage, cur model.PropertyMessage) bool {
	if len(device.Product.Items) == 0 {
		return true
	}
	hasChange := false
	for _, item := range device.Product.Items {
		hasChange = hasChange || !utils.PropCompareEQ(old.Properties[item.Key], cur.Properties[item.Key], item.DataType, item.StepSize)
		if hasChange {
			break
		}
	}
	return hasChange
}
