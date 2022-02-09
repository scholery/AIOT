package service

import (
	"koudai-box/web/sysconfig"
	"sync"
	"time"

	"koudai-box/iot/db"
	"koudai-box/iot/gateway/model"
	status "koudai-box/iot/gateway/status"

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
		SetDeviceRunningStatus(msg.DeviceId, msg.Status)
	}
}
