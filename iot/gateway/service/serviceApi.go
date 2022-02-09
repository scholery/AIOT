package service

import (
	"koudai-box/iot/gateway/driver"
	"koudai-box/iot/gateway/model"
	"koudai-box/iot/service"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"

	status "koudai-box/iot/gateway/status"
)

var dataAsyncSaveCron cron.Cron

const dataAsyncSaveCronStr = "0 0/10 * * * *"

func InitIot() {
	// Connect()
	// defer Close()
	go StartPull()
}

/**
 *开启所有属性获取及事件监听
 */
func StartPull() {
	logrus.Info("***************************开启IOT数据监听***************************")
	gatewayConfigs, ok := GetGatewayConfigs()
	if !ok {
		logrus.Error("GetGatewayConfigs is null")
		return
	}
	status.Start()
	for _, gateway := range gatewayConfigs {
		StartGatewayPull(gateway)
	}
	ASyncStatusDatas()
	go MessageListener()
}

func MessageListener() {
	for status.IsRunning() {
		select {
		case tmp := <-model.PropMessChan: //解析后的设备属性后处理
			model.StatusMsgChan <- model.StatusMsg{
				DeviceId: tmp.Device.Id,
				Status:   model.STATUS_ACTIVE,
			}
			go ExecDevicePropCalc(tmp.PropertyMessage, *tmp.Device)
		case evt := <-model.EventMessChan: //解析后的事件后处理
			go ExecDeviceEventCalc(evt.EventMessage, *evt.Device)
		case msg := <-model.PushMsgChan: //推送消息解析
			if len(msg.GatewayKey) == 0 {
				logrus.Errorf("gateway is nil,msg:%s", msg.Msg)
				continue
			}
			if msg.Msg == nil {
				logrus.Errorf("ws:gateway[%s]'s msg is null,msg:%s", msg.GatewayKey, msg.Msg)
				continue
			}
			if len(msg.DeviceKey) == 0 {
				if len(msg.Type) == 0 || msg.Type == model.Msg_Type_Props {
					go PushGatewayDeviceProps(msg.GatewayKey, msg.Msg)
				}
				if len(msg.Type) == 0 || msg.Type == model.Msg_Type_Events {
					go PushGatewayEvents(msg.GatewayKey, msg.Msg)
				}
			} else {
				if len(msg.Type) == 0 || msg.Type == model.Msg_Type_Props {
					go PushDeviceProps(msg.GatewayKey, msg.DeviceKey, msg.Msg)
				}
				if len(msg.Type) == 0 || msg.Type == model.Msg_Type_Events {
					go PushDeviceEvents(msg.GatewayKey, msg.DeviceKey, msg.Msg)
				}
			}
		default:
		}
	}
}

func ASyncStatusDatas() {
	//加载api记录
	status.ReloadCacheData()
	dataAsyncSaveCron = *cron.New()
	err := dataAsyncSaveCron.AddFunc(dataAsyncSaveCronStr, func() {
		status.CacheDataAsyncSave()
	})
	if err != nil {
		logrus.Errorf("dataAsyncSaveCron cron[%s] is error", dataAsyncSaveCronStr)
	}
	dataAsyncSaveCron.Start()
}

/**
 *开启指定网关属性获取及事件监听
 */
func StartGatewayPull(gateway model.GatewayConfig) {
	logrus.Infof("StartGatewayPull GetGateway[%s]", gateway.Key)
	driver, ok := driver.GetDriver(&gateway)
	if !ok {
		logrus.Errorf("gateway[%s]'s driver is error.protocol is '%s'", gateway.Key, gateway.Protocol)
		return
	}
	err := driver.Start()
	if nil != err {
		logrus.Errorf("gateway[%s]'s driver start error.protocol is '%s'", gateway.Key, gateway.Protocol)
		return
	}
	status.PutDriver(gateway.Id, driver)
	devices, ok := GetDevices(gateway.Id, model.STATUS_ACTIVE)
	if !ok {
		logrus.Errorf("GetGateway[%s]'s related device is null", gateway.Key)
		return
	}
	//拉取属性信息
	for _, device := range devices {
		StartDevicePull(driver, *device)
	}

}

/**
 *停止指定网关属性获取及事件监听
 */
func StopGatewayPull(gateway model.GatewayConfig) bool {
	logrus.Infof("StopGatewayPull GetGateway %s", gateway.Key)
	devices, ok := GetDevices(gateway.Id, model.STATUS_ALL)
	if !ok {
		logrus.Error("GetGateway[%s]'s related device is null", gateway.Key)
		return false
	}
	for _, device := range devices {
		StopDevicePull(*device)
	}
	return true
}

// func stopTest(device model.Device) {
// 	time.Sleep(30 * time.Second)
// 	StopDevicePull(device)
// }

/**
 *开启指定设备属性获取及事件监听
 */
func StartDevicePull(driver driver.Driver, device model.Device) {
	//拉取设备属性
	logrus.Infof("StartDevicePull device[%s]", device.Key)
	gateway := driver.GetGatewayConfig()

	switch gateway.Protocol {
	case model.Geteway_Protocol_HTTP_Client:
		getPropApi, ok := gateway.ApiConfigs[model.API_GetProp]
		logrus.Debugf("getPropApi:%+v", getPropApi)
		if !ok {
			logrus.Errorf("GetGateway[%s]'s API %s is null", gateway.Key, model.API_GetProp)
		} else {
			if getPropApi.CollectType == model.CollectType_Poll {
				startExecDevicePropPull(driver, device, getPropApi.DataCombination, getPropApi.CollectPeriod, "")
			} else if getPropApi.CollectType == model.CollectType_Schedule {
				startExecDevicePropPull(driver, device, getPropApi.DataCombination, -1, getPropApi.Cron)
			}

		}
		getEventApi, ok := gateway.ApiConfigs[model.API_GetEvent]
		logrus.Debugf("getEventApi:%+v", getEventApi)
		if !ok {
			logrus.Errorf("GetGateway[%s]'s API %s is null", gateway.Key, model.API_GetEvent)
		} else {
			if getEventApi.CollectType == model.CollectType_Poll {
				startExecEventPull(driver, device, getEventApi.DataCombination, getEventApi.CollectPeriod, "")
			} else if getEventApi.CollectType == model.CollectType_Schedule {
				startExecEventPull(driver, device, getEventApi.DataCombination, -1, getEventApi.Cron)
			}

		}
		//测试停止
		//go stopTest(device)
	case model.Geteway_Protocol_HTTP_Server:
	case model.Geteway_Protocol_MQTT:
	case model.Geteway_Protocol_MQTTSN:
	case model.Geteway_Protocol_ModbusTCP, model.Geteway_Protocol_ModbusRTU, model.Geteway_Protocol_OPCUA:
		if gateway.CollectType == model.CollectType_Poll {
			startExecDevicePropPull(driver, device, model.DataCombination_Single, gateway.CollectPeriod, "")
		} else if gateway.CollectType == model.CollectType_Schedule {
			startExecDevicePropPull(driver, device, model.DataCombination_Single, -1, gateway.Cron)
		}
	case model.Geteway_Protocol_WebSocket_Client:
	case model.Geteway_Protocol_WebSocket_Server:
	case model.Geteway_Protocol_CoAP:
	case model.Geteway_Protocol_LwM2M:
	case model.Geteway_Protocol_BACnet_IP:
	}

	status.StartDevice(&device)
	status.StartGateway(gateway, &device)
}

/**
 *停止指定设备属性获取及事件监听
 */
func StopDevicePull(device model.Device) bool {
	logrus.Infof("StopDevicePull device %s", device.Key)
	status.StopDevice(&device)
	return true
}

/*************************************对外接口 开始**********************************************/
/**
 *停止所有属性获取及事件监听
 */
func StopPull() {
	status.Stop()
	dataAsyncSaveCron.Stop()
	logrus.Info("***************************停止IOT数据监听***************************")
}

/**
 *启动网关
 */
func StartGateway(getewayId int) bool {
	gateway := GetGatewayConfig(getewayId)
	if gateway == nil {
		logrus.Errorf("StartGateway error,gateway is null,id[%d]", getewayId)
		return false
	}
	StartGatewayPull(*gateway)
	return true
}
func StopGateway(getewayId int) bool {
	gateway := GetGatewayConfig(getewayId)
	if gateway == nil {
		logrus.Errorf("StopGateway error,gateway is null,id[%d]", getewayId)
		return false
	}
	StopGatewayPull(*gateway)
	return true
}
func StartProduct(productId int) bool {
	product, ok := GetProduct(productId, model.STATUS_ALL)
	if !ok {
		logrus.Errorf("StartProduct error,product is null or disactive,productId[%d]", productId)
		return false
	}

	var ids []int
	ids = append(ids, productId)
	devices, ok := GetDeviceByProductIds(ids, model.STATUS_ACTIVE)
	if !ok {
		logrus.Errorf("StartProduct error,devices is null,productId[%d]", productId)
		return false
	}

	dri, ok := status.GetDriver(product.GatewayId)
	if !ok {
		gateway := GetGatewayConfig(product.GatewayId)
		if gateway == nil {
			logrus.Errorf("StartProduct error,gateway is null,productId[%d],GatewayId[%d]", productId, product.GatewayId)
			return false
		}
		dri, ok = driver.GetDriver(gateway)
		if !ok {
			logrus.Errorf("StartProduct error,gateway %s 's driver is error.protocol is '%s'", gateway.Key, gateway.Protocol)
			return false
		}
		err := dri.Start()
		if nil != err {
			logrus.Errorf("gateway[%s]'s driver start error.protocol is '%s'", gateway.Key, gateway.Protocol)
			return false
		}
		status.PutDriver(gateway.Id, dri)
	}
	for _, d := range devices {
		StartDevicePull(dri, *d)
		// service.SetDeviceRunningStatus(d.Id, model.STATUS_ACTIVE)
	}
	return true
}
func StopProduct(productId int) bool {
	var ids []int
	ids = append(ids, productId)
	devices, ok := GetDeviceByProductIds(ids, model.STATUS_ALL)
	if !ok {
		logrus.Errorf("StopProduct error,devices is null,productId[%d]", productId)
		return false
	}
	for _, d := range devices {
		ok := StopDevicePull(*d)
		if ok {
			service.SetDeviceRunningStatus(d.Id, model.STATUS_DISACTIVE)
		}
	}
	return true
}
func StartDevice(deviceId int) bool {
	d, ok := GetDeviceById(deviceId, model.STATUS_ALL)
	if !ok {
		logrus.Errorf("StartDevice error,device is null,deviceId[%d]", deviceId)
		return false
	}
	logrus.Debugf("StartDevice[%d],%s", deviceId, d.Key)
	dri, ok := status.GetDriver(d.Product.GatewayId)
	if !ok {
		gateway := GetGatewayConfig(d.Product.GatewayId)
		if gateway == nil {
			logrus.Errorf("StartDevice error,gateway is null,productId[%d],GatewayId[%d]", d.Product.Id, d.Product.GatewayId)
			return false
		}
		dri, ok = driver.GetDriver(gateway)
		if !ok {
			logrus.Errorf("StartDevice error,gateway %s 's driver is error.protocol is '%s'", gateway.Key, gateway.Protocol)
			return false
		}
		err := dri.Start()
		if nil != err {
			logrus.Errorf("gateway[%s]'s driver start error.protocol is '%s'", gateway.Key, gateway.Protocol)
			return false
		}
		status.PutDriver(gateway.Id, dri)
	}
	StartDevicePull(dri, *d)
	return true
}
func StopDevice(deviceId int) bool {
	d, ok := GetDeviceById(deviceId, model.STATUS_ALL)
	logrus.Debugf("StopDevice[%d],%s", deviceId, d.Key)
	if !ok {
		logrus.Errorf("StopDevice error,device is null,deviceId[%d]", deviceId)
		return false
	}
	StopDevicePull(*d)
	return true
}

func PushDeviceProps(gatewayKey string, deviceKey string, data interface{}) error {
	return ExecPropPush(gatewayKey, deviceKey, data)
}

func PushGatewayDeviceProps(gatewayKey string, data interface{}) error {
	return ExecPropPushBatch(gatewayKey, data)
}

func PushGatewayEvents(gatewayKey string, data interface{}) error {
	return ExecDeviceEventPush(gatewayKey, "", data)
}

func PushDeviceEvents(gatewayKey string, deviceKey string, data interface{}) error {
	return ExecDeviceEventPush(gatewayKey, deviceKey, data)
}

func SetZeroStatus(deviceId int) {
	dStatus := status.GetDeviceStatus(deviceId)
	dStatus.ZeroStatus = dStatus.LastStatus
}

func SetPredayStatus(deviceId int, preday *model.PropertyMessage) {
	dStatus := status.GetDeviceStatus(deviceId)
	dStatus.PreDayStatus = preday
}

/*************************************对外接口 结束**********************************************/
