package gateway

import (
	"encoding/json"
	"sync"
	"time"

	"koudai-box/iot/db"
	"koudai-box/iot/gateway/driver"
	"koudai-box/iot/gateway/model"
	"koudai-box/iot/gateway/utils"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

/****状态****/
//运行状态
var run bool = true
var runLock sync.Mutex

//设备运行状态控制
var deviceRunLock sync.RWMutex
var deviceStatusLock sync.RWMutex
var deviceTimerLock sync.Mutex
var gatewayRunLock sync.RWMutex
var deviceThreads map[int]*model.Device = make(map[int]*model.Device)
var deviceCrons map[int]*cron.Cron = make(map[int]*cron.Cron)
var deviceTimers map[int]*time.Timer = make(map[int]*time.Timer)
var gatewayThreads map[int]map[int]*model.Device = make(map[int]map[int]*model.Device)
var gatewayCrons map[int]*cron.Cron = make(map[int]*cron.Cron)
var gatewayTimers map[int]*time.Timer = make(map[int]*time.Timer)

var apilogLock sync.Mutex

type API_Status struct {
	Timestamp    int64
	Status       int
	PreTimestamp int64
}

//gatewayId：API_key:deviceID(-1 batch):times(-1:error)
var gatewayApiRecords map[int]map[string]map[int]*API_Status = make(map[int]map[string]map[int]*API_Status)

/****缓存****/
//网关
var cache_drivers map[int]driver.Driver = make(map[int]driver.Driver)

//设备最新消息
var cache_deviceStatus map[int]*model.DeviceStatus = make(map[int]*model.DeviceStatus)

func IsRunning() bool {
	return run
}

func IsDeviceRunning(id int) bool {
	if !run {
		return false
	}

	deviceRunLock.RLock()
	defer deviceRunLock.RUnlock()
	_, ok := deviceThreads[id]
	return ok
}

func IsGatewayRunning(gatewayId int) bool {
	if !run {
		return false
	}
	devices, ok := gatewayThreads[gatewayId]
	if !ok || len(devices) == 0 {
		return false
	}
	return true
}

func Start() {
	runLock.Lock()
	defer runLock.Unlock()
	run = true
}
func Stop() {
	runLock.Lock()
	defer runLock.Unlock()
	run = false
	for _, d := range deviceThreads {
		StopDevice(d)
	}
}

func StartGateway(gateway *model.GatewayConfig, device *model.Device) {
	gatewayRunLock.Lock()
	defer gatewayRunLock.Unlock()
	devices, ok := gatewayThreads[gateway.Id]
	if !ok {
		devices = make(map[int]*model.Device)
	}
	devices[device.Id] = device
	gatewayThreads[gateway.Id] = devices
}

func StopGateway(id int) {
	driver, ok := GetDriver(id)
	//放前面会死锁
	gatewayRunLock.Lock()
	defer gatewayRunLock.Unlock()
	if ok {
		driver.Stop()
		delete(cache_drivers, id)
	}
	cron, ok := gatewayCrons[id]
	if ok {
		logrus.Debugf("stop gateway[%d]'s cron", id)
		cron.Stop()
		delete(gatewayCrons, id)
	}
	timer, ok := gatewayTimers[id]
	if ok {
		logrus.Debugf("stop gateway[%d]'s timer", id)
		timer.Stop()
		delete(gatewayTimers, id)
	}
	delete(gatewayThreads, id)
}

func StartDevice(device *model.Device) {
	deviceRunLock.Lock()
	defer deviceRunLock.Unlock()
	deviceThreads[device.Id] = device
	//状态重置
	model.StatusMsgChan <- model.StatusMsg{
		DeviceId: device.Id,
		Status:   model.STATUS_UNKNOWN,
	}
}

func StopDevice(device *model.Device) {
	deviceRunLock.Lock()
	defer deviceRunLock.Unlock()
	delete(deviceThreads, device.Id)
	cron, ok := deviceCrons[device.Id]
	if ok {
		logrus.Debugf("stop device[%s]'s cron", device.Key)
		cron.Stop()
		delete(deviceCrons, device.Id)
	}
	deviceTimerLock.Lock()
	defer deviceTimerLock.Unlock()
	timer, ok := deviceTimers[device.Id]
	if ok {
		logrus.Debugf("stop device[%s]'s timer", device.Key)
		timer.Stop()
		delete(deviceTimers, device.Id)
	}
	devices, ok := gatewayThreads[device.Product.GatewayId]
	if ok {
		_, ok = devices[device.Id]
		if ok {
			delete(devices, device.Id)
		}
		gatewayThreads[device.Product.GatewayId] = devices
		if len(devices) == 0 {
			StopGateway(device.Product.GatewayId)
		}
	}
	//状态重置
	model.StatusMsgChan <- model.StatusMsg{
		DeviceId: device.Id,
		Status:   model.STATUS_UNKNOWN,
	}
}

func PutDeviceLastProp(deviceId int, prop model.PropertyMessage) {
	deviceStatusLock.Lock()
	defer deviceStatusLock.Unlock()
	status, ok := cache_deviceStatus[deviceId]
	if !ok {
		status = &model.DeviceStatus{
			Id:     deviceId,
			Key:    prop.DeviceSign,
			Status: model.STATUS_UNKNOWN,
		}
	}
	status.LastStatus = &prop
	if status.ZeroStatus == nil {
		zeroStatus := *status.LastStatus
		status.ZeroStatus = &zeroStatus
	}
	if status.PreDayStatus == nil {
		preDayStatus := *status.LastStatus
		status.PreDayStatus = &preDayStatus
	}
	cache_deviceStatus[deviceId] = status
}

func GetDeviceLastProp(deviceId int) (model.PropertyMessage, bool) {
	deviceStatusLock.RLock()
	defer deviceStatusLock.RUnlock()
	status, ok := cache_deviceStatus[deviceId]
	if !ok || status.LastStatus == nil {
		return model.PropertyMessage{}, ok
	}
	return *status.LastStatus, ok
}

func GetDeviceStatus(deviceId int) *model.DeviceStatus {
	deviceStatusLock.RLock()
	status, ok := cache_deviceStatus[deviceId]
	deviceStatusLock.RUnlock()
	if !ok {
		deviceStatusLock.Lock()
		defer deviceStatusLock.Unlock()
		status = &model.DeviceStatus{Id: deviceId, Status: model.STATUS_UNKNOWN, Timestamp: time.Now().Unix()}
		cache_deviceStatus[deviceId] = status
	}
	return status
}

func GetDevice(deviceId int) *model.Device {
	device, ok := deviceThreads[deviceId]
	if ok {
		return device
	}
	return nil
}

func PutDriver(gatewayId int, driver driver.Driver) {
	gatewayRunLock.Lock()
	defer gatewayRunLock.Unlock()
	cache_drivers[gatewayId] = driver
}

func GetDriver(gatewayId int) (driver.Driver, bool) {
	gatewayRunLock.RLock()
	defer gatewayRunLock.RUnlock()
	driver, ok := cache_drivers[gatewayId]
	return driver, ok
}

func GetDeviceCron(deviceId int) *cron.Cron {
	c, ok := deviceCrons[deviceId]
	if ok {
		return c
	}
	c = cron.New()
	deviceCrons[deviceId] = c
	return c
}

func SetDeviceTimer(deviceId int, timer *time.Timer) {
	deviceTimerLock.Lock()
	defer deviceTimerLock.Unlock()
	deviceTimers[deviceId] = timer
}

func GetGatewayCron(gatewayId int) *cron.Cron {
	c, ok := gatewayCrons[gatewayId]
	if ok {
		return c
	}
	c = cron.New()
	gatewayCrons[gatewayId] = c
	return c
}

func SetGatewayTimer(gatewayId int, timer *time.Timer) {
	gatewayTimers[gatewayId] = timer
}

func RecordApiRecord(gatewayId int, api string, deviceId int, timestamp int64, status int) {
	apilogLock.Lock()
	defer apilogLock.Unlock()
	records, ok := gatewayApiRecords[gatewayId]
	if !ok {
		records = make(map[string]map[int]*API_Status)
	}
	records1, ok := records[api]
	if !ok {
		records1 = make(map[int]*API_Status)
	}
	tmp := records1[deviceId]
	if nil == tmp {
		tmp = &API_Status{PreTimestamp: -1}
	}
	tmp.Status = status
	tmp.Timestamp = timestamp
	if status == model.STATUS_SUCCESS {
		tmp.PreTimestamp = tmp.Timestamp
	}
	records1[deviceId] = tmp
	records[api] = records1
	gatewayApiRecords[gatewayId] = records
}

func GetPreApiTS(gatewayId int, api string, deviceId int) int64 {
	records, ok := gatewayApiRecords[gatewayId]
	if !ok {
		return -1
	}
	records1, ok := records[api]
	if !ok {
		return -1
	}
	ts, ok := records1[deviceId]
	if !ok {
		return -1
	}
	if ts.Status == model.STATUS_ERROR {
		return ts.PreTimestamp
	}
	return ts.Timestamp
}

func GetGatewayDevices(id int) map[int]*model.Device {
	devices, ok := gatewayThreads[id]
	if !ok {
		devices = make(map[int]*model.Device)
	}
	return devices
}

func ReloadCacheData() {
	apilogLock.Lock()
	defer apilogLock.Unlock()
	//reload api data
	_, histories := db.QueryAllGatewayApiHistory()
	tmp := make(map[int]map[string]map[int]*API_Status)
	for _, val := range histories {
		records, ok := tmp[val.GatewayId]
		if !ok {
			records = make(map[string]map[int]*API_Status)
		}
		records1, ok := records[val.Name]
		if !ok {
			records1 = make(map[int]*API_Status)
		}
		api := &API_Status{
			Status:       val.Status,
			Timestamp:    val.UpdateTime.Local().Unix(),
			PreTimestamp: val.Timestamp,
		}
		if val.Status == model.STATUS_ERROR {
			api.Timestamp = val.UpdateTime.Local().Unix()
		}
		records1[val.DeviceId] = api
		records[val.Name] = records1
		tmp[val.GatewayId] = records
	}
	gatewayApiRecords = tmp
	//reload device status data
	_, items := db.QueryAllDeviceStatus()
	tmp1 := make(map[int]*model.DeviceStatus)
	for _, item := range items {
		status := &model.DeviceStatus{
			Id:        item.DeviceId,
			Status:    model.STATUS_UNKNOWN,
			Timestamp: item.Timestamp,
		}
		var t1, t2, t3 model.PropertyMessage
		err := json.Unmarshal([]byte(item.LastStatus), &t1)
		if err == nil && t1.Timestamp > 0 {
			status.LastStatus = &t1
		}
		err = json.Unmarshal([]byte(item.ZeroStatus), &t2)
		if err == nil && t2.Timestamp > 0 {
			status.ZeroStatus = &t2
		}
		err = json.Unmarshal([]byte(item.PreDayStatus), &t3)
		if err == nil && t3.Timestamp > 0 {
			status.PreDayStatus = &t3
		}
		tmp1[item.DeviceId] = status
	}
	deviceStatusLock.Lock()
	defer deviceStatusLock.Unlock()
	cache_deviceStatus = tmp1
}

func CacheDataAsyncSave() {
	//数据库记录检查并清理
	logrus.Info("CacheDataAsyncSave clean data")
	model.CheckChan <- struct{}{}
	model.PushStatusChan <- struct{}{}

	logrus.Info("SaveGatewayApiHistory sync status")
	// hasSync := false
	for gatewayId, val := range gatewayApiRecords {
		// hasSync = false
		for api, val1 := range val {
			for deviceId, apiStatus := range val1 {
				err := db.SaveGatewayApiHistory(db.GatewayApiHistory{
					Name:       api,
					GatewayId:  gatewayId,
					DeviceId:   deviceId,
					Timestamp:  apiStatus.PreTimestamp,
					UpdateTime: time.Unix(apiStatus.Timestamp, 0),
					Status:     apiStatus.Status,
				})
				if err != nil {
					logrus.Errorf("SaveGatewayApiHistory error:%+v,%+v", apiStatus, err)
				}
				// if deviceId > 0 {
				// 	model.StatusMsgChan <- model.StatusMsg{
				// 		DeviceId: deviceId,
				// 		Status:   model.STATUS_ACTIVE,
				// 	}
				// } else if !hasSync {
				// 	devices := GetGatewayDevices(gatewayId)
				// 	for _, d := range devices {
				// 		model.StatusMsgChan <- model.StatusMsg{
				// 			DeviceId: d.Id,
				// 			Status:   model.STATUS_ACTIVE,
				// 		}
				// 	}
				// 	hasSync = true
				// }
			}
		}
	}

	logrus.Info("device status save")

	cache_deviceStatus_tmp := cache_deviceStatus
	for _, status := range cache_deviceStatus_tmp {
		err := db.SaveDeviceStatus(db.DeviceStatus{
			DeviceId:     status.Id,
			Status:       status.Status,
			Timestamp:    status.Timestamp,
			LastStatus:   utils.ToString(status.LastStatus),
			ZeroStatus:   utils.ToString(status.ZeroStatus),
			PreDayStatus: utils.ToString(status.PreDayStatus),
		})
		if err != nil {
			logrus.Errorf("SaveDeviceStatus error:%+v,%+v", status, err)
		}
	}
}
