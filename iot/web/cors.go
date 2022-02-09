package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//Cors 跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部

		// var headerKeys []string // 声明请求头keys
		// for k, _ := range c.Request.Header {
		// 	headerKeys = append(headerKeys, k)
		// }
		// headerStr := strings.Join(headerKeys, ", ")
		// println(headerStr)
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Content-Disposition")
			c.Header("Access-Control-Max-Age", "172800")          // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false") //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()

	}
}
