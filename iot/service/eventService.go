package service

import (
	"errors"
	"sync"
	"time"

	"koudai-box/cache"

	"koudai-box/iot/db"
	"koudai-box/iot/web/common"
	"koudai-box/iot/web/dto"

	"github.com/sirupsen/logrus"
)

const (
	EVENT_CACHE_KEY string = "eventCache"
)

var eventLock = sync.Mutex{}

func AddEventService(request dto.SaveEventRequest) (int64, error) {
	eventLock.Lock()
	defer eventLock.Unlock()

	stamp, _ := time.ParseInLocation(timeTemplate, request.CreateTime, time.Local)

	event := db.Event{
		Sign:       request.Sign,
		Title:      request.Title,
		Type:       request.Type,
		Level:      request.Level,
		DeviceId:   request.DeviceId,
		ProductId:  request.ProductId,
		CreateTime: stamp,
	}
	autoIncEventId, err := db.InsertEvent(event)
	ClearEventCache()
	return autoIncEventId, err
}

func UpdateEventService(request dto.UpdateEventRequest) error {
	event := GetEventFromCache(request.Id)
	if event == nil {
		return errors.New("事件不存在")
	}
	eventLock.Lock()
	defer eventLock.Unlock()

	stamp, _ := time.ParseInLocation(timeTemplate, request.CreateTime, time.Local)

	db.UpdateEvent(&db.Event{
		Id:         request.Id,
		Sign:       request.Sign,
		Title:      request.Title,
		Type:       request.Type,
		Level:      request.Level,
		DeviceId:   request.DeviceId,
		ProductId:  request.ProductId,
		CreateTime: stamp,
	})
	return nil
}

func DeleteEventService(ids []int) error {
	eventLock.Lock()
	defer eventLock.Unlock()

	for _, c := range ids {
		err := deleteOneEvent(c)
		if err != nil {
			logrus.Error(err)
			continue
		}
	}
	ClearEventCache()
	return nil
}

func QueryEventSerivce(request dto.QueryEventDataRequest) (int64, []*dto.EventItem) {
	offset, limit := common.Page2Offset(request.PageNo, request.PageSize)
	totalSize, events := db.QueryEventsByPage(offset, limit, request.Search, request.DeviceId, request.Level, request.StartTime, request.EndTime)
	var eventItems []*dto.EventItem
	for _, event := range events {
		eventItem := fixEventInfo(event)
		eventItems = append(eventItems, &eventItem)
	}
	return totalSize, eventItems
}

func deleteOneEvent(eventId int) error {
	gateway := GetEventFromCache(eventId)
	if gateway == nil {
		return errors.New("事件不存在")
	}
	err := db.DeleteEvent(eventId)
	if err != nil {
		return errors.New("删除失败")
	}
	return nil

}

func ClearEventCache() {
	cache.Delete(EVENT_CACHE_KEY)
}

func GetEventFromCache(eventId int) *dto.EventItem {
	c := GetEventCache()[eventId]
	return c
}

func GetEventCache() map[int]*dto.EventItem {
	m, err := cache.Get(EVENT_CACHE_KEY)
	if err != nil {
		InitEventCache()
		m, _ = cache.Get(EVENT_CACHE_KEY)
		if m == nil {
			return make(map[int]*dto.EventItem)
		} else {
			return m.(map[int]*dto.EventItem)
		}
	}
	return m.(map[int]*dto.EventItem)
}

func InitEventCache() {
	events := ListAllEvent()
	eventMap := make(map[int]*dto.EventItem)
	for _, c := range events {
		eventMap[c.EventId] = c
	}
	err := SetEventCache(eventMap)
	if err != nil {
		logrus.Errorln("缓存事件数据失败:", err)
	}
}

func SetEventCache(value map[int]*dto.EventItem) error {
	return cache.SetWithNoExpire(EVENT_CACHE_KEY, value)
}

func ListAllEvent() []*dto.EventItem {
	eventItems := make([]*dto.EventItem, 0)
	_, events := db.QueryAllEvents()
	for _, event := range events {
		eventItem := fixEventInfo(event)
		eventItems = append(eventItems, &eventItem)
	}
	return eventItems
}

func QueryEventByIDService(eventID int) *dto.EventItem {
	event := GetEventFromCache(eventID)
	return event
}

func fixEventInfo(event *db.Event) dto.EventItem {
	eventItem := dto.EventItem{
		EventId:     event.Id,
		EventSign:   event.Sign,
		EventTitle:  event.Title,
		EventType:   event.Type,
		EventLevel:  event.Level,
		DeviceId:    event.DeviceId,
		DeviceName:  event.DeviceName,
		DeviceSign:  event.DeviceSign,
		ProductId:   event.ProductId,
		ProductName: event.ProductName,
		CreateTime:  event.CreateTime.Local().Format(timeTemplate),
		MessageId:   event.MessageId,
		Message:     event.Message,
		Timestamp:   event.Timestamp,
	}
	return eventItem
}
