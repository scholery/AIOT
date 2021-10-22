package service

import (
	"main/model"

	"github.com/sirupsen/logrus"
)

func GetGatewayConfigs() ([]model.GatewayConfig, bool) {
	var gatewayConfig model.GatewayConfig
	gatewayConfig.Key = "http-test"
	gatewayConfig.Name = "http测试"
	gatewayConfig.Ip = "124.160.72.210"
	gatewayConfig.Port = 8012
	gatewayConfig.Protocol = "http"
	gatewayConfig.ApiConfigs = map[string]model.ApiConfig{model.API_GetProp: {Key: model.API_GetProp, Method: "get", Name: "获取告警", Path: "/cm-admin/alarm/event/eventsRel/214", CollectType: model.CollectType_Schedule, CollectPeriod: 20, DataCombination: model.DataCombination_Single}}
	gatewayConfig.Parameters = []model.Parameter{{Key: "token", Name: "token", Value: "zhejiang_data_push_token"}}

	return []model.GatewayConfig{gatewayConfig}, true
}

func GetDevices(gateway string) ([]model.Device, bool) {
	p, ok := GetProduct("")
	if !ok {
		logrus.Errorf("product %s is not exist")
	}
	device := model.Device{Key: "test1", Name: "测试设备1", SourceId: "b6992eaf2fe2464da6189d7ea2dfdd1a", Product: p}
	return []model.Device{device}, true
}

func GetProduct(key string) (model.Product, bool) {

	//初始化设备和产品定义
	items := []model.ItemConfig{{Key: "captureId", Name: "图片ID", Source: "data[10].captureId", DataType: model.ItemDataType{Type: model.Text}},
		{Key: "deviceName", Name: "设备名称", Source: "data[10].deviceName", DataType: model.ItemDataType{Type: model.Text}},
		{Key: "deviceStatus", Name: "设备状态", Source: "data[10].deviceStatus", DataType: model.ItemDataType{Type: model.Int32, Dict: map[string]string{"0": "离线", "1": "在线"}}}}

	operationConfigs := []model.OperationConfig{}

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
	functionConfigs := map[string]model.FunctionConfig{model.Function_Extract: {Key: model.Function_Extract, Name: "数据抽取", Function: extract},
		model.Function_Calc: {Key: model.Function_Calc, Name: "数据计算", Function: calc}}

	operations := []model.OperationConfig{}
	conditions := []model.Condition{{Key: "deviceStatus", Name: "设备状态", DataType: model.Int32, Compare: "=", Value: 0}}
	conditions1 := []model.Condition{{Key: "deviceStatus", Name: "设备状态", DataType: model.Int32, Compare: "=", Value: 1}}

	alarmConfigs := []model.AlarmConfig{{Key: "offline", Name: "设备离线", Level: "1", Type: "event", Conditions: conditions, Operations: operations, Message: "设备离线"},
		{Key: "online", Name: "设备上线", Level: "1", Type: "event", Conditions: conditions1, Message: "设备上线"}}

	product := model.Product{Key: "p1", Name: "产品1",
		Items: items, OperationConfigs: operationConfigs, AlarmConfigs: alarmConfigs, FunctionConfigs: functionConfigs}
	return product, true
}
