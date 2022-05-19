package dto

type SaveAlarmRequest struct {
	Code       string `json:"code" form:"code" binding:"required"`
	Title      string `json:"title" form:"title" binding:"required"`
	Type       string `json:"type" form:"type" binding:"required"` //event、alarm
	Level      string `json:"level" form:"level" binding:"required"`
	DeviceId   string `json:"deviceId" form:"deviceId" binding:"required"`
	ProductId  int    `json:"productId" form:"productId" binding:"required"`
	CreateTime string `json:"createTime" form:"createTime" binding:"required"`
}

type UpdateAlarmRequest struct {
	Id         int    `json:"id" form:"id" binding:"required"`
	Code       string `json:"code" form:"code" binding:"required"`
	Title      string `json:"title" form:"title" binding:"required"`
	Type       string `json:"type" form:"type" binding:"required"` //event、alarm
	Level      string `json:"level" form:"level" binding:"required"`
	DeviceId   string `json:"deviceId" form:"deviceId" binding:"required"`
	CreateTime string `json:"createTime" form:"createTime" binding:"required"`
}

type DeleteAlarmRequest struct {
	AlarmIds []string `json:"alarmIds" form:"alarmIds" binding:"required"`
}

type QueryAlarmDataRequest struct {
	Search      string `json:"search" form:"search"`
	ProductId   string `json:"productId" form:"productId"`
	ProductName string `json:"productName" form:"productName"`
	DeviceId    string `json:"deviceId" form:"deviceId"`
	DeviceName  string `json:"deviceName" form:"deviceName"`
	Level       string `json:"level" form:"level"`
	Type        string `json:"type" form:"type"`
	StartTime   string `json:"startTime" form:"startTime"`
	EndTime     string `json:"endTime" form:"endTime"`
	PageNo      int    `json:"pageNo" form:"pageNo"`
	PageSize    int    `json:"pageSize" form:"pageSize"`
}
