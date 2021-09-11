package service

import (
	. "main/model"

	"github.com/sirupsen/logrus"
)

func GetGatewayConfigs() ([]GatewayConfig, bool) {
	var gatewayConfig GatewayConfig
	gatewayConfig.Key = "http-test"
	gatewayConfig.Name = "http测试"
	gatewayConfig.Ip = "124.160.72.210"
	gatewayConfig.Port = 8012
	gatewayConfig.Protocol = "http"
	gatewayConfig.ApiConfigs = map[string]ApiConfig{API_GetProp: {Key: API_GetProp, Method: "get", Name: "获取告警", Value: "/cm-admin/alarm/event/eventsRel/214"}}
	gatewayConfig.Parameters = []Parameter{{Key: "token", Name: "token", Value: "45c89b7fa77ab3d52b7eed083195c107"}}

	return []GatewayConfig{gatewayConfig}, true
}

func GetDevices(gateway string) ([]Device, bool) {
	p, ok := GetProduct("")
	if !ok {
		logrus.Errorf("product %s is not exist")
	}
	device := Device{Key: "test1", Name: "测试设备1", SourceId: "b6992eaf2fe2464da6189d7ea2dfdd1a", Product: p}
	return []Device{device}, true
}

func GetProduct(key string) (Product, bool) {

	//初始化设备和产品定义
	items := []ItemConfig{{Key: "captureId", Name: "图片ID", Source: "data[10].captureId", DataType: "string"},
		{Key: "deviceName", Name: "设备名称", Source: "data[10].deviceName", DataType: "string"},
		{Key: "deviceStatus", Name: "设备状态", Source: "data[10].deviceStatus", DataType: "string"}}

	operationConfigs := []OperationConfig{}

	extract := `function extract(data){
					console.log("hahah length=",data.data.length);
					//return data.data[0];
					if(data.data[10].deviceStatus == '离线'){
						data.data[10].deviceStatus = 0;
					}else{
						data.data[10].deviceStatus = 1;
					}
					return data;
				}`
	calc := `function calculate(data){
					console.log("calculate:",data.MessageId,data.DeviceId);
					return data;
				}`
	functionConfigs := map[string]FunctionConfig{Function_Extract: {Key: Function_Extract, Name: "数据抽取", Function: extract},
		Function_Calc: {Key: Function_Calc, Name: "数据计算", Function: calc}}

	operations := []OperationConfig{}
	conditions := []Condition{{Key: "deviceStatus", Name: "设备状态", DataType: "int", Compare: "=", Value: "0"}}

	alarmConfigs := []AlarmConfig{{Key: "offline", Name: "设备离线", Level: "1", Conditions: conditions, Operations: operations, Message: "设备离线"}}

	product := Product{Key: "p1", Name: "产品1", CollectPeriod: 20, DataCombination: "array",
		Items: items, OperationConfigs: operationConfigs, AlarmConfigs: alarmConfigs, FunctionConfigs: functionConfigs}
	return product, true
}
