package db

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/sirupsen/logrus"
)

func InsertProperty(poprs DeviceProperty) (int64, error) {
	id, err := webOrm.Insert(&poprs)
	if err != nil {
		logger.Errorln(err)
		id = -1
	}
	return id, err
}

func GetLastestProperty(deviceId string) (*DeviceProperty, error) {
	var prop DeviceProperty
	qs := webOrm.QueryTable("device_property")
	err := qs.Filter("device_id", deviceId).OrderBy("-timestamp").One(&prop)
	if err != nil {
		return nil, err
	}
	return &prop, nil
}

func GetProperties(deviceId string, count, begin, end int64) ([]*DeviceProperty, error) {
	var props []*DeviceProperty
	qs := webOrm.QueryTable("device_property")
	filter := qs.Filter("device_id", deviceId)
	if begin > 0 {
		filter = filter.Filter("timestamp__gte", begin)
	}
	if end > 0 {
		filter = filter.Filter("timestamp__lte", end)
	}
	if count > 0 {
		num, _ := filter.Count()
		offset := num - count
		if offset < 0 {
			offset = 0
		}
		filter = filter.Limit(count, offset)
	}
	_, err := filter.All(&props)
	if err != nil {
		return nil, err
	}
	return props, nil
}

func DeleteOldProps(capacity int) error {
	var item DeviceProperty
	querySelector := webOrm.QueryTable("device_property")

	err := querySelector.OrderBy("-timestamp").Offset(capacity).One(&item)
	if err != nil {
		return err
	}
	res, err := webOrm.Raw("delete from device_property where timestamp <=?", item.Timestamp).Exec()
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	logrus.Debugf("delete device_property count:%d", count)
	return nil
}

func DevicdePropsPushed(times [2]int64, failIds []string) (int64, bool) {
	if times[0] == 0 {
		return 0, false
	}
	querySelector := webOrm.QueryTable("device_property")
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

func GetPredayProps(deviceId int, day time.Time) []*DeviceProperty {
	var props []*DeviceProperty
	qs := webOrm.QueryTable("device_property")
	y, m, d := day.Date()
	end := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	start := end.AddDate(0, 0, -1)
	_, err := qs.Filter("device_id", deviceId).Filter("timestamp__gte", start.Unix()).Filter("timestamp__lt", end.Unix()).All(&props)
	if err != nil {
		return nil
	}
	return props
}
