package driver

import (
	"main/model"
)

type Driver interface {
	//数据抽取接口,批量
	FetchPropBatch() (interface{}, error)
	//数据抽取接口,单条
	FetchProp(device model.Device) (interface{}, error)
	//数据抽取接口
	ExtracterProp(data interface{}, device model.Device) (interface{}, error)
	//数据转换接口
	TransformerProp(data interface{}, device model.Device) (interface{}, error)

	GetCollectPeriod(key string) int
}

func GetDriver(gateway model.GatewayConfig, device model.Device) (Driver, bool) {
	var driver Driver
	ok := false
	switch gateway.Protocol {
	case "http":
		driver = &HttpDriver{gateway}
		ok = true
	case "Modbus":
	default:
	}
	return driver, ok
}
