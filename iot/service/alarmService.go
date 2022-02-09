package service

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"time"

	"koudai-box/cache"

	"koudai-box/iot/db"
	"koudai-box/iot/web/common"
	"koudai-box/iot/web/dto"

	"github.com/sirupsen/logrus"
)

const (
	ALARM_CACHE_KEY string = "alarmCache"
)

var alarmLock = sync.Mutex{}
var timeTemplate = "2006-01-02 15:04:05"

func AddAlarmService(request dto.SaveAlarmRequest) (int64, error) {
	alarmLock.Lock()
	defer alarmLock.Unlock()

	stamp, _ := time.ParseInLocation(timeTemplate, request.CreateTime, time.Local)

	alarm := db.Alarm{
		Code:       request.Code,
		Title:      request.Title,
		Type:       request.Type,
		Level:      request.Level,
		DeviceId:   request.DeviceId,
		CreateTime: stamp,
	}
	autoIncAlarmId, err := db.InsertAlarm(alarm)
	ClearAlarmCache()
	return autoIncAlarmId, err
}

func UpdateAlarmService(request dto.UpdateAlarmRequest) error {
	alarm := GetAlarmFromCache(request.Id)
	if alarm == nil {
		return errors.New("告警不存在")
	}
	alarmLock.Lock()
	defer alarmLock.Unlock()

	stamp, _ := time.ParseInLocation(timeTemplate, request.CreateTime, time.Local)

	db.UpdateAlarm(&db.Alarm{
		Id:         request.Id,
		Code:       request.Code,
		Title:      request.Title,
		Type:       request.Type,
		Level:      request.Level,
		DeviceId:   request.DeviceId,
		CreateTime: stamp,
	})
	return nil
}

func DeleteAlarmService(ids []string) error {
	alarmLock.Lock()
	defer alarmLock.Unlock()

	for _, c := range ids {
		err := deleteOneAlarm(c)
		if err != nil {
			logrus.Error(err)
			continue
		}
	}
	ClearAlarmCache()
	return nil
}

func QueryAlarmSerivce(request dto.QueryAlarmDataRequest) (int64, []*dto.AlarmItem) {
	offset, limit := common.Page2Offset(request.PageNo, request.PageSize)
	totalSize, alarms := db.QueryAlarmsByPage(offset, limit, request.Search, request.DeviceId, request.Level, request.StartTime, request.EndTime)
	var alarmItems []*dto.AlarmItem
	for _, alarm := range alarms {
		alarmItem := fixAlarmInfo(alarm)
		alarmItems = append(alarmItems, &alarmItem)
	}
	return totalSize, alarmItems
}

func deleteOneAlarm(alarmId string) error {
	gateway := GetAlarmFromCacheByMsgId(alarmId)
	if gateway == nil {
		return errors.New("告警不存在")
	}
	err := db.DeleteAlarmByAlarmId(alarmId)
	if err != nil {
		return errors.New("删除失败")
	}
	return nil

}

func ClearAlarmCache() {
	cache.Delete(ALARM_CACHE_KEY)
}

func GetAlarmFromCache(alarmId int) *dto.AlarmItem {
	c := GetAlarmCache()[strconv.Itoa(alarmId)]
	return c
}

func GetAlarmFromCacheByMsgId(alarmMsgId string) *dto.AlarmItem {
	c := GetAlarmCache()[alarmMsgId]
	return c
}

func GetAlarmCache() map[string]*dto.AlarmItem {
	m, err := cache.Get(ALARM_CACHE_KEY)
	if err != nil {
		InitAlarmCache()
		m, _ = cache.Get(ALARM_CACHE_KEY)
		if m == nil {
			return make(map[string]*dto.AlarmItem)
		} else {
			return m.(map[string]*dto.AlarmItem)
		}
	}
	return m.(map[string]*dto.AlarmItem)
}

func InitAlarmCache() {
	alarms := ListAllAlarm()
	alarmMap := make(map[string]*dto.AlarmItem)
	for _, c := range alarms {
		alarmMap[c.AlarmId] = c
	}
	err := SetAlarmCache(alarmMap)
	if err != nil {
		logrus.Errorln("缓存告警数据失败:", err)
	}
}

func SetAlarmCache(value map[string]*dto.AlarmItem) error {
	return cache.SetWithNoExpire(ALARM_CACHE_KEY, value)
}

func ListAllAlarm() []*dto.AlarmItem {
	alarmItems := make([]*dto.AlarmItem, 0)
	_, alarms := db.QueryAllAlarms()
	for _, alarm := range alarms {
		alarmItem := fixAlarmInfo(alarm)
		alarmItems = append(alarmItems, &alarmItem)
	}
	return alarmItems
}

func QueryAlarmByIDService(alarmID int) *dto.AlarmItem {
	alarm := GetAlarmFromCache(alarmID)
	return alarm
}

func fixAlarmInfo(alarm *db.Alarm) dto.AlarmItem {
	var props []interface{}
	json.Unmarshal([]byte(alarm.Properties), &props)
	var cons []interface{}
	json.Unmarshal([]byte(alarm.Conditions), &cons)
	alarmItem := dto.AlarmItem{
		AlarmId:     alarm.MessageId,
		AlarmSign:   alarm.Code,
		AlarmTitle:  alarm.Title,
		AlarmType:   alarm.Type,
		AlarmLevel:  alarm.Level,
		DeviceId:    alarm.DeviceId,
		DeviceSign:  alarm.DeviceSign,
		DeviceName:  alarm.DeviceName,
		ProductId:   alarm.ProductId,
		ProductName: alarm.ProductName,
		CreateTime:  time.Unix(alarm.Timestamp, 0).Local().Format(timeTemplate),
		MessageId:   alarm.MessageId,
		Message:     alarm.Message,
		Timestamp:   alarm.Timestamp,
		Props:       props,
		Conditions:  cons,
	}
	return alarmItem
}
