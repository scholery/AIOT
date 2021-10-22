package service

import (
	"main/driver"
	"main/model"
	"main/utils"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type PropertyChan struct {
	PropertyMessage model.PropertyMessage
	Device          model.Device
}

type EventChan struct {
	EventMessage model.EventMessage
	Device       model.Device
}

//运行标记
var run bool
var runLock sync.Mutex

//设备运行状态控制
var deviceThreads map[string]bool = make(map[string]bool)
var deviceRunLock sync.RWMutex

//设备状态
var deviceProps map[string]model.PropertyMessage = make(map[string]model.PropertyMessage)

//设备属性消息通道
var propMessChan = make(chan PropertyChan, Message_queen_size)
var eventMessChan = make(chan EventChan, Message_queen_size)

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
			go ExecDevicePropCalc(tmp.PropertyMessage, tmp.Device)
		case evt := <-eventMessChan:
			logrus.Info(evt)
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
	//拉取设备属性
	switch gateway.Protocol {
	case model.Geteway_Protocol_HTTP:
		getPropApi, ok := gateway.ApiConfigs[model.API_GetProp]
		if gateway.Protocol == "" && !ok {
			logrus.Error("GetGateway %s's API %s is null", gateway.Key, model.API_GetProp)
			return
		}
		if getPropApi.DataCombination == model.DataCombination_Single {
			//拉取属性信息
			for _, device := range devices {
				StartDevicePull(gateway, device)
			}
		} else {
			StartDevicePullBatch(gateway, devices)
		}
	case model.Geteway_Protocol_MQTT:
	case model.Geteway_Protocol_ModbusTCP:
	case model.Geteway_Protocol_OPCUA:
		//拉取属性信息
		for _, device := range devices {
			StartDevicePull(gateway, device)
		}
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
	switch gateway.Protocol {
	case model.Geteway_Protocol_HTTP:
		getPropApi, ok := gateway.ApiConfigs[model.API_GetProp]
		if gateway.Protocol == "" && !ok {
			logrus.Error("GetGateway %s's API %s is null", gateway.Key, model.API_GetProp)
			return false
		}
		if getPropApi.DataCombination == model.DataCombination_Single {
			for _, device := range devices {
				StopDevicePull(gateway, device)
			}
		} else {
			StartDevicePullBatch(gateway, devices)
		}
	case model.Geteway_Protocol_MQTT:
	case model.Geteway_Protocol_ModbusTCP:
	case model.Geteway_Protocol_OPCUA:
		for _, device := range devices {
			StopDevicePull(gateway, device)
		}
	}
	return true
}

/**
 *开启指定设备属性获取及事件监听
 */
func StartDevicePullBatch(gateway model.GatewayConfig, device []model.Device) {
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
	go ExecDevicePropPull(driver, device)
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
func ExecDevicePropPull(driver driver.Driver, device model.Device) {
	//判断是否停止
	if !isDeviceRun(device.Key) {
		logrus.Infof("Stop PullDeviceProp device %s", device.Key)
		return
	}
	logrus.Infof("PullDeviceProp device %s", device.Key)
	//下一次轮询
	period := driver.GetCollectPeriod(model.API_GetProp)
	time.AfterFunc(time.Duration(period)*time.Second, func() { ExecDevicePropPull(driver, device) })
	logrus.Infof("预约%d秒后执行下一次抽取", period)
	start := time.Now() // 获取当前时间
	//连接网络抽取数据
	data, err := driver.FetchProp(device)
	if err != nil {
		logrus.Errorf("FetchData error,device is %s,err:%s ", device.Key, err)
		return
	}
	//处理抽取的数据
	PostDevicePropPull(data, driver, device)
	elapsed := time.Since(start)
	logrus.Info("ExecPullDeviceProp 抽取数据执行完成耗时：", elapsed)
}

/**
 *处理抽取的设备属性数据
 */
func PostDevicePropPull(props interface{}, driver driver.Driver, device model.Device) {
	//预处理返回数据
	data, err := driver.ExtracterProp(props, device)
	if err != nil {
		data = props
		logrus.Errorf("Extracter data error,device is %s,data:\r\n%s ", device.Key, data)
		logrus.Error(err)
	}
	//根据物模型数据转换
	data, err = driver.TransformerProp(data, device)
	if err != nil {
		logrus.Errorf("Transformer data error,device is %s,data:\r\n%s ", device.Key, data)
		logrus.Error(err)
		return
	}
	tmp, ok := data.(model.PropertyMessage)
	//发送数据获取成功通知
	if ok {
		propMessChan <- PropertyChan{tmp, device}
	}
}

/**
 *执行计算
 */
func ExecDevicePropCalc(data model.PropertyMessage, device model.Device) {
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
	dataGateway.LoaderProperty(tmpP)
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
