package service

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"koudai-box/iot/db"
	"koudai-box/iot/gateway/model"
	"koudai-box/iot/gateway/utils"
	"koudai-box/iot/web/common"
	"koudai-box/iot/web/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var deviceLock = sync.Mutex{}

//添加设备
func AddDeviceService(request dto.AddDeviceRequest) (int64, error) {
	deviceLock.Lock()
	defer deviceLock.Unlock()

	device := db.Device{
		Name:           request.Name,
		Code:           request.Code,
		SourceId:       request.SourceId,
		Geo:            request.Geo,
		Locale:         request.Locale,
		Desc:           request.Desc,
		ExtProps:       utils.ToString(request.ExtProps),
		ProductId:      request.ProductId,
		ActivateStatus: 0,
		RunningStatus:  0,
		DelFlag:        0,
	}

	autoIncDeviceId, err := db.InsertDevice(device)
	device.Id = int(autoIncDeviceId)
	//推送
	PushDevice(&device)
	return autoIncDeviceId, err
}

//更新
func UpdateDeviceService(request dto.UpdateDeviceRequest) (int, error) {
	device, err := db.QueryDeviceByID(request.Id)
	if err != nil {
		return 0, err
	}
	deviceLock.Lock()
	defer deviceLock.Unlock()

	device.Name = request.Name
	device.Code = request.Code
	device.ProductId = request.ProductId
	device.Desc = request.Desc
	device.SourceId = request.SourceId
	device.Geo = request.Geo
	device.Locale = request.Locale
	device.ExtProps = utils.ToString(request.ExtProps)

	err = db.UpdateDevice(device)
	if err != nil {
		return 0, err
	}
	//推送
	PushDevice(device)
	return device.Id, err
}

//查询设备列表
func QueryDeviceSerivce(request dto.QueryDeviceDataRequest) (int64, []gin.H) {
	offset, limit := common.Page2Offset(request.PageNo, request.PageSize)
	totalSize, devices := db.QueryDeviceByPage(offset, limit, request.Search, request.ActivateStatus, request.RunningStatus, request.ProductId)
	var list []gin.H = make([]gin.H, 0)

	if totalSize == 0 {
		return totalSize, list
	}

	var productIds []int = make([]int, 0)
	for _, device := range devices {
		productIds = append(productIds, device.ProductId)
	}
	if len(productIds) == 0 {
		return 0, list
	}
	//查询产品
	products := db.QueryProductByIDs(productIds)
	for _, device := range devices {

		var selectProduct *db.Product
		for _, product := range products {
			if product.Id == device.ProductId {
				selectProduct = product
				break
			}
		}

		logrus.Infoln(device.CreateTime)
		item := gin.H{
			"id":             device.Id,
			"name":           device.Name,
			"activateStatus": device.ActivateStatus,
			"runningStatus":  device.RunningStatus,
			"code":           device.Code,
			"SourceId":       device.SourceId,
			"Geo":            device.Geo,
			"Locale":         device.Locale,
			"createTime":     device.CreateTime.Local().Format("2006-01-02 15:04:05"),
			"updateTime":     device.UpdateTime.Local().Format("2006-01-02 15:04:05"),
			"desc":           device.Desc,
			"ExtProps":       utils.ToMap(device.ExtProps),
		}
		if selectProduct != nil {
			item["productName"] = selectProduct.Name
			item["productCode"] = selectProduct.Code
			item["productId"] = selectProduct.Id
		}

		list = append(list, item)
	}

	return totalSize, list
}

func QueryDeviceIdsByActivateStatus(activateStatus int) []int {
	devices := db.QueryDevicesByStatus(activateStatus, -1)
	var list []int = make([]int, 0)
	for _, device := range devices {
		list = append(list, device.Id)
	}
	return list
}

func QueryDeviceById(deviceId int) (*db.Device, error) {
	deviceInfo, err := db.QueryDeviceByID(deviceId)
	if err != nil {
		return nil, err
	}
	return deviceInfo, err
}

func QueryDeviceInfoByID(deviceId int) (gin.H, error) {
	deviceInfo, err := db.QueryDeviceByID(deviceId)
	if err != nil {
		return nil, err
	}
	item := gin.H{
		"id":             deviceInfo.Id,
		"name":           deviceInfo.Name,
		"activateStatus": deviceInfo.ActivateStatus,
		"runningStatus":  deviceInfo.RunningStatus,
		"code":           deviceInfo.Code,
		"SourceId":       deviceInfo.SourceId,
		"Geo":            deviceInfo.Geo,
		"Locale":         deviceInfo.Locale,
		"createTime":     deviceInfo.CreateTime.Local().Format("2006-01-02 15:04:05"),
		"updateTime":     deviceInfo.UpdateTime.Local().Format("2006-01-02 15:04:05"),
		"desc":           deviceInfo.Desc,
		"ExtProps":       utils.ToMap(deviceInfo.ExtProps),
	}
	if deviceInfo.ProductId != 0 {
		selectProduct, _ := db.QueryProductByID(deviceInfo.ProductId)
		// if err != nil {
		// 	return nil, errors.New("查询产品失败")
		// }
		if selectProduct != nil {
			item["productName"] = selectProduct.Name
			item["productCode"] = selectProduct.Code
			item["image"] = selectProduct.Image
			item["productId"] = selectProduct.Id

			//产品品类
			categorys := strings.Split(selectProduct.Category, ",")
			lastCategory := categorys[len(categorys)-1]
			dictList, _ := db.ListDictAll()
			for _, d := range dictList {
				if lastCategory == d.Value {
					item["categoryName"] = d.Name
					break
				}
			}

			if selectProduct.GatewayId != 0 {
				//网关
				gateway, _ := db.QueryGatewayByID(selectProduct.GatewayId)
				if gateway != nil {
					item["gatewayName"] = gateway.Name
					item["gateWayId"] = selectProduct.GatewayId
					for _, d := range dictList {
						if gateway.Protocol == d.Value {
							item["gatewayProtocol"] = gateway.Protocol
							item["gatewayProtocolName"] = d.Name
							break
						}
					}
				}
			}
		}
	}

	return item, nil
}

//统计设备信息
func StatisticsDeviceSerivce() gin.H {
	total := db.QueryAllDeviceCount()
	publicCount := db.QueryAllDeviceByStateCount(1)
	noPublicCount := db.QueryAllDeviceByStateCount(0)
	offlineCount := db.QueryAllDeviceByOnlineCount(0)
	onlineCount := db.QueryAllDeviceByOnlineCount(1)
	// fmt.Println("offlineCount", offlineCount)
	return gin.H{
		"total":         total,
		"publicCount":   publicCount,
		"noPublicCount": noPublicCount,
		"offlineCount":  offlineCount,
		"onlineCount":   onlineCount,
	}
}

//删除项目
func DeleteDeviceSerivce(id int) error {
	device, err := db.QueryDeviceByID(id)
	if err != nil {
		return err
	}
	if device.RunningStatus == 1 {
		return errors.New("设备未停止")
	}
	deviceLock.Lock()
	defer deviceLock.Unlock()

	//设置
	device.DelFlag = 1
	//推送
	PushDevice(device)
	return db.UpdateDevice(device)
}

//设置状态
func SetDeviceStatusSerivce(id, activateStatus int) error {
	device, err := db.QueryDeviceByID(id)
	if err != nil {
		return err
	}
	deviceLock.Lock()
	defer deviceLock.Unlock()

	//设置
	device.ActivateStatus = activateStatus
	device.RunningStatus = model.STATUS_DISACTIVE
	//推送
	PushDevice(device)

	return db.UpdateDevice(device)
}

//设置运行状态
func SetDeviceRunningStatus(id, runningStatus int) error {
	device, err := db.QueryDeviceByID(id)
	if err != nil {
		return err
	}
	deviceLock.Lock()
	defer deviceLock.Unlock()

	//设置
	device.RunningStatus = runningStatus
	//推送
	PushDevice(device)

	return db.UpdateDevice(device)
}

// func SetDeviceActivateStatus(deviceId, activateStatus int) error {
// 	device, err := db.QueryDeviceByID(deviceId)
// 	if err != nil {
// 		return errors.New("查询设备失败")
// 	}
// 	deviceLock.Lock()
// 	defer deviceLock.Unlock()

// 	if device == nil {
// 		return errors.New("未查询到设备")
// 	}

// 	if device.ActivateStatus == activateStatus {
// 		return nil
// 	}

// 	ok := false
// 	if activateStatus == 0 {
// 		ok = iot.StopDevice(deviceId)
// 	} else {
// 		ok = iot.StartDevice(deviceId)
// 	}
// 	if !ok {
// 		return errors.New("设备启停失败")
// 	}

// 	err = SetDeviceStatusSerivce(deviceId, activateStatus)
// 	if err != nil {
// 		return errors.New("设备状态更新失败")
// 	}
// 	return nil
// }

func PushAllDevices() error {
	devices, err := db.QueryDevices()
	if err != nil {
		return nil
	}
	for _, device := range devices {
		PushDevice(device)
	}
	return nil
}

func PushDevice(device *db.Device) {
	selectProduct, _ := db.QueryProductByID(device.ProductId)
	devMsg := model.DeviceMessage{
		DeviceId:       device.Id,
		Key:            device.Code,
		Name:           device.Name,
		SourceId:       device.SourceId,
		Geo:            device.Geo,
		Locale:         device.Locale,
		ActivateStatus: device.ActivateStatus,
		RunningStatus:  device.RunningStatus,
		Desc:           device.Desc,
		ExtProps:       utils.ToMap(device.ExtProps),
		UpdateTime:     device.UpdateTime.Format("2006-01-02 15:04:05"),
		DelFlag:        device.DelFlag,
	}
	if selectProduct != nil {
		devMsg.Category = selectProduct.Category
	}
	model.PushOutMsgChan <- model.Message{
		SN:    utils.GetSN(),
		Type:  model.Message_Type_Device,
		Msg:   devMsg,
		MsgId: strconv.FormatInt(int64(devMsg.DeviceId), 10),
	}
}
