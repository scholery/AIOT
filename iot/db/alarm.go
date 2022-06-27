package db

import (
	"koudai-box/global"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/sirupsen/logrus"
)

func InsertAlarm(alarm Alarm) (int64, error) {
	id, err := webOrm.Insert(&alarm)
	if err != nil {
		logger.Errorln(err)
		id = -1
	}
	return id, err
}

func UpdateAlarm(alarm *Alarm) (int64, error) {
	webOrm.Begin()
	num, err := webOrm.Update(alarm, "title", "sign", "type", "level", "device_id", "product_id", "create_time")
	if err != nil {
		webOrm.Rollback()
	}
	webOrm.Commit()
	return num, err
}

func DeleteAlarm(alarmId int) error {
	_, err := webOrm.Delete(&Alarm{Id: alarmId}, "id")
	if err != nil {
		webOrm.Rollback()
	}
	return err
}

func DeleteAlarmByAlarmId(alarmId string) error {
	_, err := webOrm.Delete(&Alarm{MessageId: alarmId}, "message_id")
	if err != nil {
		webOrm.Rollback()
	}
	return err
}

func QueryAlarmsByPage(offset, limit int, search, productId, productName, deviceId, deviceName, alarmType, level, startTime, endTime string) (int64, []*Alarm) {
	var childrenItem []*Alarm
	querySelector := webOrm.QueryTable("alarm")

	cond := orm.NewCondition()

	if len(search) > 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.Or("code__contains", search).Or("title__contains", search)
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
	if len(alarmType) > 0 {
		cond = cond.And("type", alarmType)
	}
	if len(level) > 0 {
		cond = cond.And("level", level)
	}
	if len(startTime) > 0 && len(endTime) > 0 {
		var timeTemplate = global.TIME_TEMPLATE
		stamp, _ := time.ParseInLocation(timeTemplate, startTime, time.Local)
		etamp, _ := time.ParseInLocation(timeTemplate, endTime, time.Local)

		cond1 := orm.NewCondition()
		cond1 = cond1.And("timestamp__gte", stamp.Unix()).And("timestamp__lte", etamp.Unix())
		cond = cond.AndCond(cond1)
	}

	_, err := querySelector.SetCond(cond).Limit(limit, offset).OrderBy("-timestamp").All(&childrenItem)
	if err != nil {
		logger.Errorln(err)
	}
	num, _ := querySelector.SetCond(cond).Count()
	return num, childrenItem
}

func QueryAllAlarms() (int64, []*Alarm) {
	var childrenItem []*Alarm
	qs := webOrm.QueryTable("alarm")
	num, err := qs.All(&childrenItem)
	if err != nil {
		logger.Errorln(err)
	}
	return num, childrenItem
}

//总告警数量
func CountTotalAlarms() int64 {
	qs := webOrm.QueryTable("alarm")
	cond := orm.NewCondition()
	count, _ := qs.SetCond(cond).Count()
	return count
}

//当日告警数量
func CountTodayAlarms() int64 {
	qs := webOrm.QueryTable("alarm")

	cond := orm.NewCondition()
	todayTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	cond = cond.And("timestamp__gte", todayTime.Unix()).And("timestamp__lte", time.Now().Unix())
	count, err := qs.SetCond(cond).Count()
	if err != nil {
		logger.Error(err)
	}
	return count
}

//当日告警最多的设备
func CountTodayMostAlarmDeviceName() string {
	o := orm.NewOrm()
	res := make(orm.Params)
	todayTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	num, err := o.Raw("SELECT name as name, count_alarm as value from (SELECT device_id, count(id) AS count_alarm FROM alarm  where timestamp >=? GROUP BY device_id ORDER BY count_alarm desc limit 1) a LEFT JOIN device b on a.device_id = b.id", todayTime.Unix()).RowsToMap(&res, "name", "value")
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

//当日告警最多的名称
func CountTodayMostAlarmName() string {
	o := orm.NewOrm()
	res := make(orm.Params)
	todayTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	num, err := o.Raw("SELECT title as name, count(id) AS value FROM alarm  where timestamp >= ? GROUP BY code ORDER BY value desc limit 1", todayTime.Unix()).RowsToMap(&res, "name", "value")
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

func DeleteOldAlarms(capacity int) error {
	var item Alarm
	querySelector := webOrm.QueryTable("alarm")

	err := querySelector.OrderBy("-timestamp").Offset(capacity).One(&item)
	if err != nil {
		return err
	}
	res, err := webOrm.Raw("delete from alarm where timestamp <=?", item.Timestamp).Exec()
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	logrus.Debugf("delete alarm count:%d", count)
	return nil
}

func AlarmPushed(times [2]int64, failIds []string) (int64, bool) {
	if times[0] == 0 {
		return 0, false
	}
	querySelector := webOrm.QueryTable("alarm")
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
