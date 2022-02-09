package sysconfig

type StorageConfig struct {
	EventStorageMaxPercentage int `json:"eventStorageMaxPercentage"` // 0.5-0.95
	EventMaxHours             int `json:"eventMaxHours"`             //<=0为无限制
	AlarmMaxCount             int `json:"alarmMaxCount"`             //要求100-20000之间， 不得超过20000
	IotAlarmMaxCount          int `json:"iotAlarmMaxCount"`          //要求100-20000之间， 不得超过20000
	IotEventMaxCount          int `json:"iotEventMaxCount"`          //要求100-20000之间， 不得超过20000
	IotPropsMaxCount          int `json:"iotPropsMaxCount"`          //要求100-20000之间， 不得超过20000
}
