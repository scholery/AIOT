package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"koudai-box/conf"

	"koudai-box/iot/db"
	"koudai-box/iot/service"
	"koudai-box/iot/web/common"
	"koudai-box/iot/web/dto"

	iot "koudai-box/iot/gateway/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

func RegisterDeviceURL(r *gin.RouterGroup) {
	r.POST("/device/add", addDevice)
	r.POST("/device/update", updateDevice)
	r.POST("/device/data", queryDevice)
	r.GET("/device/info/:id", queryDeviceInfo)

	r.POST("/device/set-activateStatus", setDeviceActivateStatus) //修改设备激活状态
	r.GET("/device/delete/:id", deleteDevice)
	r.DELETE("/device/deletes", deletesDevice)
	r.GET("/device/statistics", statisticsDevice)             //统计设备信息
	r.POST("/device/batchActivates", setDeviceBatchActivates) //批量激活  批量停止
	r.GET("/device/allActivate", allDeviceActivate)           //全部激活
	r.GET("/device/currentSituation/:id", queryDeviceCurrentSituation)

	r.GET("/device/excelTemplate", excelTemplate)
	r.POST("/device/export", export)
	r.POST("/device/importExcel", importExcelHandler)

	r.Any("/device/postOperation/:id/:operationCode", postOperation) //操作下发
	r.Any("/device/sync", syncDevice)
	//操作下发
	r.Any("/device/zerostatus/:id", setZerostatus)
	r.Any("/device/calcpredayavg/:id", CalcPredayAvg)

}

//添加设备
func addDevice(c *gin.Context) {
	var request dto.AddDeviceRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	eventId, err := service.AddDeviceService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("保存成功", eventId))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

// 更新设备
func updateDevice(c *gin.Context) {
	var request dto.UpdateDeviceRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	deviceId, err := service.UpdateDeviceService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("保存成功", gin.H{
			"deviceId": deviceId,
		}))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//查询设备列表
func queryDevice(c *gin.Context) {
	var request dto.QueryDeviceDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	if request.PageNo < 1 {
		request.PageNo = 1
	}
	if request.PageSize < 1 {
		println(request.PageSize)
		request.PageSize = 10
	}

	count, page := service.QueryDeviceSerivce(request)
	c.JSON(http.StatusOK, common.OkPage(request.PageNo, request.PageSize, count, "", page))
}

//查询设备详情
func queryDeviceInfo(c *gin.Context) {
	id := c.Param("id")
	deviceID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	deviceInfo, _ := service.QueryDeviceInfoByID(deviceID)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("查询错误", err))
	}
	c.JSON(http.StatusOK, common.Ok("获取成功", deviceInfo))
}

//统计产品信息
func statisticsDevice(c *gin.Context) {
	data := service.StatisticsDeviceSerivce()
	c.JSON(http.StatusOK, common.Ok("获取成功", data))
}

//删设备
func deleteDevice(c *gin.Context) {
	id := c.Param("id")
	deviceID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}

	err = service.DeleteDeviceSerivce(deviceID)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Ok("删除成功", nil))
}

func deletesDevice(c *gin.Context) {
	var request dto.DeleteDeviceRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	for _, v := range request.Ids {
		err := service.DeleteDeviceSerivce(v)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
	}
	c.JSON(http.StatusOK, common.Ok("删除成功", nil))
}

func setDeviceActivateStatus(c *gin.Context) {
	var request dto.SetDeviceActivateRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	deviceInfo, err := service.QueryDeviceById(request.Id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("查询设备失败", nil))
		return
	}
	if deviceInfo == nil {
		c.JSON(http.StatusOK, common.Error("设备不存在", nil))
		return
	}
	if deviceInfo.ActivateStatus == *request.ActivateStatus {
		c.JSON(http.StatusOK, common.Ok("设置完成", nil))
		return
	}
	err = service.SetDeviceStatusSerivce(request.Id, *request.ActivateStatus)
	if err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	ok := false
	if *request.ActivateStatus == 0 {
		ok = iot.StopDevice(request.Id)
	} else {
		ok = iot.StartDevice(request.Id)
	}
	if !ok {
		logrus.Errorf("change device[%d]'s state error", request.Id)
		return
	}

	c.JSON(http.StatusOK, common.Ok("设置成功", gin.H{
		"deviceId": request.Id,
	}))
}

func setDeviceBatchActivates(c *gin.Context) {
	var requests dto.SetBatchDeviceActivateRequest
	if err := c.ShouldBind(&requests); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	if *requests.ActivateStatus != 0 && *requests.ActivateStatus != 1 {
		c.JSON(http.StatusOK, common.Error("状态异常", nil))
		return
	}

	for _, deviceId := range requests.Ids {
		deviceInfo, err := service.QueryDeviceById(deviceId)
		if err != nil {
			c.JSON(http.StatusOK, common.Error("查询设备失败", nil))
			return
		}
		if deviceInfo == nil {
			c.JSON(http.StatusOK, common.Error("设备不存在", nil))
			return
		}
		if deviceInfo.ActivateStatus == *requests.ActivateStatus {
			c.JSON(http.StatusOK, common.Ok("设置完成", nil))
			return
		}
		err = service.SetDeviceStatusSerivce(deviceId, *requests.ActivateStatus)
		if err != nil {
			c.JSON(http.StatusOK, common.Error(err.Error(), nil))
			return
		}
		ok := false
		if *requests.ActivateStatus == 0 {
			ok = iot.StopDevice(deviceId)
		} else {
			ok = iot.StartDevice(deviceId)
		}
		if !ok {
			logrus.Errorf("change device[%d]'s state error", deviceId)
			continue
		}

	}
	c.JSON(http.StatusOK, common.Ok("设置成功", nil))
}

func allDeviceActivate(c *gin.Context) {
	unActivateDeviceIds := service.QueryDeviceIdsByActivateStatus(0)
	for _, deviceId := range unActivateDeviceIds {

		err := service.SetDeviceStatusSerivce(deviceId, 1)
		if err != nil {
			c.JSON(http.StatusOK, common.Error(err.Error(), nil))
			return
		}
		ok := false
		ok = iot.StartDevice(deviceId)
		if !ok {
			logrus.Errorf("change device[%d]'s state error", deviceId)
			continue
		}

	}
	c.JSON(http.StatusOK, common.Ok("设置成功", nil))
}

func queryDeviceCurrentSituation(c *gin.Context) {
	deviceID := c.Param("id")
	deviceProps, _ := service.GetLastestProperty(deviceID)
	c.JSON(http.StatusOK, common.Ok("获取成功", deviceProps))
}

func excelTemplate(c *gin.Context) {
	file := "template/" + "设备" + ".xlsx"
	if exists, _ := common.FileExist(file); !exists {
		c.JSON(STATUS_OK, common.Error("模板文件未找到", nil))
		return
	}
	f, err := os.Open(file)
	if err != nil {
		logrus.Error(err)
		c.JSON(STATUS_OK, common.Error("打开模板文件失败", nil))
		return
	}
	defer f.Close()

	fileContent := bytes.NewBuffer(nil)
	_, err = io.Copy(fileContent, f)
	if err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	} else {
		c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape("设备")+".xlsx")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileContent.Bytes())
	}
}

func importExcelHandler(c *gin.Context) {
	cameraType, xlsx, err := parseImportExcelParams(c)
	if err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	importResult := importExcel(cameraType, xlsx)
	c.JSON(http.StatusOK, common.Ok(importResult.String(), importResult))
}

func parseImportExcelParams(c *gin.Context) (string, *excelize.File, error) {
	datas, _ := c.MultipartForm()

	headers, ok := datas.File["file"]
	if !ok {
		return "", nil, errors.New("请选择要上传的文件")
	}
	if !strings.HasSuffix(headers[0].Filename, ".xlsx") {
		return "", nil, errors.New("请上传xlsx文件")
	}
	fi, err := headers[0].Open()
	if err != nil {
		logrus.Error(err)
		return "", nil, errors.New("xlsx文件无法打开")
	}
	defer fi.Close()
	excel, err := excelize.OpenReader(fi)
	if err != nil {
		logrus.Error(err)
		return "", nil, errors.New("上传的文件不是一个合法的excel文件")
	}
	return "", excel, nil
}

func importExcel(cameraType string, xlsx *excelize.File) common.ImportResult {
	sheetName := xlsx.GetSheetMap()[1]
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		logrus.Error(err)
		return common.ImportResult{}
	}
	return importDevice(rows[1:])
}

func importDevice(rows [][]string) common.ImportResult {
	result := common.ImportResult{}
	for index, row := range rows {
		addDevice := dto.AddDeviceRequest{}
		for index, value := range row {
			if index == 0 {
				addDevice.Code = value
			} else if index == 1 {
				addDevice.Name = value
			} else if index == 2 {
				procuct, err := db.QueryProductByName(value)
				if err != nil {
					addDevice.ProductId = 0
					logrus.Error(err)
				} else {
					addDevice.ProductId = procuct.Id
				}
			} else if index == 3 {
				addDevice.Desc = value
			}
		}
		_, err := service.AddDeviceService(addDevice)
		if err != nil {
			logrus.Error("行 ", index, " 数据错误，", err)
			result.FailCount++
			result.FailedRows = append(result.FailedRows, index+2)
			continue
		}
		result.SuccessCount++
	}
	return result
}

func export(c *gin.Context) {
	var request dto.DeleteDeviceRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	devices := db.QueryDevicesByIds(request.Ids)
	timeFormat := time.Now().Format("2006-01-02")
	excelFilePath := conf.GetConf().TempPath + "/" + timeFormat + ".xlsx"
	_, err := os.Create(excelFilePath)
	if err != nil {
		logrus.Errorln(err)
		c.JSON(STATUS_OK, common.Error("创建excel文件失败", nil))
		return
	}
	f := excelize.NewFile()
	sheetName := "Sheet1"
	index := f.NewSheet(sheetName)
	f.SetCellValue(sheetName, "A1", "设备标识")
	f.SetCellValue(sheetName, "B1", "设备名称")
	f.SetCellValue(sheetName, "C1", "产品名称")
	f.SetCellValue(sheetName, "D1", "产品标识")
	f.SetCellValue(sheetName, "E1", "描述")
	f.SetCellValue(sheetName, "F1", "注册时间")
	for index, v := range devices {
		col := index + 2
		f.SetCellValue(sheetName, "A"+strconv.Itoa(col), v[0])
		f.SetCellValue(sheetName, "B"+strconv.Itoa(col), v[1])
		f.SetCellValue(sheetName, "C"+strconv.Itoa(col), v[2])
		f.SetCellValue(sheetName, "D"+strconv.Itoa(col), v[3])
		f.SetCellValue(sheetName, "E"+strconv.Itoa(col), v[4])
		f.SetCellValue(sheetName, "F"+strconv.Itoa(col), v[5])
	}
	f.SetActiveSheet(index)
	if err := f.SaveAs(excelFilePath); err != nil {
		fmt.Println(err)
	}

	tmpf, err := os.Open(excelFilePath)
	if err != nil {
		logrus.Error(err)
		c.JSON(STATUS_OK, common.Error("打开excel文件失败", nil))
		return
	}
	defer tmpf.Close()
	fileContent := bytes.NewBuffer(nil)
	_, _ = io.Copy(fileContent, tmpf)
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(timeFormat)+".xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileContent.Bytes())
}

func postOperation(c *gin.Context) {
	id := c.Param("id")
	operationCode := c.Param("operationCode")
	deviceId, err := strconv.Atoi(id)
	if err != nil || len(operationCode) == 0 {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	err = iot.PostOperation(deviceId, operationCode)
	if nil != err {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, common.Ok("操作成功", nil))
}

func syncDevice(c *gin.Context) {
	err := service.PushAllDevices()
	if nil != err {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, common.Ok("同步成功", nil))
}

func setZerostatus(c *gin.Context) {
	id := c.Param("id")
	deviceId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	iot.SetZeroStatus(deviceId)
	c.JSON(http.StatusOK, common.Ok("操作成功", nil))
}

func CalcPredayAvg(c *gin.Context) {
	id := c.Param("id")
	deviceId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	preday := service.CalcPredayAvg(deviceId)
	iot.SetPredayStatus(deviceId, preday)
	c.JSON(http.StatusOK, common.Ok("计算成功", nil))
}
