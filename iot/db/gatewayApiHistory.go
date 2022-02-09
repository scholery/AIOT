package db

func SaveGatewayApiHistory(history GatewayApiHistory) error {
	var item GatewayApiHistory
	err := webOrm.QueryTable("gateway_api_history").Filter("gateway_id", history.GatewayId).Filter("device_id", history.DeviceId).Filter("name", history.Name).One(&item)
	if err != nil {
		_, err := webOrm.Insert(&history)
		if err != nil {
			logger.Errorln(err)
		}
		return err
	} else {
		history.Id = item.Id
		_, err := webOrm.Update(&history)
		return err
	}
}

func QueryAllGatewayApiHistory() (int64, []*GatewayApiHistory) {
	var childrenItem []*GatewayApiHistory
	qs := webOrm.QueryTable("gateway_api_history")
	num, err := qs.All(&childrenItem)
	if err != nil {
		logger.Errorln(err)
	}
	return num, childrenItem
}
