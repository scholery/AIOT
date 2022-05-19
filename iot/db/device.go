package db

import (
	"fmt"
	"koudai-box/iot/gateway/model"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

func InsertDevice(device Device) (int64, error) {
	id, err := webOrm.Insert(&device)
	if err != nil {
		logger.Errorln(err)
		id = -1
	}
	return id, err
}

func QueryDevicetByProductIds(productIds []int) []*Device {
	var devices []*Device
	qs := webOrm.QueryTable("device")

	cond := orm.NewCondition()
	cond = cond.And("productId__in", productIds).And("del_flag", 0)
	qs = qs.SetCond(cond)

	qs.All(&devices)
	return devices
}

//更加id查询
func QueryDeviceByID(id int) (*Device, error) {
	var device Device
	err := webOrm.QueryTable("device").Filter("id", id).Filter("del_flag", 0).One(&device)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	return &device, nil
}

//更新
func UpdateDevice(device *Device) error {
	_, err := webOrm.Update(device)
	if err != nil {
		logger.Errorln(err)
		webOrm.Rollback()
	}
	return err
}

//查询数据库
func QueryDeviceByPage(offset, limit int, search string, activateStatus int, runningStatus int, productId int) (int64, []*Device) {
	var devices []*Device
	qs := webOrm.QueryTable("device")
	cond := orm.NewCondition()
	cond = cond.And("del_flag", 0)

	if activateStatus != model.STATUS_ALL {
		cond = cond.And("activateStatus", activateStatus)
	}

	if runningStatus != model.STATUS_ALL {
		cond = cond.And("runningStatus", runningStatus)
	}

	if productId != 0 {
		cond = cond.And("product_id", productId)
	}

	if len(search) > 0 {
		cond1 := orm.NewCondition()
		cond1 = cond1.Or("code__contains", search).Or("name__contains", search)
		cond = cond.AndCond(cond1)
	}

	qs = qs.SetCond(cond)
	_, err := qs.Limit(limit, offset).OrderBy("-createTime").All(&devices)
	if err != nil {
		logger.Errorln(err)
	}
	num, _ := qs.Count()
	return num, devices
}

func QueryDevices() ([]*Device, error) {
	var devices []*Device
	qs := webOrm.QueryTable("device")
	_, err := qs.Filter("del_flag", 0).OrderBy("-createTime").All(&devices)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	return devices, nil
}

func QueryDevicesByStatus(activateStatus int, runningStatus int) []*Device {
	var devices []*Device
	qs := webOrm.QueryTable("device")
	cond := orm.NewCondition()
	cond = cond.And("del_flag", 0)

	if activateStatus != model.STATUS_ALL {
		cond = cond.And("activateStatus", activateStatus)
	}

	if runningStatus != model.STATUS_ALL {
		cond = cond.And("runningStatus", runningStatus)
	}
	qs = qs.SetCond(cond)
	_, err := qs.OrderBy("-createTime").All(&devices)
	if err != nil {
		logger.Errorln(err)
	}
	return devices
}

//查询所有产品数量
func QueryAllDeviceCount() int64 {
	num, err := webOrm.QueryTable("device").Filter("del_flag", 0).Count()
	if err != nil {
		logger.Errorln(err)
		return 0
	}
	return num
}

//查询所有产品数量
func QueryAllDeviceByStateCount(state int) int64 {
	num, err := webOrm.QueryTable("device").Filter("activateStatus", state).Filter("del_flag", 0).Count()
	if err != nil {
		logger.Errorln(err)
		return 0
	}
	return num
}

//查询所有产品数量
func QueryAllDeviceByOnlineCount(online int) int64 {
	num, err := webOrm.QueryTable("device").Filter("runningStatus", online).Filter("del_flag", 0).Count()
	if err != nil {
		logger.Errorln(err)
		return 0
	}
	return num
}

//根据ids查询
func QueryDevicesByIds(ids []int) []orm.ParamsList {
	o := orm.NewOrm()
	var lists []orm.ParamsList
	slice1 := make([]string, 0)
	for _, id := range ids {
		str := strconv.Itoa(id)
		slice1 = append(slice1, str)
	}
	slice2 := strings.Join(slice1, ",")
	slice2 = fmt.Sprintf("(%s)", slice2)
	o.Raw("SELECT a.code ,a.name, b.name as product_name, b.code as product_code , b.desc as desc, strftime('%Y-%m-%d %H:%M:%S',a.create_time) as create_time FROM device a LEFT JOIN product b on a.product_id = b.id where a.id IN " + slice2).ValuesList(&lists)
	return lists
}

//根据产品名查询产品
func QueryProductByName(name string) (*Product, error) {
	var product Product
	err := webOrm.QueryTable("product").Filter("name", name).Filter("del_flag", 0).One(&product)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	return &product, nil
}
