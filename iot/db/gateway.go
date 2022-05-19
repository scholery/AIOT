package db

import "strconv"

func InsertGateway(gateway Gateway) (int64, error) {
	id, err := webOrm.Insert(&gateway)
	if err != nil {
		logger.Errorln(err)
		id = -1
	}
	return id, err
}

func QueryAllGateways() (int64, []*Gateway) {
	var childrenItem []*Gateway
	qs := webOrm.QueryTable("gateway")
	num, err := qs.Filter("del_flag", 0).All(&childrenItem)
	if err != nil {
		logger.Errorln(err)
	}
	return num, childrenItem
}

func CheckIPPort(gatewayIp string, gatewayPort int, gatewayId int) (int64, error) {
	var childrenItem []*Gateway
	qs := webOrm.QueryTable("gateway")
	querySelector := qs.Filter("del_flag", 0).Filter("ip", gatewayIp).Filter("port", gatewayPort).Filter("protocol__in", "http_server", "WebSocket")
	if gatewayId != 0 {
		querySelector = querySelector.Exclude("id", gatewayId)
	}
	num, err := querySelector.All(&childrenItem)
	if err != nil {
		logger.Errorln(err)
	}
	return num, err
}

func CheckSign(gatewaySign string, gatewayId int) (int64, error) {
	var childrenItem []*Gateway
	qs := webOrm.QueryTable("gateway")
	querySelector := qs.Filter("del_flag", 0).Filter("sign", gatewaySign)
	if gatewayId != 0 {
		querySelector = querySelector.Exclude("id", gatewayId)
	}
	num, err := querySelector.All(&childrenItem)
	if err != nil {
		logger.Errorln(err)
	}
	return num, err
}

func UpdateGateway(gateway *Gateway) (int64, error) {
	webOrm.Begin()
	num, err := webOrm.Update(gateway, "name", "sign", "protocol", "ip", "port", "authInfo", "routers", "collectType", "collectPeriod", "cron", "modbusConfig", "describe")
	if err != nil {
		webOrm.Rollback()
	}
	webOrm.Commit()
	return num, err
}

func UpdateGatewayStatus(gateway *Gateway) (int64, error) {
	webOrm.Begin()
	num, err := webOrm.Update(gateway, "status")
	if err != nil {
		webOrm.Rollback()
	}
	webOrm.Commit()
	return num, err
}

func DeleteGateway(gatewayId int) error {
	// _, err := webOrm.Delete(&Gateway{Id: gatewayId})
	_, err := webOrm.Update(&Gateway{Id: gatewayId, DelFlag: 1}, "del_flag")
	if err != nil {
		webOrm.Rollback()
	}
	return err
}

func QueryGatewayByID(gatewayId int) (*Gateway, error) {
	var gateway Gateway
	err := webOrm.QueryTable("gateway").Filter("id", gatewayId).Filter("del_flag", 0).One(&gateway)
	if err != nil {
		return nil, err
	}
	return &gateway, nil
}

func QueryGatewaysByPage(offset, limit int, gatewayName string, gatewayProtocol string, gatewayStatus string) []*Gateway {
	var childrenItem []*Gateway
	qs := webOrm.QueryTable("gateway")
	querySelector := qs.Filter("del_flag", 0).Filter("name__contains", gatewayName)
	if gatewayProtocol != "" {
		querySelector = querySelector.Filter("protocol", gatewayProtocol)
	}
	if gatewayStatus != "" {
		statusInt, err := strconv.Atoi(gatewayStatus)
		if err != nil {
			return nil
		}
		querySelector = querySelector.Filter("status", statusInt)
	}
	_, err := querySelector.Limit(limit, offset).All(&childrenItem)
	if err != nil {
		logger.Errorln(err)
	}
	return childrenItem
}

func ListGatewaysCount(gatewayName string, gatewayProtocol string, gatewayStatus string) int64 {
	qs := webOrm.QueryTable("gateway")
	querySelector := qs.Filter("del_flag", 0).Filter("name__contains", gatewayName)
	if gatewayProtocol != "" {
		querySelector = querySelector.Filter("protocol", gatewayProtocol)
	}
	if gatewayStatus != "" {
		statusInt, err := strconv.Atoi(gatewayStatus)
		if err != nil {
			return 0
		}
		querySelector = querySelector.Filter("status", statusInt)
	}
	num, err := querySelector.Count()
	if err != nil {
		logger.Errorln(err)
	}
	return num
}

func QueryGatewaysByStatus(gatewayStatus string) []*Gateway {
	var childrenItem []*Gateway
	qs := webOrm.QueryTable("gateway")
	querySelector := qs.Filter("del_flag", 0)
	if gatewayStatus != "" {
		querySelector = querySelector.Filter("status", gatewayStatus)
	}
	_, err := querySelector.All(&childrenItem)
	if err != nil {
		logger.Errorln(err)
	}
	return childrenItem
}
