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

func RegisterEventURL(r *gin.RouterGroup) {
	r.POST("/event/add", saveEvent)
	r.POST("/event/update", updateEvent)
	r.POST("/event/delete", deleteEvent)
	r.POST("/event/data", queryEventData)

	r.GET("/event/info/:id", queryEventInfo)

	r.GET("/event/statisticsEvent", StatisticsEvent)
}

//保存事件
func saveEvent(c *gin.Context) {
	var request dto.SaveEventRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	eventId, err := service.AddEventService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("保存成功", eventId))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//更新事件
func updateEvent(c *gin.Context) {
	var request dto.UpdateEventRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.UpdateEventService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("保存成功", nil))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//删除事件
func deleteEvent(c *gin.Context) {
	var request dto.DeleteEventRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.DeleteEventService(request.Ids)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("删除成功", nil))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//查询事件
func queryEventData(c *gin.Context) {
	var request dto.QueryEventDataRequest
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

	count, page := service.QueryEventSerivce(request)
	c.JSON(http.StatusOK, common.OkPage(request.PageNo, request.PageSize, count, "", page))
}

//事件详情
func queryEventInfo(c *gin.Context) {
	ID := c.Param("id")
	eventID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	event := service.QueryEventByIDService(eventID)
	c.JSON(http.StatusOK, common.Ok("", event))
}

//事件统计
func StatisticsEvent(c *gin.Context) {
	totalEvents := db.CountTotalEvents()
	todayEvents := db.CountTodayEvents()
	mostEventDeviceName := db.CountTodayMostEventDeviceName()
	todayMostEventName := db.CountTodayMostEventName()
	data := gin.H{
		"totalEvents":         totalEvents,
		"todayEvents":         todayEvents,
		"mostEventDeviceName": mostEventDeviceName,
		"todayMostEventName":  todayMostEventName,
	}
	c.JSON(http.StatusOK, common.Ok("", data))
}
