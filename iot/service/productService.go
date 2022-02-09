package service

import (
	"errors"
	"strings"
	"sync"
	"time"

	"koudai-box/cache"

	"koudai-box/iot/db"
	"koudai-box/iot/web/common"
	"koudai-box/iot/web/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	PRODUCT_CACHE_KEY string = "productCache"
)

var productLock = sync.Mutex{}
var logger = logrus.New()

//添加产品
func AddProductService(request dto.AddProductDataRequest, imagePath string) (int64, error) {
	productLock.Lock()
	defer productLock.Unlock()

	t := time.Now()

	product := db.Product{
		Name:             request.Name,
		Code:             request.Code,
		Image:            imagePath,
		Category:         request.Category,
		GatewayId:        request.GatewayId,
		Desc:             request.Desc,
		PublishTime:      &t,
		Items:            "[]",
		OperationConfigs: "[]",
		EventConfigs:     "[]",
		AlarmConfigs:     "[]",
		FunctionConfigs:  "",
		State:            0,
	}

	autoIncProductId, err := db.InsertProduct(product)
	ClearProductCache()
	return autoIncProductId, err
}

//更新
func UpdateProductService(request dto.UpdateProductDataRequest, imagePath string) (int, error) {
	product, err := db.QueryProductByID(request.Id)
	if err != nil {
		return 0, err
	}
	productLock.Lock()
	defer productLock.Unlock()

	if len(imagePath) > 0 {
		product.Image = imagePath
	}

	product.Name = request.Name
	product.Code = request.Code
	product.Category = request.Category
	product.GatewayId = request.GatewayId
	product.Desc = request.Desc

	err = db.UpdateProduct(product)
	if err != nil {
		return 0, err
	}
	return product.Id, err
}

//查询产品列表
func QueryProductSerivce(request dto.QueryProductDataRequest) (int64, []*dto.ProductItem) {
	var list []*dto.ProductItem = make([]*dto.ProductItem, 0)

	offset, limit := common.Page2Offset(request.PageNo, request.PageSize)
	totalSize, products := db.QueryProductByPage(offset, limit, request.Search, request.State)
	if totalSize == 0 {
		return totalSize, list
	}

	var productIds []int = make([]int, 0)
	for _, product := range products {
		productIds = append(productIds, product.Id)
	}

	//查询设备数量
	devices := db.QueryDevicetByProductIds(productIds)

	for _, product := range products {
		deviceCount := 0
		disabledDeviceCount := 0
		onlineDeviceCount := 0
		offlineDeviceCount := 0

		//设备的状态
		for _, device := range devices {
			if device.ProductId == product.Id {
				deviceCount = deviceCount + 1

				if device.ActivateStatus == 0 {
					disabledDeviceCount = disabledDeviceCount + 1
				}

				if device.RunningStatus == 1 {
					onlineDeviceCount = onlineDeviceCount + 1
				} else {
					offlineDeviceCount = offlineDeviceCount + 1
				}
			}
		}

		productItem := dto.ProductItem{
			Id:                  product.Id,
			State:               product.State,
			Name:                product.Name,
			Desc:                product.Desc,
			Code:                product.Code,
			Image:               product.Image,
			GatewayId:           product.GatewayId,
			Category:            strings.Split(product.Category, ","),
			DeviceCount:         deviceCount,
			DisabledDeviceCount: disabledDeviceCount,
			OnlineDeviceCount:   onlineDeviceCount,
			OfflineDeviceCount:  offlineDeviceCount,
		}
		list = append(list, &productItem)
	}
	return totalSize, list
}

//查询
func DetailProductSerivce(id int) (map[string]interface{}, error) {
	//查询
	product, err := db.QueryProductByID(id)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{
		"id":          product.Id,
		"name":        product.Name,
		"state":       product.State,
		"code":        product.Code,
		"gatewayId":   product.GatewayId,
		"image":       product.Image,
		"category":    strings.Split(product.Category, ","),
		"desc":        product.Desc,
		"createTime":  product.CreateTime.Local().Format("2006-01-02 15:04:05"),
		"publishTime": product.PublishTime.Local().Format("2006-01-02 15:04:05"),
	}

	//设备数量
	data["deviceCount"] = 0

	//产品品类
	categorys := strings.Split(product.Category, ",")
	lastCategory := categorys[len(categorys)-1]
	dictList, _ := db.ListDictAll()
	for _, d := range dictList {
		if lastCategory == d.Value {
			data["categoryName"] = d.Name
			break
		}
	}

	//网关
	gateway, _ := db.QueryGatewayByID(product.GatewayId)
	if gateway != nil {
		data["gatewayName"] = gateway.Name
		for _, d := range dictList {
			if gateway.Protocol == d.Value {
				data["gatewayProtocol"] = gateway.Protocol
				data["gatewayProtocolName"] = d.Name
				break
			}
		}
	}
	return data, err
}

func ClearProductCache() {
	cache.Delete(PRODUCT_CACHE_KEY)
}

//统计产品信息
func StatisticsProductSerivce() gin.H {
	total := db.QueryAllProductCount()
	publicCount := db.QueryAllProductByStateCount(1)
	noPublicCount := db.QueryAllProductByStateCount(0)
	return gin.H{
		"total":         total,
		"publicCount":   publicCount,
		"noPublicCount": noPublicCount,
	}
}

//设置状态
func SetProductStateSerivce(id, state int) error {
	product, err := db.QueryProductByID(id)
	if err != nil {
		return err
	}
	productLock.Lock()
	defer productLock.Unlock()

	//设置
	product.State = state
	t := time.Now()
	if state == 1 {
		product.PublishTime = &t
	}

	return db.UpdateProduct(product)
}

//删除项目
func DeleteProductSerivce(id int) error {
	var productIds []int = make([]int, 0)
	productIds = append(productIds, id)
	devices := db.QueryDevicetByProductIds(productIds)
	if len(devices) > 0 {
		return errors.New("产品下包含设备")
	}
	product, err := db.QueryProductByID(id)
	if err != nil {
		return err
	}
	productLock.Lock()
	defer productLock.Unlock()

	//设置
	product.DelFlag = 1
	return db.UpdateProduct(product)
}

//根据gateWayId查询产品
func QueryProductsByGateWayId(gateWayId int) []*db.Device {
	var productIds []int = make([]int, 0)
	productIds = append(productIds, gateWayId)
	devices := db.QueryDevicetByProductIds(productIds)
	return devices
}
