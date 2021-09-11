package driver

import (
	"main/model"
)

type Driver interface { //数据抽取接口
	FetchData() (interface{}, error)
	//数据抽取接口
	Extracter(data interface{}) (interface{}, error)
	//数据转换接口
	Transformer(data interface{}) (interface{}, error)
}

func GetDriver(gateway model.GatewayConfig, device model.Device) (Driver, bool) {
	var driver Driver
	ok := false
	switch gateway.Protocol {
	case "http":
		driver = &HttpDriver{gateway, device}
		ok = true
	case "Modbus":
	default:
	}
	return driver, ok
}
