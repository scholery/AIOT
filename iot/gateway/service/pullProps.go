package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"koudai-box/iot/gateway/driver"
	"koudai-box/iot/gateway/model"
	status "koudai-box/iot/gateway/status"
	"koudai-box/iot/gateway/utils"

	"github.com/sirupsen/logrus"
)

var runLock sync.Mutex

func startExecDevicePropPull(driver driver.Driver, device model.Device, dataCombination string, period int, cron string) {
	if dataCombination != model.DataCombination_Single && dataCombination != model.DataCombination_Array {
		return
	}

	if dataCombination == model.DataCombination_Single {
		//单条
		logrus.Debug("单条执行")
		if period <= 0 {
			c := status.GetDeviceCron(device.Id)
			err := c.AddFunc(cron, func() {
				ExecDevicePropPull(driver, device, period)
			})
			if err != nil {
				logrus.Errorf("cron[%s] of gateway %s is error", cron, driver.GetGatewayConfig().Key)
				return
			}
			c.Start()
		} else {
			go ExecDevicePropPull(driver, device, period)
		}

	} else if dataCombination == model.DataCombination_Array {
		runLock.Lock()
		defer runLock.Unlock()
		ok := status.IsGatewayRunning(driver.GetGatewayConfig().Id)
		gateway := driver.GetGatewayConfig()
		status.StartGateway(gateway, &device)
		if !ok {
			//多条
			logrus.Debugf("启动批量执行，gateway[%s],devcie[%s]", driver.GetGatewayConfig().Key, device.Key)
			if period <= 0 {
				c := status.GetGatewayCron(gateway.Id)
				err := c.AddFunc(cron, func() {
					ExecDevicePropPullBatch(driver, period)
				})
				if err != nil {
					logrus.Errorf("cron[%s] of gateway %s is error", cron, driver.GetGatewayConfig().Key)
					return
				}
				c.Start()
			} else {
				go ExecDevicePropPullBatch(driver, period)
			}
		} else {
			logrus.Debugf("批量执行，已经在运行，gateway[%s],devcie[%s]", driver.GetGatewayConfig().Key, device.Key)
		}
	}
}

/**
 *执行设备连接并抽取数据，批量:http client
 */
func ExecDevicePropPullBatch(driver driver.Driver, period int) {
	gateway := driver.GetGatewayConfig()
	logrus.Debugf("ExecDevicePropPullBatch gateway: %s", gateway.Key)
	//下一次轮询
	if period > 0 {
		timer := time.AfterFunc(time.Duration(period)*time.Second, func() { ExecDevicePropPullBatch(driver, period) })
		status.SetGatewayTimer(gateway.Id, timer)
		logrus.Debugf("ExecDevicePropPullBatch 预约%d秒后执行下一次批量抽取", period)
	}
	start := time.Now() // 获取当前时间
	//连接网络抽取数据
	ts := status.GetPreApiTS(gateway.Id, model.API_GetProp, -1)
	data, err := driver.FetchPropBatch(ts)
	if err != nil {
		logrus.Errorf("ExecDevicePropPullBatch FetchPropBatch error,gateway is %s,err:%s ", gateway.Key, err)
		//设备离线
		devices := status.GetGatewayDevices(gateway.Id)
		for id, _ := range devices {
			model.StatusMsgChan <- model.StatusMsg{
				DeviceId: id,
				Status:   model.STATUS_DISACTIVE,
			}
		}
		//记录接口访问时间
		status.RecordApiRecord(gateway.Id, model.API_GetProp, -1, start.Unix(), model.STATUS_ERROR)
		return
	}
	//记录接口访问时间
	status.RecordApiRecord(gateway.Id, model.API_GetProp, -1, start.Unix(), model.STATUS_SUCCESS)

	doExecDevicePropBatch(driver, gateway, data)

	elapsed := time.Since(start)
	logrus.Debugf("ExecDevicePropPullBatch[%s]抽取数据执行完成耗时：%+v", gateway.Key, elapsed)
}

/**
 *执行设备连接并抽取数据，单条
 */
func ExecDevicePropPull(driver driver.Driver, device model.Device, period int) {
	//判断是否停止
	if !status.IsDeviceRunning(device.Id) {
		logrus.Infof("device %s is Stopped", device.Key)
		return
	}
	logrus.Debugf("ExecDevicePropPull device %s", device.Key)
	//下一次轮询
	if period > 0 {
		timer := time.AfterFunc(time.Duration(period)*time.Second, func() { ExecDevicePropPull(driver, device, period) })
		status.SetDeviceTimer(device.Id, timer)
		logrus.Debugf("ExecDevicePropPull 预约%d秒后执行下一次抽取", period)
	}
	start := time.Now() // 获取当前时间
	//连接网络抽取数据
	ts := status.GetPreApiTS(device.Product.GatewayId, model.API_GetProp, device.Id)
	data, err := driver.FetchProp(&device, ts)
	if err != nil || data == nil {
		logrus.Errorf("ExecDevicePropPull FetchData error,device is %s,err:%s ", device.Key, err)
		model.StatusMsgChan <- model.StatusMsg{
			DeviceId: device.Id,
			Status:   model.STATUS_DISACTIVE,
		}
		//记录接口访问时间
		status.RecordApiRecord(device.Product.GatewayId, model.API_GetProp, device.Id, start.Unix(), model.STATUS_ERROR)
		return
	}
	//记录接口访问时间
	status.RecordApiRecord(device.Product.GatewayId, model.API_GetProp, device.Id, start.Unix(), model.STATUS_SUCCESS)

	doExecDeviceProp(driver, &device, data)

	elapsed := time.Since(start)
	logrus.Debugf("ExecPullDeviceProp[%s] 抽取数据执行完成耗时：%+v", device.Key, elapsed)
}

func ExecPropPush(gatewayKey string, deviceKey string, data interface{}) error {
	gateway := GetGatewayConfigByKey(gatewayKey)
	if gateway == nil {
		logrus.Errorf("ExecPropPush gateway[%s] is not exist.", gatewayKey)
		return fmt.Errorf("gateway[%s] is not exist", gatewayKey)
	}
	gatewayId := gateway.Id
	logrus.Debugf("ExecPropPush device[%s]", deviceKey)
	devices := status.GetGatewayDevices(gatewayId)
	var device *model.Device
	for _, d := range devices {
		if d.Key == deviceKey {
			device = d
			break
		}
	}

	if device == nil {
		logrus.Errorf("ExecPropPush error,devicde[%s] is not exist.", deviceKey)
		return errors.New("devicde is not exist")
	}
	driver, ok := status.GetDriver(gatewayId)
	if !ok {
		logrus.Errorf("ExecPropPush gateway[%s]'s driver is not exist.devicde[%s] ", gatewayId, deviceKey)
		return errors.New("gateway's driver not exist")
	}
	return doExecDeviceProp(driver, device, data)
}

/**
 *单设备-处理推送的设备属性数据
 */
func doExecDeviceProp(driver driver.Driver, device *model.Device, data interface{}) error {
	props, err := driver.ExtracterProp(data, device.Product)
	logrus.Debugf("doExecDeviceProp ExtracterProp,device[%s],data:%+v,props:%+v", device.Key, data, props)
	if err != nil {
		props = data
		logrus.Errorf("Extracter data error,device is %s,data:%+v,err:%+v ", device.Key, data, err)
	}
	//处理抽取的数据
	PostDevicePropPull(props, driver, *device)
	return nil
}
func ExecPropPushBatch(gatewayKey string, data interface{}) error {
	gateway := GetGatewayConfigByKey(gatewayKey)
	if gateway == nil {
		logrus.Errorf("ExecPropPushBatch gateway[%s] is not exist.", gatewayKey)
		return fmt.Errorf("gateway[%s] is not exist", gatewayKey)
	}
	driver, ok := status.GetDriver(gateway.Id)
	if !ok {
		logrus.Errorf("ExecPropPushBatch gateway[%s]'s driver is not exist.", gatewayKey)
		return fmt.Errorf("gateway[%s]'s driver not exist", gatewayKey)
	}
	return doExecDevicePropBatch(driver, gateway, data)
}

/**
 *多设备-处理推送的设备属性数据，批量
 */
func doExecDevicePropBatch(driver driver.Driver, gateway *model.GatewayConfig, data interface{}) error {
	gatewayId := gateway.Id
	devices := status.GetGatewayDevices(gatewayId)
	if len(devices) == 0 {
		logrus.Errorf("doExecDevicePropBatch gateway[%d]'s devicde is null.", gatewayId)
		return errors.New("devicde is not exist")
	}
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
		datas, err := driver.ExtracterProp(data, p)
		if err != nil {
			logrus.Errorf("doExecDevicePropBatch ExtracterProp product[%s] error[%v]", p.Key, err)
			continue
		}
		dataMap, ok := datas.(map[string]interface{})
		logrus.Debug("doExecDevicePropBatch dataMap:", dataMap)
		if !ok {
			logrus.Errorf("doExecDevicePropBatch ExtracterProp:transformer product[%s] data error", p.Key)
			continue
		}
		devicesMap, ok := proDevs[p.Id]
		if !ok {
			logrus.Debug("doExecDevicePropBatch 产品[%s:%s]没有设备在运行", p.Key, p.Name)
			continue
		}
		//预处理返回数据
		for key, dev := range devicesMap {
			data, ok := dataMap[key]
			if !ok {
				logrus.Errorf("doExecDevicePropBatch 设备[%s:%s]没有数据", key, dev.Name)
				continue
			}
			logrus.Debug("doExecDevicePropBatch 处理数据:", key, data)
			//处理抽取的数据
			PostDevicePropPull(data, driver, *dev)
		}
		// for key, val := range dataMap {
		// 	d, ok := proDevs[p.Id][key]
		// 	logrus.Debug("key=", key)
		// 	if !ok {
		// 		logrus.Errorf("doExecDevicePropBatch ExtracterProp:device key[%s] not exist", key)
		// 		continue
		// 	}
		// 	logrus.Debug("doExecDevicePropBatch 处理数据:", key, val)
		// 	//处理抽取的数据
		// 	PostDevicePropPull(val, driver, *d)
		// }
	}
	return nil
}

/**
 *单设备-处理抽取的设备属性数据
 */
func PostDevicePropPull(props interface{}, driver driver.Driver, device model.Device) {
	//根据物模型数据转换
	data := props
	data, err := driver.TransformerProp(data, &device)
	if err != nil {
		logrus.Errorf("Transformer data error,device is %s,data:\r\n%s ", device.Key, data)
		logrus.Error(err)
		return
	}
	tmp, ok := data.(model.PropertyMessage)
	//发送数据获取成功通知
	if ok {
		model.PropMessChan <- model.PropertyChan{PropertyMessage: tmp, Device: &device}
	}
}

/**
 *执行计算
 */
func ExecDevicePropCalc(data model.PropertyMessage, device model.Device) {
	start := time.Now() // 获取当前时间
	dataGateway := &DataGateway{Device: &device}
	// logrus.Debugf("PropertyMessage：%+v", data)
	// 执行计算函数
	tmpP := data
	res, err := dataGateway.Calculater(data)
	if err != nil {
		logrus.Errorf("Calculater error,res:%+v,err:%+v", res, err)
		return
	} else {
		r, ok := res.(model.PropertyMessage)
		if !ok {
			byte, err := json.Marshal(res)
			if err != nil {
				logrus.Errorf("Calculater error,Marshal convert res error:%+v,err:%+v", res, err)
				return
			}
			var m model.PropertyMessage
			err = json.Unmarshal(byte, &m)
			if err != nil {
				logrus.Errorf("Calculater error,Unmarshal convert res error:%+v,err:%+v", res, err)
				return
			} else {
				logrus.Debugf("Calculater,Marshal convert res:%+v", res)
			}
			r = m
		}
		tmpP = r
	}
	//变化上报
	old, ok := status.GetDeviceLastProp(device.Id)
	if ok && !HasChange(device, old, tmpP) {
		logrus.Debugf("device[%s]'s prop no change", device.Key)
		return
	}
	//数据存储
	dataGateway.LoaderProperty(tmpP, true)
	//更新缓存的最新状态
	status.PutDeviceLastProp(device.Id, tmpP)
	//告警过滤
	alarms, err := dataGateway.FilterAlarm(tmpP)
	if err != nil {
		logrus.Error(err)
		return
	}
	//计算告警
	for _, alarm := range alarms {
		tmpA, ok := alarm.(model.IotEventMessage)
		if !ok {
			logrus.Error("alarm is null")
		}
		//告警存储
		dataGateway.LoaderAlarm(tmpA, true)
	}
	elapsed := time.Since(start)
	logrus.Debugf("ExecCalc[%s]数据计算执行完成耗时：%+v", device.Key, elapsed)
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
		hasChange = strings.Compare(model.DataReportType_Schedule, item.DataReportType) == 0 || !utils.PropCompareEQ(old.Properties[item.Code], cur.Properties[item.Code], item.DataType)
		if hasChange {
			break
		}
	}
	return hasChange
}
