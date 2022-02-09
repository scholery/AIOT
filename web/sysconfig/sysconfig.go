package sysconfig

import "koudai-box/conf"

func GetStorageConfig() *StorageConfig {
	conf := conf.GetConf()
	//eventStorageMaxPercentage, _ := strconv.Atoi(GetSystemConfig(EVENT_STORAGE_MAX_PERCENTAGE, DEFAULT_EVENT_STORAGE_MAX_PERCENTAGE))
	storageConfig := &StorageConfig{
		AlarmMaxCount:    5000,
		EventMaxHours:    5000,
		IotAlarmMaxCount: conf.IotAlarmMaxCount,
		IotEventMaxCount: conf.IotEventMaxCount,
		IotPropsMaxCount: conf.IotPropsMaxCount,
		//EventStorageMaxPercentage: eventStorageMaxPercentage,
	}
	return storageConfig
}
