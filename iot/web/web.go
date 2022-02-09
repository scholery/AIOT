package web

import (
	"fmt"
	"io"
	"net/http"

	"koudai-box/iot/web/api"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GinHead() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		c.Next()
	}
}

func Init(port int) *http.Server {
	gin.DefaultWriter = io.Discard

	r := gin.Default()
	r.MaxMultipartMemory = 32 << 20
	r.Use(GinHead())
	//开启跨域设置
	r.Use(Cors())

	r.BasePath()

	gr := r.Group("/app/api/v1/iot")
	api.RegisterGatewayURL(gr)
	api.RegisterAlarmURL(gr)
	api.RegisterEventURL(gr)
	api.RegisterDictURL(gr)
	api.RegisterProductURL(gr)
	api.RegisterDeviceURL(gr)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	logrus.Infof("***************************开启web应用，端口:%d***************************", port)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	return httpServer
}
