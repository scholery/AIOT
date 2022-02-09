package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"koudai-box/iot/gateway/driver"
	"koudai-box/iot/gateway/model"
	status "koudai-box/iot/gateway/status"
	"koudai-box/iot/gateway/utils"

	"github.com/sirupsen/logrus"
)

func startExecEventPull(driver driver.Driver, device model.Device, dataCombination string, period int, cron string) {
	if dataCombination != model.DataCombination_Single && dataCombination != model.DataCombination_Array {
		return
	}

	if dataCombination == model.DataCombination_Single {
		//单条
		logrus.Debug("事件单条执行")
		if period <= 0 {
			c := status.GetDeviceCron(device.Id)
			err := c.AddFunc(cron, func() {
				ExecDeviceEventPull(driver, &device, period)
			})
			if err != nil {
				logrus.Errorf("cron[%s] of gateway %s is error", cron, driver.GetGatewayConfig().Key)
			}
			c.Start()
		} else {
			go ExecDeviceEventPull(driver, &device, period)
		}

	} else if dataCombination == model.DataCombination_Array {
		runLock.Lock()
		defer runLock.Unlock()
		ok := status.IsGatewayRunning(driver.GetGatewayConfig().Id)
		gateway := driver.GetGatewayConfig()
		status.StartGateway(gateway, &device)
		if !ok {
			//多条
			logrus.Debugf("启动事件批量执行，gateway[%s],devcie[%s]", driver.GetGatewayConfig().Key, device.Key)
			if period <= 0 {
				c := status.GetGatewayCron(gateway.Id)
				err := c.AddFunc(cron, func() {
					ExecDeviceEventPull(driver, nil, period)
				})
				if err != nil {
					logrus.Errorf("cron[%s] of gateway %s is error", cron, driver.GetGatewayConfig().Key)
				}
				c.Start()
			} else {
				go ExecDeviceEventPull(driver, nil, period)
			}
		} else {
			logrus.Debugf("事件批量执行，已经在运行，gateway[%s],devcie[%s]", driver.GetGatewayConfig().Key, device.Key)
		}
	}
}

/**
 *执行设备连接并抽取事件，单条
 */
func ExecDeviceEventPush(gatewayKey string, deviceKey string, data interface{}) error {
	gateway := GetGatewayConfigByKey(gatewayKey)
	if gateway == nil {
		logrus.Errorf("ExecDeviceEventPush gateway[%s] is not exist.", gatewayKey)
		return fmt.Errorf("gateway[%s] is not exist", gatewayKey)
	}
	var device *model.Device
	if len(gatewayKey) != 0 {
		devices := status.GetGatewayDevices(gateway.Id)
		for _, d := range devices {
			if d.Key == deviceKey {
				device = d
				break
			}
		}
	}
	driver, ok := status.GetDriver(gateway.Id)
	if !ok {
		logrus.Errorf("ExecDeviceEventPush gateway[%s]'s driver is not exist.", gatewayKey)
		return fmt.Errorf("gateway[%s]'s driver not exist", gatewayKey)
	}
	return doExecDeviceEvent(driver, device, data)
}

/**
 *执行设备连接并抽取事件，单条
 */
func ExecDeviceEventPull(driver driver.Driver, device *model.Device, period int) {
	gateway := driver.GetGatewayConfig()
	//判断是否停止
	if nil != device && !status.IsDeviceRunning(device.Id) {
		logrus.Infof("device %s is Stopped", device.Key)
		return
	} else {
		logrus.Debugf("ExecDeviceEventPull device %s", device.Key)
	}
	//下一次轮询
	if period > 0 {
		time.AfterFunc(time.Duration(period)*time.Second, func() { ExecDeviceEventPull(driver, device, period) })
		logrus.Debugf("ExecDeviceEventPull 预约%d秒后执行下一次抽取", period)
	}
	start := time.Now() // 获取当前时间
	//连接网络抽取数据
	deviceId := -1
	if nil != device {
		deviceId = device.Id
	}
	ts := status.GetPreApiTS(gateway.Id, model.API_GetEvent, deviceId)
	data, err := driver.FetchEvent(device, ts)
	if err != nil {
		logrus.Errorf("ExecDeviceEventPull FetchData error,gateway is %s,err:%s ", gateway.Key, err)
		//记录接口访问时间
		status.RecordApiRecord(gateway.Id, model.API_GetEvent, deviceId, start.Unix(), model.STATUS_ERROR)
		return
	}
	//记录接口访问时间
	status.RecordApiRecord(gateway.Id, model.API_GetEvent, deviceId, start.Unix(), model.STATUS_SUCCESS)

	doExecDeviceEvent(driver, device, data)

	elapsed := time.Since(start)
	logrus.Debugf("ExecPullDeviceProp[%s] 抽取事件执行完成耗时：%+v", device.Key, elapsed)
}

/**
 *执行设备连接并抽取事件，单条
 */
func doExecDeviceEvent(driver driver.Driver, device *model.Device, data interface{}) error {
	//单条
	if device != nil {
		tmp, err := driver.ExtracterEvent(data, device.Product)
		if err != nil {
			logrus.Errorf("doExecDeviceEvent devicde[%s] error.%+v", device.Key, err)
			return fmt.Errorf("ExtracterEvent devicde[%s]'s data err.%+v", device.Key, err)
		}
		var event model.EventMessage
		byte, err := json.Marshal(tmp)
		if err != nil {
			return fmt.Errorf("ExtracterEvent devicde[%s]'s data err.%+v", device.Key, err)
		}
		err = json.Unmarshal(byte, &event)
		if err != nil {
			logrus.Errorf("ExtracterEvent devicde[%s]'s data err.%+v", device.Key, err)
		}
		logrus.Debug("doExecDeviceEvent event:", event)
		if len(event.DeviceSign) != 0 {
			//处理抽取的数据
			doOneEventTreat(driver, device, event)
		}
		var events []model.EventMessage
		err = json.Unmarshal(byte, &events)
		logrus.Debug("doExecDeviceEvent events:", events)
		if nil == err {
			//预处理返回数据
			for _, val := range events {
				//处理抽取的数据
				doOneEventTreat(driver, device, val)
			}
		}

	} else { //多条
		devices := status.GetGatewayDevices(driver.GetGatewayConfig().Id)
		products := make(map[int]*model.Product)
		proDevs := make(map[int]map[string]*model.Device)
		for _, d := range devices {
			products[d.Product.Id] = d.Product
			des, ok := proDevs[d.Product.Id]
			if !ok {
				des = make(map[string]*model.Device)
			}
			des[d.SourceId] = d
			proDevs[d.Product.Id] = des
		}
		for _, p := range products {
			datas, err := driver.ExtracterEvent(data, p)
			if err != nil {
				logrus.Errorf("ExtracterEvent product[%s]'s data err.%+v", p.Key, err)
				continue
			}
			byte, err := json.Marshal(datas)
			if err != nil {
				logrus.Errorf("ExtracterEvent devicde[%s]'s data err.%+v", device.Key, err)
				continue
			}
			var events []model.EventMessage
			err = json.Unmarshal(byte, &events)
			logrus.Debug("doExecDeviceEvent events:", events)
			if nil != err && len(events) == 0 {
				logrus.Errorf("doExecDeviceEvent ExtracterEvent: product[%s] data error,not array.", p.Key)
				continue
			}
			//预处理返回数据
			for key, val := range events {
				d, ok := proDevs[p.Id][val.DeviceSign]
				if !ok {
					logrus.Errorf("doExecDeviceEvent ExtracterProp:device[%s] not exist", val.DeviceSign)
					continue
				}
				logrus.Debug("doExecDeviceEvent 处理数据:", key, val)
				//处理抽取的数据
				doOneEventTreat(driver, d, val)
			}
		}
	}
	return nil
}

func doOneEventTreat(driver driver.Driver, device *model.Device, event model.EventMessage) {
	event.DeviceId = strconv.Itoa(device.Id)
	event.MessageId = utils.GetUUID()
	event.Timestamp = time.Now().Unix()
	PostDeviceEventPull(event, driver, device)
}

/**
 *单设备-处理抽取的设备属性数据
 */
func PostDeviceEventPull(event interface{}, driver driver.Driver, device *model.Device) {
	//根据物模型数据转换
	data := event
	tmp, ok := data.(model.EventMessage)
	//发送数据获取成功通知
	if ok {
		model.EventMessChan <- model.EventChan{EventMessage: tmp, Device: device}
	}
}

/**
 *执行事件计算
 */
func ExecDeviceEventCalc(data model.EventMessage, device model.Device) {
	start := time.Now() // 获取当前时间
	dataGateway := &DataGateway{Device: &device}
	logrus.Debugf("EventMessage：%+v", data)
	//数据存储
	dataGateway.LoadeEvent(data)
	//属性推送
	Push(data, model.Message_Type_Event)
	elapsed := time.Since(start)
	logrus.Debugf("ExecDeviceEventCalc[%s]事件计算执行完成耗时：%+v", device.Key, elapsed)
}
