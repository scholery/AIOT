package service

import (
	"encoding/json"
	"errors"
	"sync"

	"koudai-box/cache"

	"koudai-box/iot/db"
	"koudai-box/iot/gateway/model"
	"koudai-box/iot/gateway/utils"
	"koudai-box/iot/web/common"
	"koudai-box/iot/web/dto"

	"github.com/sirupsen/logrus"
)

const (
	GATEWAY_CACHE_KEY string = "gatewayCache"
)

var gatewayLock = sync.Mutex{}

func AddGatewayService(request dto.SaveGatewayRequest) (int64, error) {
	num, _ := db.CheckIPPort(request.Ip, request.Port, 0)
	if num != 0 {
		return 0, errors.New("端口已存在")
	}
	num, _ = db.CheckSign(request.Sign, 0)
	if num != 0 {
		return 0, errors.New("标识已存在")
	}
	gatewayLock.Lock()
	defer gatewayLock.Unlock()
	// var autoIncGatewayId int64
	gateway := db.Gateway{
		Name:          request.GatewayName,
		Sign:          request.Sign,
		Status:        request.Status,
		Protocol:      request.Protocol,
		Ip:            request.Ip,
		Port:          request.Port,
		AuthInfo:      getAuthInfoStr(request.AuthInfo),
		Routers:       getRouterStr(request.Routers),
		CollectType:   request.CollectType,
		CollectPeriod: request.CollectPeriod,
		Cron:          request.Cron,
		ModbusConfig:  utils.ToString(request.ModbusConfig),
		Describe:      request.Describe,
	}
	autoIncGatewayId, err := db.InsertGateway(gateway)
	ClearGatewayCache()
	return autoIncGatewayId, err
}

func getAuthInfoStr(authInfos []dto.AuthInfoItem) string {
	authInfoStr, err := json.Marshal(authInfos)
	if err != nil {
		return ""
	}
	return string(authInfoStr)
}

func getRouterStr(routerInfos []dto.RouterItem) string {
	routerInfoStr, err := json.Marshal(routerInfos)
	if err != nil {
		return ""
	}
	return string(routerInfoStr)
}

func UpdateGatewayService(request dto.SaveGatewayRequest) error {
	gateway := GetGatewayFromCache(request.GatewayId)
	if gateway == nil {
		return errors.New("网关不存在")
	}
	if request.Protocol == "http_server" {
		num, _ := db.CheckIPPort(request.Ip, request.Port, request.GatewayId)
		if num != 0 {
			return errors.New("端口已存在")
		}
	}
	num, _ := db.CheckSign(request.Sign, request.GatewayId)
	if num != 0 {
		return errors.New("标识已存在")
	}
	gatewayLock.Lock()
	defer gatewayLock.Unlock()

	db.UpdateGateway(&db.Gateway{
		Id:            request.GatewayId,
		Name:          request.GatewayName,
		Sign:          request.Sign,
		Protocol:      request.Protocol,
		Ip:            request.Ip,
		Port:          request.Port,
		AuthInfo:      getAuthInfoStr(request.AuthInfo),
		Routers:       getRouterStr(request.Routers),
		CollectType:   request.CollectType,
		CollectPeriod: request.CollectPeriod,
		Cron:          request.Cron,
		ModbusConfig:  utils.ToString(request.ModbusConfig),
		Describe:      request.Describe,
	})
	ClearGatewayCache()
	return nil
}

func UpdateGatewayStatusService(request dto.UpdateGatewayStatusRequest) error {
	gateway := GetGatewayFromCache(request.GatewayId)
	if gateway == nil {
		return errors.New("网关不存在")
	}
	gatewayLock.Lock()
	defer gatewayLock.Unlock()

	db.UpdateGatewayStatus(&db.Gateway{
		Id:     request.GatewayId,
		Status: request.Status,
	})
	ClearGatewayCache()
	return nil
}

func DeleteGatewayService(ids []int) error {
	gatewayLock.Lock()
	defer gatewayLock.Unlock()

	deleteGatewayIdList := make([]int, 0)

	for _, c := range ids {
		products, err := db.QueryProductByGatewayID(c)
		if err != nil {
			logrus.Error(err)
			continue
		}
		if len(products) > 0 {
			logrus.Warning("网关下包含产品")
			continue
		}
		err = deleteOneGateway(c)
		if err != nil {
			logrus.Error(err)
			continue
		}
		deleteGatewayIdList = append(deleteGatewayIdList, c)
	}
	ClearGatewayCache()
	return nil
}

func QueryGatewaySerivce(request dto.QueryGatewayDataRequest) (int64, []*dto.GatewayItem) {
	offset, limit := common.Page2Offset(request.PageNo, request.PageSize)
	gateways := db.QueryGatewaysByPage(offset, limit, request.GatewayName, request.GatewayProtocol, request.GatewayStatus)
	totalSize := db.ListGatewaysCount(request.GatewayName, request.GatewayProtocol, request.GatewayStatus)
	var gatewayItems []*dto.GatewayItem
	for _, gateway := range gateways {
		gatewayItem := dto.GatewayItem{
			GatewayId:     gateway.Id,
			GatewayName:   gateway.Name,
			Sign:          gateway.Sign,
			Status:        gateway.Status,
			Protocol:      gateway.Protocol,
			ProtocolName:  GetDictName("gateway_protocol_type", gateway.Protocol),
			Ip:            gateway.Ip,
			Port:          gateway.Port,
			AuthInfo:      getAuthInfo(gateway.AuthInfo),
			Routers:       GetRouters(gateway.Routers),
			CollectType:   gateway.CollectType,
			CollectPeriod: gateway.CollectPeriod,
			Cron:          gateway.Cron,
			ModbusConfig:  getModbusConfig(gateway.ModbusConfig),
			Describe:      gateway.Describe,
		}
		gatewayItems = append(gatewayItems, &gatewayItem)
	}
	return totalSize, gatewayItems
}

func getModbusConfig(config string) model.ModbusConfig {
	str := []byte(config)
	var modbusConfig model.ModbusConfig
	err := json.Unmarshal(str, &modbusConfig)
	if err != nil {
		return model.ModbusConfig{}
	}
	return modbusConfig
}

func getAuthInfo(authInfoStr string) []dto.AuthInfoItem {
	str := []byte(authInfoStr)
	authInfoItems := []dto.AuthInfoItem{}
	err := json.Unmarshal(str, &authInfoItems)
	if err != nil {
		return nil
	}
	return authInfoItems
}

func GetRouters(routerStr string) []dto.RouterItem {
	str := []byte(routerStr)
	routerItems := []dto.RouterItem{}
	err := json.Unmarshal(str, &routerItems)
	if err != nil {
		return nil
	}
	return routerItems
}

func QueryGatewayByIDService(gatewayID int) *dto.GatewayItem {
	gateway := GetGatewayFromCache(gatewayID)
	return gateway
}

func deleteOneGateway(gatewayId int) error {
	gateway := GetGatewayFromCache(gatewayId)
	if gateway == nil {
		return errors.New("网关不存在")
	}
	err := db.DeleteGateway(gatewayId)
	if err != nil {
		return errors.New("删除失败")
	}
	return nil

}

func ClearGatewayCache() {
	cache.Delete(GATEWAY_CACHE_KEY)
}

func GetGatewayFromCache(gatewayId int) *dto.GatewayItem {
	c := GetGatewayCache()[gatewayId]
	return c
}

func GetGatewayCache() map[int]*dto.GatewayItem {
	m, err := cache.Get(GATEWAY_CACHE_KEY)
	if err != nil {
		InitGatewayCache()
		m, _ = cache.Get(GATEWAY_CACHE_KEY)
		if m == nil {
			return make(map[int]*dto.GatewayItem)
		} else {
			return m.(map[int]*dto.GatewayItem)
		}
	}
	return m.(map[int]*dto.GatewayItem)
}

func InitGatewayCache() {
	gateways := ListAllGateway()
	gatewayMap := make(map[int]*dto.GatewayItem)
	for _, c := range gateways {
		gatewayMap[c.GatewayId] = c
	}
	err := SetGatewayCache(gatewayMap)
	if err != nil {
		logrus.Errorln("缓存中网关数据失败:", err)
	}
}

func SetGatewayCache(value map[int]*dto.GatewayItem) error {
	return cache.SetWithNoExpire(GATEWAY_CACHE_KEY, value)
}

func ListAllGateway() []*dto.GatewayItem {
	gatewayItems := make([]*dto.GatewayItem, 0)
	_, gateways := db.QueryAllGateways()
	for _, gateway := range gateways {
		gatewayItem := dto.GatewayItem{
			GatewayId:     gateway.Id,
			GatewayName:   gateway.Name,
			Sign:          gateway.Sign,
			Status:        gateway.Status,
			Protocol:      gateway.Protocol,
			Ip:            gateway.Ip,
			Port:          gateway.Port,
			AuthInfo:      getAuthInfo(gateway.AuthInfo),
			Routers:       GetRouters(gateway.Routers),
			Describe:      gateway.Describe,
			CollectType:   gateway.CollectType,
			CollectPeriod: gateway.CollectPeriod,
			Cron:          gateway.Cron,
			ModbusConfig:  getModbusConfig(gateway.ModbusConfig),
		}
		gatewayItems = append(gatewayItems, &gatewayItem)
	}
	return gatewayItems
}
