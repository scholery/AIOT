package httpServer

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"koudai-box/iot/gateway/model"
	utils "koudai-box/iot/gateway/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//httpserver
var cache_servers map[int]*http.Server = make(map[int]*http.Server)

type HttpServerDriver struct {
	Gateway *model.GatewayConfig
}

func (httpDriver *HttpServerDriver) GetGatewayConfig() *model.GatewayConfig {
	return httpDriver.Gateway
}

func (httpDriver *HttpServerDriver) Start() error {
	gateway := httpDriver.GetGatewayConfig()
	if _, ok := cache_servers[gateway.Id]; ok {
		return fmt.Errorf("http server[%s] has been started", gateway.Key)
	}
	gin.DefaultWriter = io.Discard
	r := gin.Default()
	r.MaxMultipartMemory = 32 << 20
	r.Use(GinHead())
	r.BasePath()
	gr := r.Group("/app/api/v1/iot")
	RegisterURL(gr)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": "Api not found"})
	})

	logrus.Infof("++++++++++++++开启Http server[%s]应用[%s:%d]++++++++++++++", gateway.Key, gateway.Ip, gateway.Port)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", gateway.Ip, gateway.Port),
		Handler: r,
	}
	cache_servers[gateway.Id] = httpServer
	go httpServer.ListenAndServe()
	return nil
}

func (httpDriver *HttpServerDriver) Stop() error {
	gateway := httpDriver.GetGatewayConfig()
	server, ok := cache_servers[gateway.Id]
	if !ok {
		return fmt.Errorf("http server[%s] has not start", gateway.Key)
	}
	server.Close()
	logrus.Infof("++++++++++++++关闭Http server[%s]应用[%s]++++++++++++++", gateway.Key, server.Addr)
	delete(cache_servers, gateway.Id)
	return nil
}
func (httpDriver *HttpServerDriver) StartDevicde(device *model.Device) error {
	return nil
}

func (httpDriver *HttpServerDriver) StopDevicde(device *model.Device) error {
	return nil
}

func GinHead() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Next()
	}
}

func (httpDriver *HttpServerDriver) FetchPropBatch(ts int64) (interface{}, error) {
	return nil, nil
}

func (httpDriver *HttpServerDriver) FetchProp(device *model.Device, ts int64) (interface{}, error) {
	return nil, nil
}

//数据抽取接口
func (httpDriver *HttpServerDriver) ExtracterProp(data interface{}, product *model.Product) (interface{}, error) {
	if len(product.FunctionConfigs) == 0 {
		return nil, errors.New("function is null")
	}
	function, ok := product.FunctionConfigs[model.Function_Extract_Prop]
	if !ok || len(function.Function) == 0 {
		return nil, errors.New("function is null")
	}
	logrus.Debugf("Extracter funtion name=%s,function[%s]", model.Function_Extract_Prop, function.Function)
	return utils.ExecJS(function.Function, function.Key, data)
}

//数据转换接口
func (httpDriver *HttpServerDriver) TransformerProp(data interface{}, device *model.Device) (interface{}, error) {
	return utils.Transformer2DeviceProp(data, device)
}

//数据抽取接口,单条
func (httpDriver *HttpServerDriver) FetchEvent(device *model.Device, ts int64) (interface{}, error) {
	return nil, errors.New("no implement")
}

//网关数据发送设备
func (httpDriver *HttpServerDriver) PostOperation(api model.ApiConfig, data interface{}, device *model.Device) (interface{}, error) {
	return nil, errors.New("no implement")
}

//事件抽取接口
func (httpDriver *HttpServerDriver) ExtracterEvent(data interface{}, product *model.Product) (interface{}, error) {
	if len(product.FunctionConfigs) == 0 {
		return nil, errors.New("function is null")
	}
	function, ok := product.FunctionConfigs[model.Function_Extract_Event]
	if !ok || len(function.Function) == 0 {
		return nil, errors.New("function is null")
	}
	logrus.Debugf("Extracter funtion name=%s,function[%s]", model.Function_Extract_Event, function.Function)
	return utils.ExecJS(function.Function, function.Key, data)
}
