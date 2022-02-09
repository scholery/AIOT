package api

import (
	"net/http"
	"strconv"

	"koudai-box/iot/db"
	"koudai-box/iot/service"
	"koudai-box/iot/web/common"
	"koudai-box/iot/web/dto"

	"github.com/gin-gonic/gin"
)

func RegisterAlarmURL(r *gin.RouterGroup) {
	r.POST("/alarm/add", saveAlarm)
	r.POST("/alarm/update", updateAlarm)
	r.POST("/alarm/delete", deleteAlarm)
	r.POST("/alarm/data", queryAlarmData)

	r.GET("/alarm/info/:id", queryAlarmInfo)

	r.GET("/alarm/statisticsAlarm", StatisticsAlarm)
}

//保存告警
func saveAlarm(c *gin.Context) {
	var request dto.SaveAlarmRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	alarmId, err := service.AddAlarmService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("保存成功", alarmId))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//更新告警
func updateAlarm(c *gin.Context) {
	var request dto.UpdateAlarmRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.UpdateAlarmService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("保存成功", nil))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//删除告警
func deleteAlarm(c *gin.Context) {
	var request dto.DeleteAlarmRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.DeleteAlarmService(request.AlarmIds)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("删除成功", nil))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//查询告警
func queryAlarmData(c *gin.Context) {
	var request dto.QueryAlarmDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	if request.PageNo < 1 {
		request.PageNo = 1
	}
	if request.PageSize < 1 {
		request.PageSize = 10
	}
	count, page := service.QueryAlarmSerivce(request)
	c.JSON(http.StatusOK, common.OkPage(request.PageNo, request.PageSize, count, "", page))
}

//告警详情
func queryAlarmInfo(c *gin.Context) {
	ID := c.Param("id")
	alarmID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	alarm := service.QueryAlarmByIDService(alarmID)
	c.JSON(http.StatusOK, common.Ok("", alarm))
}

//告警统计
func StatisticsAlarm(c *gin.Context) {
	totalAlarms := db.CountTotalAlarms()
	todayAlarms := db.CountTodayAlarms()
	mostAlarmDeviceName := db.CountTodayMostAlarmDeviceName()
	todayMostAlarmName := db.CountTodayMostAlarmName()
	data := gin.H{
		"totalAlarms":         totalAlarms,
		"todayAlarms":         todayAlarms,
		"mostAlarmDeviceName": mostAlarmDeviceName,
		"todayMostAlarmName":  todayMostAlarmName,
	}
	c.JSON(http.StatusOK, common.Ok("", data))
}
