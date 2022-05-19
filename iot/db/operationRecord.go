package db

import (
	"koudai-box/global"
	"time"

	"github.com/astaxie/beego/orm"
)

func InsertOperationRecord(record OperationRecord) (int64, error) {
	id, err := webOrm.Insert(&record)
	if err != nil {
		logger.Errorln(err)
		id = -1
	}
	return id, err
}

func QueryOperationRecordByPage(offset, limit int, search, deviceId, level, startTime, endTime string) (int64, []*Event) {
	var childrenItem []*Event
	querySelector := webOrm.QueryTable("operation_record")

	cond := orm.NewCondition()

	if len(search) > 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.Or("code__contains", search).Or("name__contains", search)
		cond = cond.AndCond(cond1)
	}
	if len(deviceId) > 0 {
		cond = cond.And("deviceId", deviceId)
	}
	if len(level) > 0 {
		cond = cond.And("level", level)
	}
	if len(startTime) > 0 && len(endTime) > 0 {
		var timeTemplate = global.TIME_TEMPLATE
		stamp, _ := time.ParseInLocation(timeTemplate, startTime, time.Local)
		etamp, _ := time.ParseInLocation(timeTemplate, endTime, time.Local)

		cond1 := orm.NewCondition()
		cond1 = cond1.And("create_time__gte", stamp.Unix()).And("create_time__gte", etamp.Unix())
		cond = cond.AndCond(cond1)
	}

	_, err := querySelector.SetCond(cond).Limit(limit, offset).All(&childrenItem)
	if err != nil {
		logger.Errorln(err)
	}
	num, _ := querySelector.Count()
	return num, childrenItem
}
