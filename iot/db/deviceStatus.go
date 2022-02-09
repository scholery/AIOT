package db

func SaveDeviceStatus(status DeviceStatus) error {
	var item DeviceStatus
	err := webOrm.QueryTable("device_status").Filter("device_id", status.DeviceId).One(&item)
	if err != nil {
		_, err := webOrm.Insert(&status)
		if err != nil {
			logger.Errorln(err)
		}
		return err
	} else {
		cols := []string{"status", "timestamp"}
		if status.LastStatus != "" {
			cols = append(cols, "last_status")
		}
		if status.PreDayStatus != "" {
			cols = append(cols, "pre_day_status")
		}
		if status.ZeroStatus != "" {
			cols = append(cols, "zero_status")
		}
		_, err = webOrm.Update(&status, cols...)
		return err
	}
}

func QueryAllDeviceStatus() (int64, []*DeviceStatus) {
	var items []*DeviceStatus
	qs := webOrm.QueryTable("device_status")
	num, err := qs.All(&items)
	if err != nil {
		logger.Errorln(err)
	}
	return num, items
}
