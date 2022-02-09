package dto

type SaveEventRequest struct {
	Sign       string `json:"sign" form:"sign" binding:"required"`
	Title      string `json:"title" form:"title" binding:"required"`
	Type       string `json:"type" form:"type" binding:"required"` //event、event
	Level      string `json:"level" form:"level" binding:"required"`
	DeviceId   string `json:"deviceId" form:"deviceId" binding:"required"`
	ProductId  int    `json:"productId" form:"productId" binding:"required"`
	CreateTime string `json:"createTime" form:"createTime" binding:"required"`
}

type UpdateEventRequest struct {
	Id         int    `json:"id" form:"id" binding:"required"`
	Sign       string `json:"sign" form:"sign" binding:"required"`
	Title      string `json:"title" form:"title" binding:"required"`
	Type       string `json:"type" form:"type" binding:"required"` //event、event
	Level      string `json:"level" form:"level" binding:"required"`
	DeviceId   string `json:"deviceId" form:"deviceId" binding:"required"`
	ProductId  int    `json:"productId" form:"productId" binding:"required"`
	CreateTime string `json:"createTime" form:"createTime" binding:"required"`
}

type DeleteEventRequest struct {
	Ids []int `json:"eventIds" form:"eventIds" binding:"required"`
}

type QueryEventDataRequest struct {
	Search    string `json:"search" form:"search"`
	DeviceId  string `json:"deviceId" form:"deviceId"`
	Level     string `json:"level" form:"level"`
	StartTime string `json:"startTime" form:"startTime"`
	EndTime   string `json:"endTime" form:"endTime"`
	PageNo    int    `json:"pageNo" form:"pageNo"`
	PageSize  int    `json:"PageSize" form:"PageSize"`
}
