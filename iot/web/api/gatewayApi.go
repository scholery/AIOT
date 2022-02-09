package api

import (
	"errors"
	"net/http"
	"strconv"

	iot "koudai-box/iot/gateway/service"
	"koudai-box/iot/service"
	"koudai-box/iot/web/common"
	"koudai-box/iot/web/dto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func RegisterGatewayURL(r *gin.RouterGroup) {
	r.POST("/gateway/add", saveGateway)
	r.POST("/gateway/update", updateGateway)
	r.POST("/gateway/updateStatus", updateStatus)
	r.POST("/gateway/delete", deleteGateway)
	r.POST("/gateway/data", queryGatewayData)
	r.GET("/gateway/info/:id", queryGatewayInfo)
}

var STATUS_OK = http.StatusOK

func saveGateway(c *gin.Context) {
	var request dto.SaveGatewayRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	gatewayId, err := service.AddGatewayService(request)
	if err == nil {
		c.JSON(STATUS_OK, common.Ok("保存成功", gatewayId))
	} else {
		c.JSON(STATUS_OK, common.Error(err.Error(), nil))
	}
}

func updateGateway(c *gin.Context) {
	var request dto.SaveGatewayRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	err := service.UpdateGatewayService(request)
	if err == nil {
		c.JSON(STATUS_OK, common.Ok("保存成功", nil))
	} else {
		c.JSON(STATUS_OK, common.Error(err.Error(), nil))
	}
}

func updateStatus(c *gin.Context) {
	var request dto.UpdateGatewayStatusRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	err := service.UpdateGatewayStatusService(request)
	if err != nil {
		c.JSON(STATUS_OK, common.Error(err.Error(), nil))
		return
	}
	ok := false
	if request.Status == 0 {
		ok = iot.StopGateway(request.GatewayId)
		SetProductStateProcessBatchByGTId(request.GatewayId, 0)
	} else {
		ok = iot.StartGateway(request.GatewayId)
		SetProductStateProcessBatchByGTId(request.GatewayId, 1)
	}
	if !ok {
		logrus.Errorf("change gateway[%d]'s status error", request.GatewayId)
		c.JSON(http.StatusOK, common.Error("stop or start error", nil))
		return
	}
	c.JSON(http.StatusOK, common.Ok("修改状态成功", gin.H{
		"gatewayId": request.GatewayId,
	}))
}

func SetProductStateProcessBatchByGTId(gateWayId int, state int) {
	devices := service.QueryProductsByGateWayId(gateWayId)
	for _, device := range devices {
		err := SetProductStateProcess(device.Id, state)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
	}
}

func SetProductStateProcess(id, state int) error {
	err := service.SetProductStateSerivce(id, state)
	if err != nil {
		return err
	}
	//设置
	ok := false
	if state == 0 {
		ok = iot.StopProduct(id)
	} else {
		ok = iot.StartProduct(id)
	}
	if !ok {
		logrus.Errorf("change product[%d]'s state error", id)
		return errors.New("设置产品状态失败")
	}
	return nil
}

func deleteGateway(c *gin.Context) {
	var request dto.DeleteGatewayRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	ids := request.GatewayIds
	if len(ids) == 0 {
		ids = []int{request.GatewayId}
	}
	err := service.DeleteGatewayService(ids)
	if err == nil {
		c.JSON(STATUS_OK, common.Ok("删除成功", nil))
	} else {
		c.JSON(STATUS_OK, common.Error(err.Error(), nil))
	}
}

func queryGatewayData(c *gin.Context) {
	var request dto.QueryGatewayDataRequest
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
	count, page := service.QueryGatewaySerivce(request)
	c.JSON(STATUS_OK, common.OkPage(request.PageNo, request.PageSize, count, "", page))
}

func queryGatewayInfo(c *gin.Context) {
	ID := c.Param("id")
	gatewayID, err := strconv.Atoi(ID)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	gateway := service.QueryGatewayByIDService(gatewayID)
	c.JSON(http.StatusOK, common.Ok("", gateway))
}
