package api

import (
	"net/http"
	"path"
	"strconv"

	"koudai-box/conf"

	iot "koudai-box/iot/gateway/service"
	"koudai-box/iot/service"
	"koudai-box/iot/web/common"
	"koudai-box/iot/web/dto"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func RegisterProductURL(r *gin.RouterGroup) {
	//产品
	r.POST("/product/add", addProduct)            //添加产品
	r.POST("/product/update", updateProduct)      //更新产品
	r.POST("/product/data", queryProduct)         //获取产品列表
	r.POST("/product/set-state", setProductState) //设置产品状态
	r.GET("/product/info/:id", queryProductInfo)  //查询产品信息

	r.GET("/product/delete/:id", deleteProduct)     //删除产品
	r.GET("/product/statistics", statisticsProduct) //统计产品信息

	//产品物模型
	r.POST("/product/add-item", addProductItem)                      //添加产品物模型
	r.POST("/product/update-item", updateProductItem)                //添加产品物模型
	r.GET("/product/query-item/:productId", queryProductItem)        //查询所有产品物模型
	r.GET("/product/delete-item/:productId/:key", deleteProductItem) //删除产品物模型
	r.GET("/product/detail-item/:productId/:key", detailProductItem) //查看产品物模型

	//产品操作
	r.POST("/product/add-operation", addProductOperation)                      //添加产品操作
	r.POST("/product/update-operation", updateProductOperation)                //添加产品操作
	r.GET("/product/query-operation/:productId", queryProductOperation)        //查询所有产品操作
	r.GET("/product/delete-operation/:productId/:key", deleteProductOperation) //删除产品操作
	r.GET("/product/detail-operation/:productId/:key", detailProductOperation) //查看产品操作

	//产品事件
	r.POST("/product/add-event", addProductEvent)                      //添加产品事件
	r.POST("/product/update-event", updateProductEvent)                //添加产品事件
	r.GET("/product/query-event/:productId", queryProductEvent)        //查询所有产品事件
	r.GET("/product/delete-event/:productId/:key", deleteProductEvent) //删除产品事件
	r.GET("/product/detail-event/:productId/:key", detailProductEvent) //查看产品事件

	//产品告警
	r.POST("/product/add-alarm", addProductAlarm)                      //添加产品告警
	r.POST("/product/update-alarm", updateProductAlarm)                //添加产品告警
	r.GET("/product/query-alarm/:productId", queryProductAlarm)        //查询所有产品告警
	r.GET("/product/delete-alarm/:productId/:key", deleteProductAlarm) //删除产品告警
	r.GET("/product/detail-alarm/:productId/:key", detailProductAlarm) //查看产品告警

	//产品函数
	r.POST("/product/add-function", addProductFunction)                      //添加产品函数
	r.POST("/product/update-function", updateProductFunction)                //添加产品函数
	r.GET("/product/query-function/:productId", queryProductFunction)        //查询所有产品函数
	r.GET("/product/detail-function/:productId/:key", detailProductFunction) //查看产品函数
}

//添加产品
func addProduct(c *gin.Context) {
	var request dto.AddProductDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	//保存图片
	imagePath := path.Join(conf.GetConf().IotImagePath, uuid.NewV4().String()+path.Ext(request.Image.Filename))
	if err := c.SaveUploadedFile(request.Image, imagePath); err != nil {
		println(imagePath)
		logrus.Error(err)
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	productId, err := service.AddProductService(request, imagePath)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("保存成功", gin.H{
			"productId": productId,
		}))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

// 更新产品
func updateProduct(c *gin.Context) {
	var request dto.UpdateProductDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	var imagePath string
	if request.Image != nil {
		imagePath = path.Join(conf.GetConf().IotImagePath, uuid.NewV4().String()+path.Ext(request.Image.Filename))
		if err := c.SaveUploadedFile(request.Image, imagePath); err != nil {
			logrus.Error(err)
			c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		}
	}

	productId, err := service.UpdateProductService(request, imagePath)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("保存成功", gin.H{
			"productId": productId,
		}))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//查询产品列表
func queryProduct(c *gin.Context) {
	var request dto.QueryProductDataRequest
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

	count, page := service.QueryProductSerivce(request)
	c.JSON(http.StatusOK, common.OkPage(request.PageNo, request.PageSize, count, "", page))
}

//查询详情
func queryProductInfo(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}

	info, err := service.DetailProductSerivce(productID)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Ok("查询成功", info))
}

func setProductState(c *gin.Context) {
	var request dto.SetProductStateRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	err := service.SetProductStateSerivce(request.Id, *request.State)
	if err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	//设置
	ok := false
	if *request.State == 0 {
		ok = iot.StopProduct(request.Id)
	} else {
		ok = iot.StartProduct(request.Id)
	}
	if !ok {
		logrus.Errorf("change product[%d]'s state error", request.Id)
		c.JSON(http.StatusOK, common.Error("stop or start error", nil))
		return
	}
	c.JSON(http.StatusOK, common.Ok("设置成功", gin.H{
		"productId": request.Id,
	}))
}

//统计产品信息
func statisticsProduct(c *gin.Context) {
	data := service.StatisticsProductSerivce()
	c.JSON(http.StatusOK, common.Ok("保存成功", data))
}

//删除产品
func deleteProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err = service.DeleteProductSerivce(productID)
	if err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, common.Ok("删除成功", nil))
}

//添加产品物模型
func addProductItem(c *gin.Context) {
	var request dto.AddProductItemDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.AddProductItemService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", gin.H{}))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//更新产品物模型
func updateProductItem(c *gin.Context) {
	var request dto.UpdateProductItemDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.UpdateProductItemService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", gin.H{}))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//查询列表
func queryProductItem(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}

	list, err := service.QueryProductItemService(productID)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", list))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//删除产品物模型
func deleteProductItem(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	itemKey := c.Param("key")

	err = service.DeleteProductItemSerivce(productID, itemKey)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Ok("删除成功", nil))
}

//物模型详情
func detailProductItem(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	itemKey := c.Param("key")

	info, err := service.DetailProductItemSerivce(productID, itemKey)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Ok("查询成功", info))
}

//添加产品操作
func addProductOperation(c *gin.Context) {
	var request dto.AddProductOperationDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.AddProductOperationService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", gin.H{}))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//更新产品操作
func updateProductOperation(c *gin.Context) {
	var request dto.UpdateProductOperationDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.UpdateProductOperationService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", gin.H{}))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//查询操作列表
func queryProductOperation(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}

	list, err := service.QueryProductOperationService(productID)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", list))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//删除产品操作
func deleteProductOperation(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	itemKey := c.Param("key")

	err = service.DeleteProductOperationSerivce(productID, itemKey)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Ok("删除成功", nil))
}

//操作详情
func detailProductOperation(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	itemKey := c.Param("key")

	info, err := service.DetailProductOperationSerivce(productID, itemKey)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Ok("查询成功", info))
}

//添加产品操作
func addProductEvent(c *gin.Context) {
	var request dto.AddProductEventDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.AddProductEventService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", gin.H{}))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//更新产品操作
func updateProductEvent(c *gin.Context) {
	var request dto.UpdateProductEventDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.UpdateProductEventService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", gin.H{}))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//查询操作列表
func queryProductEvent(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}

	list, err := service.QueryProductEventService(productID)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", list))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//删除产品操作
func deleteProductEvent(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	itemKey := c.Param("key")

	err = service.DeleteProductEventSerivce(productID, itemKey)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Ok("删除成功", nil))
}

//操作详情
func detailProductEvent(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	itemKey := c.Param("key")

	info, err := service.DetailProductEventSerivce(productID, itemKey)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Ok("查询成功", info))
}

//添加产品告警
func addProductAlarm(c *gin.Context) {
	var request dto.AddProductAlarmDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.AddProductAlarmService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", gin.H{}))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//更新产品操作
func updateProductAlarm(c *gin.Context) {
	var request dto.UpdateProductAlarmDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.UpdateProductAlarmService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", gin.H{}))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//查询操作列表
func queryProductAlarm(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}

	list, err := service.QueryProductAlarmService(productID)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", list))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//删除产品操作
func deleteProductAlarm(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	itemKey := c.Param("key")

	err = service.DeleteProductAlarmSerivce(productID, itemKey)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Ok("删除成功", nil))
}

//操作详情
func detailProductAlarm(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	itemKey := c.Param("key")

	info, err := service.DetailProductAlarmSerivce(productID, itemKey)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Ok("查询成功", info))
}

//添加产品告警
func addProductFunction(c *gin.Context) {
	var request dto.AddProductFunctionDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.AddProductFunctionService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", gin.H{}))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//更新产品操作
func updateProductFunction(c *gin.Context) {
	var request dto.UpdateProductFunctionDataRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}

	err := service.UpdateProductFunctionService(request)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", gin.H{}))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//查询操作列表
func queryProductFunction(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}

	list, err := service.QueryProductFunctionService(productID)
	if err == nil {
		c.JSON(http.StatusOK, common.Ok("设置成功", list))
	} else {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	}
}

//操作详情
func detailProductFunction(c *gin.Context) {
	id := c.Param("productId")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	itemKey := c.Param("key")

	info, err := service.DetailProductFunctionSerivce(productID, itemKey)
	if err != nil {
		c.JSON(http.StatusOK, common.Error("参数错误", err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Ok("查询成功", info))
}
