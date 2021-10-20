package service

import (
	"errors"
	"fmt"
	"main/driver"
	"main/model"
	"main/utils"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const Message_queen_size = 10

//设备属性消息通道
var propMessChan = make(chan PropertyChan, Message_queen_size)

//运行标记
var run bool
var runLock sync.Mutex

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

type PropertyChan struct {
	PropertyMessage model.PropertyMessage
	Device          model.Device
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
func (gateway *DataGateway) LoaderMessage(data model.PropertyMessage) (bool, error) {
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

//设备运行状态控制
var deviceThreads map[string]bool = make(map[string]bool)
var deviceRunLock sync.RWMutex

//设备状态
var deviceProps map[string]model.PropertyMessage = make(map[string]model.PropertyMessage)

/**
 *开启所有属性获取及事件监听
 */
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
 */
func StopPull() {
	logrus.Info("StopPull")
	runLock.Lock()
	defer runLock.Unlock()
	run = false
	for k := range deviceThreads {
		setDeviceStop(k)
	}
}

/**
 *开启指定网关属性获取及事件监听
 */
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

/**
 *停止指定网关属性获取及事件监听
 */
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

/**
 *开启指定设备属性获取及事件监听
 */
func StartDevicePull(gateway model.GatewayConfig, device model.Device) {
	logrus.Infof("StartDevicePull device %s", device.Key)
	setDeviceRun(device.Key)
	driver, ok := driver.GetDriver(gateway, device)
	if !ok {
		logrus.Errorf("driver init failed,type is %s,device is %s", gateway.Protocol, device.Key)
		return
	}
	go ExecPull(driver, device, 0)
}

/**
 *停止指定设备属性获取及事件监听
 */
func StopDevicePull(gateway model.GatewayConfig, device model.Device) bool {
	logrus.Infof("StopDevicePull device %s", device.Key)
	setDeviceStop(device.Key)
	return true
}

func isDeviceRun(key string) bool {
	deviceRunLock.RLock()
	defer deviceRunLock.RUnlock()
	st, ok := deviceThreads[key]
	return ok && st
}

func setDeviceRun(key string) {
	deviceRunLock.Lock()
	defer deviceRunLock.Unlock()
	deviceThreads[key] = true
}

func setDeviceStop(key string) {
	deviceRunLock.Lock()
	defer deviceRunLock.Unlock()
	delete(deviceThreads, key)
}

/**
 *执行设备连接并抽取数据
 */
func ExecPull(driver driver.Driver, device model.Device, sleep int) {
	//判断是否停止
	if !isDeviceRun(device.Key) {
		logrus.Infof("Stop ExecPull device %s", device.Key)
		return
	}
	//采集间隔
	if sleep > 0 {
		time.Sleep(time.Duration(sleep) * time.Second)
		logrus.Infof("预约%d秒后执行下一次抽取", sleep)
	}
	//预约下一次执行
	go ExecPull(driver, device, device.Product.CollectPeriod)

	logrus.Infof("ExecPull device %s", device.Key)
	start := time.Now() // 获取当前时间
	//连接网络抽取数据
	data, err := driver.FetchData()
	if err != nil {
		logrus.Errorf("FetchData error,device is %s,err:%s ", device.Key, err)
		return
	}
	//预处理返回数据
	data, err = driver.Extracter(data)
	if err != nil {
		logrus.Errorf("Extracter data error,device is %s,data:\r\n%s ", device.Key, data)
		logrus.Error(err)
	}
	//根据物模型数据转换
	data, err = driver.Transformer(data)
	if err != nil {
		logrus.Errorf("Transformer data error,device is %s,data:\r\n%s ", device.Key, data)
		logrus.Error(err)
		return
	}
	elapsed := time.Since(start)
	logrus.Info("ExecPull 抽取数据执行完成耗时：", elapsed)
	tmp, ok := data.(model.PropertyMessage)
	//发送数据获取成功通知
	if ok {
		propMessChan <- PropertyChan{tmp, device}
	}
}

/**
 *执行计算
 */
func ExecCalc(data model.PropertyMessage, device model.Device) {
	start := time.Now() // 获取当前时间
	dataGateway := &DataGateway{Device: device}
	//执行计算函数
	res, err := dataGateway.Calculater(data)
	if err != nil {
		logrus.Error("Calculater error.", err)
	}
	tmpP, ok := res.(model.PropertyMessage)
	if !ok {
		logrus.Error("calc error")
		return
	}
	//变化上报
	old, ok := deviceProps[device.Key]
	if ok && !HasChange(device, old, tmpP) {
		logrus.Info("no change")
		return
	}
	//数据存储
	dataGateway.LoaderMessage(tmpP)
	//属性推送
	Public(tmpP, Router_prop)
	//更新缓存的最新状态
	deviceProps[device.Key] = tmpP
	//告警过滤
	alarms, err := dataGateway.FilterAlarm(tmpP)
	logrus.Info("alarms", alarms)
	if err != nil {
		logrus.Error(err)
		return
	}
	//计算告警
	for _, alarm := range alarms {
		tmpA, ok := alarm.(model.AlarmMessage)
		if !ok {
			logrus.Error("alarm is null")
		}
		//告警推送
		Public(tmpA, Router_alarm)
		//告警存储
		dataGateway.LoaderAlarm(tmpA)
	}
	elapsed := time.Since(start)
	logrus.Info("ExecCalc 数据计算执行完成耗时：", elapsed)
}

/**
 *属性是否变化
 */
func HasChange(device model.Device, old model.PropertyMessage, cur model.PropertyMessage) bool {
	if len(device.Product.Items) == 0 {
		return true
	}
	hasChange := false
	for _, item := range device.Product.Items {
		hasChange = strings.Compare(model.DataReportType_Schedule, item.DataReportType) == 0 || !utils.PropCompareEQ(old.Properties[item.Key], cur.Properties[item.Key], item.DataType)
		if hasChange {
			break
		}
	}
	return hasChange
}
