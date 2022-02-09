package api

import (
	"encoding/json"
	"net/http"

	"koudai-box/global"

	"koudai-box/iot/service"
	"koudai-box/iot/web/common"

	"github.com/gin-gonic/gin"
)

func RegisterDictURL(r *gin.RouterGroup) {
	r.POST("/dict/codes", dictCodes)
}

func dictCodes(c *gin.Context) {
	var codes []string
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	err = json.Unmarshal(data, &codes)
	if err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
		return
	}
	response, err := service.DictCodesService(codes)
	if err != nil {
		c.JSON(http.StatusOK, common.Error(err.Error(), nil))
	} else {
		c.JSON(http.StatusOK, common.Ok(global.LIST_DIC_SUCCESS_MSG, response))
	}
}
