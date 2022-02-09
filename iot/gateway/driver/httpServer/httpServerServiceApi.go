package httpServer

import (
	"net/http"

	"koudai-box/iot/gateway/model"
	"koudai-box/iot/web/common"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RegisterURL(r *gin.RouterGroup) {
	r.POST("/gateway/:gatewayKey/props", batchDeviceProps)

	r.POST("/gateway/:gatewayKey/props/:deviceKey", singleDeviceProps)

	r.POST("/gateway/:gatewayKey/events", batchRecieveEvents)

	r.POST("/gateway/:gatewayKey/events/:deviceKey", recieveEvents)
}

//单条推送
func singleDeviceProps(c *gin.Context) {
	gatewayKey := c.Param("gatewayKey")
	deviceKey := c.Param("deviceKey")
	if len(gatewayKey) == 0 || len(deviceKey) == 0 {
		c.JSON(http.StatusOK, common.Error("路径参数错误", nil))
		return
	}
	var request interface{}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	logrus.Debugf("gateway[%s] deviceKey[%s]'s props:", gatewayKey, deviceKey, request)

	var pushMsg model.PushMsg
	pushMsg.Msg = request
	pushMsg.Type = model.Msg_Type_Props
	pushMsg.GatewayKey = gatewayKey
	pushMsg.DeviceKey = deviceKey
	model.PushMsgChan <- pushMsg
	c.JSON(http.StatusOK, common.Ok("消息已接收，正在处理", ""))
}

//多条推送
func batchDeviceProps(c *gin.Context) {
	gatewayKey := c.Param("gatewayKey")
	if len(gatewayKey) == 0 {
		c.JSON(http.StatusOK, common.Error("路径参数错误", nil))
		return
	}
	var request interface{}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	logrus.Debugf("gateway[%s]'s props", gatewayKey, request)
	var pushMsg model.PushMsg
	pushMsg.Msg = request
	pushMsg.Type = model.Msg_Type_Props
	pushMsg.GatewayKey = gatewayKey
	model.PushMsgChan <- pushMsg
	c.JSON(http.StatusOK, common.Ok("消息已接收，正在处理", ""))
}

//单条推送
func recieveEvents(c *gin.Context) {
	gatewayKey := c.Param("gatewayKey")
	deviceKey := c.Param("deviceKey")
	if len(gatewayKey) == 0 || len(deviceKey) == 0 {
		c.JSON(http.StatusOK, common.Error("路径参数错误", nil))
		return
	}
	var request interface{}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	logrus.Debugf("gateway[%s] deviceKey[%s]'s events:", gatewayKey, deviceKey, request)
	var pushMsg model.PushMsg
	pushMsg.Msg = request
	pushMsg.Type = model.Msg_Type_Events
	pushMsg.GatewayKey = gatewayKey
	pushMsg.DeviceKey = deviceKey
	model.PushMsgChan <- pushMsg
	c.JSON(http.StatusOK, common.Ok("消息已接收，正在处理", ""))
}

//多条推送
func batchRecieveEvents(c *gin.Context) {
	gatewayKey := c.Param("gatewayKey")
	if len(gatewayKey) == 0 {
		c.JSON(http.StatusOK, common.Error("路径参数错误", nil))
		return
	}
	var request interface{}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	logrus.Debugf("gateway[%s]'s events", gatewayKey, request)

	var pushMsg model.PushMsg
	pushMsg.Msg = request
	pushMsg.Type = model.Msg_Type_Events
	pushMsg.GatewayKey = gatewayKey
	model.PushMsgChan <- pushMsg
	c.JSON(http.StatusOK, common.Ok("消息已接收，正在处理", ""))
}
