package db

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/sirupsen/logrus"
)

func InsertEvent(event Event) (int64, error) {
	id, err := webOrm.Insert(&event)
	if err != nil {
		logger.Errorln(err)
		id = -1
	}
	return id, err
}

func UpdateEvent(event *Event) (int64, error) {
	num, err := webOrm.Update(event, "title", "sign", "type", "level", "device_id", "product_id", "create_time")
	if err != nil {
		webOrm.Rollback()
	}
	return num, err
}

func DeleteEvent(eventId int) error {
	_, err := webOrm.Delete(&Event{Id: eventId}, "id")
	if err != nil {
		webOrm.Rollback()
	}
	return err
}

func QueryEventsByPage(offset, limit int, search, productId, productName, deviceId, deviceName, eventType, level, startTime string, endTime string) (int64, []*Event) {
	var childrenItem []*Event
	querySelector := webOrm.QueryTable("event")

	cond := orm.NewCondition()

	if len(search) > 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.Or("sign__contains", search).Or("title__contains", search)
		cond = cond.AndCond(cond1)
	}
	if len(productId) > 0 {
		cond = cond.And("productId", productId)
	}
	if len(productName) > 0 {
		cond = cond.And("productName__contains", productName)
	}
	if len(deviceId) > 0 {
		cond = cond.And("deviceId", deviceId)
	}
	if len(deviceName) > 0 {
		cond = cond.And("deviceName__contains", deviceName)
	}
	if len(eventType) > 0 {
		cond = cond.And("type", eventType)
	}
	if len(level) > 0 {
		cond = cond.And("level", level)
	}
	if startTime != "" && endTime != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("create_time__gte", startTime).And("create_time__lte", endTime)
		cond = cond.AndCond(cond1)
	}

	_, err := querySelector.SetCond(cond).Limit(limit, offset).OrderBy("-timestamp").All(&childrenItem)
	if err != nil {
		logger.Errorln(err)
	}
	// logger.Info(childrenItem)
	num, _ := querySelector.SetCond(cond).Count()
	return num, childrenItem
}

func QueryAllEvents() (int64, []*Event) {
	var childrenItem []*Event
	qs := webOrm.QueryTable("event")
	num, err := qs.Filter("del_flag", 0).All(&childrenItem)
	if err != nil {
		logger.Errorln(err)
	}
	return num, childrenItem
}

//总事件数量
func CountTotalEvents() int64 {
	qs := webOrm.QueryTable("event")
	count, _ := qs.Count()
	return count
}

//当日事件数量
func CountTodayEvents() int64 {
	qs := webOrm.QueryTable("event")

	cond := orm.NewCondition()
	todayTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	cond = cond.And("timestamp__gte", todayTime.Unix()).And("timestamp__lte", time.Now().Unix())
	count, err := qs.SetCond(cond).Count()
	if err != nil {
		logger.Error(err)
	}
	return count
}

//当日事件最多的设备
func CountTodayMostEventDeviceName() string {
	o := orm.NewOrm()
	res := make(orm.Params)
	todayTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	num, err := o.Raw("SELECT name as name, count_event as value from (SELECT device_id, count(id) AS count_event FROM event where timestamp >=? GROUP BY device_id ORDER BY count_event desc limit 1) a LEFT JOIN device b on a.device_id = b.id", todayTime.Unix()).RowsToMap(&res, "name", "value")
	if err != nil {
		logger.Error(err)
		return ""
	}
	if num == 0 {
		logger.Error(err)
		return ""
	}
	var deviceName string
	for k, _ := range res {
		deviceName = k
	}
	return deviceName
}

//当日事件最多的名称
func CountTodayMostEventName() string {
	o := orm.NewOrm()
	res := make(orm.Params)
	todayTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	num, err := o.Raw("SELECT title as name, count(id) AS value FROM event where timestamp >=? GROUP BY title ORDER BY value desc limit 1", todayTime.Unix()).RowsToMap(&res, "name", "value")
	if err != nil {
		logger.Error(err)
		return ""
	}
	if num == 0 {
		logger.Error(err)
		return ""
	}
	var alarmCode string
	for k, _ := range res {
		alarmCode = k
	}
	return alarmCode
}

func DeleteOldEvents(capacity int) error {
	var item Event
	querySelector := webOrm.QueryTable("event")
	err := querySelector.OrderBy("-timestamp").Offset(capacity).One(&item)
	if err != nil {
		return err
	}
	res, err := webOrm.Raw("delete from event where timestamp <=?", item.Timestamp).Exec()
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	logrus.Debugf("delete event count:%d", count)
	return nil
}

func EventPushed(times [2]int64, failIds []string) (int64, bool) {
	if times[0] == 0 {
		return 0, false
	}
	querySelector := webOrm.QueryTable("event")
	cond := orm.NewCondition()
	if len(failIds) > 0 {
		cond = cond.AndNot("message_id__in", failIds)
	}
	if times[1] == 0 {
		cond = cond.And("timestamp", times[0])
	} else {
		cond = cond.And("timestamp__gte", times[0]).And("timestamp__lte", times[1])
	}
	num, err := querySelector.SetCond(cond).Update(orm.Params{"push_flag": 1})
	if err != nil {
		logrus.Error(err)
		return 0, false
	}
	return num, true
}
