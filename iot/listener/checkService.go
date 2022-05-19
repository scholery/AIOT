package listener

import (
	"fmt"
	"koudai-box/global"
	"koudai-box/web/sysconfig"
	"strconv"
	"sync"
	"time"

	"koudai-box/iot/db"
	"koudai-box/iot/gateway/model"
	gatewayService "koudai-box/iot/gateway/service"
	status "koudai-box/iot/gateway/status"
	"koudai-box/iot/gateway/utils"
	"koudai-box/iot/service"

	"github.com/sirupsen/logrus"
)

var capacityLimitLock sync.RWMutex

func CheckListener() {
	for range model.CheckChan {
		CheckCapacityLimit()
	}
}

func CheckCapacityLimit() {
	capacityLimitLock.Lock()
	defer capacityLimitLock.Unlock()
	storageConfig := sysconfig.GetStorageConfig()
	err := db.DeleteOldAlarms(storageConfig.IotAlarmMaxCount)
	if err != nil {
		logrus.Errorf("DeleteOldAlarms error:%+v", err)
	}
	err = db.DeleteOldEvents(storageConfig.IotEventMaxCount)
	if err != nil {
		logrus.Errorf("DeleteOldEvents error:%+v", err)
	}
	err = db.DeleteOldProps(storageConfig.IotPropsMaxCount)
	if err != nil {
		logrus.Errorf("DeleteOldProps error:%+v", err)
	}
}

func DeviceStatusUpdate() {
	for msg := range model.StatusMsgChan {
		tmp := status.GetDeviceStatus(msg.DeviceId)
		logrus.Debugf("DeviceStatusUpdate-L(%d) before:%d,msgstatus:%d", msg.DeviceId, tmp.Timestamp, msg.Status)
		if tmp.Status == msg.Status {
			continue
		}
		tmp.Timestamp = time.Now().Unix()
		tmp.Status = msg.Status
		tmp1 := status.GetDeviceStatus(msg.DeviceId)
		logrus.Debugf("DeviceStatusUpdate-L(%d) after:%d", msg.DeviceId, tmp1.Timestamp)
		service.SetDeviceRunningStatus(msg.DeviceId, msg.Status)
		//上线、下限事件
		generateAlarm(tmp)
	}
}

/**
* 上线、下限事件
**/
func generateAlarm(statusMsg *model.DeviceStatus) (bool, error) {
	device := status.GetDevice(statusMsg.Id)
	if device == nil {
		return false, fmt.Errorf("device %s is not start", statusMsg.Key)
	}
	gateway := &gatewayService.DataGateway{Device: device}
	var title, code, msg string
	if statusMsg.Status == model.STATUS_ACTIVE {
		title = fmt.Sprintf("设备[%s]上线", device.Name)
		code = "online"
		msg = fmt.Sprintf("设备[%s]上线，上线时间：%s", device.Name, time.Now().Local().Format(global.TIME_TEMPLATE))
	} else if statusMsg.Status == model.STATUS_DISACTIVE {
		title = fmt.Sprintf("设备[%s]下线", device.Name)
		code = "offline"
		msg = fmt.Sprintf("设备[%s]下线，下线时间：%s", device.Name, time.Now().Local().Format(global.TIME_TEMPLATE))
	} else {
		return false, nil
	}
	alarm := model.IotEventMessage{DeviceId: strconv.Itoa(device.Id), DeviceSign: device.Key, MessageId: utils.GetUUID(), Code: code,
		Timestamp: time.Now().Unix(), Type: model.Alarm_Type_Event, Level: "common", Title: title, Message: msg}
	return gateway.LoaderAlarm(alarm, true)
}
