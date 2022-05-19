package service

import (
	"encoding/json"

	"koudai-box/iot/gateway/model"
	"koudai-box/iot/gateway/utils"
	"koudai-box/iot/service"
	"koudai-box/iot/web/dto"

	"github.com/sirupsen/logrus"

	db "koudai-box/iot/db"
)

const (
	Status_Stop    = 0
	Status_Running = 1
)

func GetGatewayConfig(gatewayId int) *model.GatewayConfig {
	gateways := service.GetGatewayCache()
	g, ok := gateways[gatewayId]
	if !ok {
		return nil
	}
	gatewayConfig := convertGateway(g)
	return &gatewayConfig
}
func GetGatewayConfigByKey(gatewayKey string) *model.GatewayConfig {
	gateways := service.GetGatewayCache()
	for _, v := range gateways {
		if v.Sign == gatewayKey {
			gatewayConfig := convertGateway(v)
			return &gatewayConfig
		}
	}
	return nil
}
func GetGatewayConfigs() ([]model.GatewayConfig, bool) {
	gateways := service.GetGatewayCache()
	if gateways == nil {
		logrus.Error("gateway is nil")
		return nil, false
	}
	var gatewayconfigs []model.GatewayConfig
	for _, gateway := range gateways {
		if gateway.Status == Status_Stop {
			continue
		}
		gatewayConfig := convertGateway(gateway)
		gatewayconfigs = append(gatewayconfigs, gatewayConfig)
	}

	return gatewayconfigs, true
}

func GetDevices(gatewayId int, status int) ([]*model.Device, bool) {
	ps, err := db.QueryProductByGatewayID(gatewayId)
	if err != nil || len(ps) == 0 {
		logrus.Errorf("gateway[%d]'s product is not exist", gatewayId)
		return nil, false
	}
	var ids []int
	for _, p := range ps {
		if p.State == 0 {
			continue
		}
		ids = append(ids, p.Id)
	}
	list, ok := GetDeviceByProductIds(ids, status)
	return list, ok
}

func GetDeviceById(deviceId int, status int) (*model.Device, bool) {
	d, err := db.QueryDeviceByID(deviceId)
	if err != nil {
		logrus.Errorf("Device %d is not exist", deviceId)
		return nil, false
	}
	if status != model.STATUS_ALL && d.ActivateStatus != status {
		logrus.Errorf("Device %d status is %d but not %d", deviceId, d.ActivateStatus, status)
		return nil, false
	}
	device := model.Device{Id: d.Id, Key: d.Code, SourceId: d.SourceId, Name: d.Name, Desc: d.Desc, ExtProps: utils.ToMap(d.ExtProps)}
	if len(device.SourceId) == 0 {
		device.SourceId = device.Key
	}
	p, ok := GetProduct(d.ProductId, status)
	if ok {
		device.Product = p
	}
	return &device, true
}

func GetDeviceByProductIds(ids []int, status int) ([]*model.Device, bool) {
	if len(ids) == 0 {
		return nil, false
	}
	ds := db.QueryDevicetByProductIds(ids)
	list := make([]*model.Device, 0)
	products := make(map[int]*model.Product)
	for _, d := range ds {
		p, ok := products[d.ProductId]
		if !ok {
			p, ok = GetProduct(d.ProductId, model.STATUS_ALL)
			if ok {
				products[d.ProductId] = p
			}
		}
		if status != model.STATUS_ALL && d.ActivateStatus != status {
			logrus.Errorf("Device %d status is %d but not %d", d.Id, d.ActivateStatus, status)
			continue
		}
		device := model.Device{Id: d.Id, Key: d.Code, SourceId: d.SourceId, Name: d.Name, Desc: d.Desc, ExtProps: utils.ToMap(d.ExtProps)}
		if len(device.SourceId) == 0 {
			device.SourceId = device.Key
		}
		if ok {
			device.Product = p
		}
		list = append(list, &device)
	}
	return list, true
}

func GetProduct(id int, status int) (*model.Product, bool) {
	p, err := db.QueryProductByID(id)
	if err != nil {
		return nil, false
	}
	if status != model.STATUS_ALL && p.State != status {
		logrus.Errorf("Product %d state is %d but not %d", id, p.State, status)
		return nil, false
	}
	tmp := model.Product{Id: p.Id, Key: p.Code, Name: p.Name, GatewayId: p.GatewayId,
		Items: convertItemConfig(p.Items), AlarmConfigs: convertAlarmConfig(p.AlarmConfigs), FunctionConfigs: converFunctionConfig(p.FunctionConfigs)}
	return &tmp, true
}

func GetProductXXXX(key string) (model.Product, bool) {

	//初始化设备和产品定义
	items := []model.ItemConfig{{Key: "code", Name: "股票代码", Source: "code", DataType: model.ItemDataType{Type: model.Text}},
		{Key: "hq", Name: "行情", Source: "hq", DataType: model.ItemDataType{Type: model.Text}},
		{Key: "status", Name: "状态", Source: "status", DataType: model.ItemDataType{Type: model.Int32, Dict: map[string]string{"0": "离线", "1": "在线"}}}}

	operationConfigs := []model.OperationConfig{}

	extract := `function extract(data){
					console.log("hahah length=",data.length);
					//return data.data[0];
					if(data[0].status == 0){
						data[0].status = 'OK';
					}else{
						data[0].status = 'NG';
					}
					return data[0];
				}`
	calc := `function calculate(data){
					console.log("calculate:",data.MessageId,data.DeviceId);
					return data;
				}`
	functionConfigs := map[string]model.FunctionConfig{model.Function_Extract_Prop: {Key: model.Function_Extract_Prop, Name: "数据抽取", Function: extract},
		model.Function_Calc: {Key: model.Function_Calc, Name: "数据计算", Function: calc}}

	operations := []model.OperationConfig{}
	conditions := []model.Condition{{Key: "status", Name: "设备状态", DataType: model.Text, Compare: "=", Value: "NG"}}
	conditions1 := []model.Condition{{Key: "status", Name: "设备状态", DataType: model.Text, Compare: "=", Value: "OK"}}

	alarmConfigs := []model.AlarmConfig{{Key: "offline", Name: "设备离线", Level: "1", Type: "event", Conditions: conditions, Operations: operations, Message: "设备离线"},
		{Key: "online", Name: "设备上线", Level: "1", Type: "event", Conditions: conditions1, Message: "设备上线"}}

	product := model.Product{Id: 1, Key: "p1", Name: "产品1",
		Items: items, OperationConfigs: operationConfigs, AlarmConfigs: alarmConfigs, FunctionConfigs: functionConfigs}
	return product, true
}

func convertGateway(gateway *dto.GatewayItem) model.GatewayConfig {
	gatewayConfig := model.GatewayConfig{
		Id:            gateway.GatewayId,
		Key:           gateway.Sign,
		Name:          gateway.GatewayName,
		Ip:            gateway.Ip,
		Port:          gateway.Port,
		Protocol:      gateway.Protocol,
		CollectType:   gateway.CollectType,
		CollectPeriod: gateway.CollectPeriod,
		Cron:          gateway.Cron,
		ModbusConfig:  gateway.ModbusConfig,
	}
	auths := gateway.AuthInfo
	for _, auth := range auths {
		gatewayConfig.Parameters = append(gatewayConfig.Parameters, model.Parameter{Key: auth.Key, Value: auth.Value})
	}
	routers := gateway.Routers
	if routers != nil {
		gatewayConfig.ApiConfigs = make(map[string]model.ApiConfig)
		for _, router := range routers {
			gatewayConfig.ApiConfigs[router.RouterName] = model.ApiConfig{
				Name:            router.RouterName,
				Method:          router.RequestType,
				Path:            router.RouterUrl,
				CollectType:     router.CollectType,
				CollectPeriod:   router.Period,
				DataCombination: router.DataCombination,
				Cron:            router.CornExpression,
			}
		}
	}
	return gatewayConfig
}
func convertItemConfig(items string) []model.ItemConfig {
	itemConfigs, err := service.ConvertItemConfig(items)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	list := make([]model.ItemConfig, 0)
	for _, item := range itemConfigs {
		d := model.ItemConfig{
			Sort:           item.Sort,
			Key:            item.Key,
			Code:           item.Code,
			Name:           item.Name,
			Source:         item.SourceCode,
			NodeId:         item.NodeId,        //opcua
			Address:        item.Address,       //modbus
			Quantity:       item.Quantity,      //modbus
			OperaterType:   item.OperaterType,  //modbus
			ZoomFactor:     item.ZoomFactor,    //modbus
			ExchangeHL:     item.ExchangeHL,    //modbus
			ExchangeOrder:  item.ExchangeOrder, //modbus
			DataType:       item.DataType,
			DataReportType: model.DataReportType_Schedule,
			Desc:           item.Desc,
		}
		if len(d.Source) == 0 {
			d.Source = item.Address //modbus
		}
		if len(d.Source) == 0 {
			d.Source = item.NodeId //opcua
		}
		list = append(list, d)
	}
	return list
}

func convertAlarmConfig(alarms string) []model.AlarmConfig {
	var alarmConfigs []model.AlarmConfig
	err := json.Unmarshal([]byte(alarms), &alarmConfigs)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return alarmConfigs
}

func converFunctionConfig(functions string) map[string]model.FunctionConfig {
	var functionConfig *dto.ProductFunctionConfig
	err := json.Unmarshal([]byte(functions), &functionConfig)
	if err != nil {
		logrus.Errorln(err)
		functionConfig = &dto.ProductFunctionConfig{}
	}
	// extract := `//将多条数据转换成map形式，{设备id:设备数据}
	// 			function extract(data){
	// 				var length = data.length;
	// 				console.log("data length=",length);
	// 				var obj = {};
	// 				for(var i = 0;i < length; i++){
	// 					var item = data[i];
	// 					obj[item.code] = item;
	// 				}
	// 				return obj;
	// 			}`
	// extractEvent := `/**
	// 					*处理设备单条数据
	// 					*return{
	// 					* 属性标识1:属性值,
	// 					* 属性标识2:属性值
	// 					*}
	// 					*/
	// 					function extractEvent(param){
	// 						data=param[0]
	// 						var obj = {};
	// 						obj.DeviceSign = data.code;
	// 						obj.Type = "event";
	// 						obj.Title = data.status+"";
	// 						obj.Message = data.code;
	// 						obj.Properties = data;
	// 						return obj;
	// 					}`
	functionConfigs := map[string]model.FunctionConfig{
		model.Function_Extract_Prop:  {Key: model.Function_Extract_Prop, Name: "数据抽取", Function: functionConfig.ExtractProp},
		model.Function_Extract_Event: {Key: model.Function_Extract_Event, Name: "事件抽取", Function: functionConfig.ExtractEvent},
		model.Function_Calc:          {Key: model.Function_Calc, Name: "数据计算", Function: functionConfig.Calculate},
		model.Function_postBody:      {Key: model.Function_postBody, Name: "请求体处理", Function: functionConfig.Body}}
	return functionConfigs
}
