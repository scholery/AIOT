package test

import (
	"fmt"
	. "main/driver"
	"main/model"
	. "main/model"
	"main/service"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

var httpDriver HttpDriver
var device Device

func Init() {
	//初始化网关
	var gatewayConfig GatewayConfig
	gatewayConfig.Key = "http-test"
	gatewayConfig.Name = "http测试"
	gatewayConfig.Ip = "124.160.72.210"
	gatewayConfig.Port = 8012
	gatewayConfig.Protocol = "http"
	gatewayConfig.ApiConfigs = map[string]model.ApiConfig{model.API_GetProp: {Key: model.API_GetProp, Method: "get", Name: "获取告警", Path: "/cm-admin/alarm/event/eventsRel/214", CollectType: model.CollectType_Schedule, CollectPeriod: 20, DataCombination: model.DataCombination_Single}}
	gatewayConfig.Parameters = []Parameter{{Key: "token", Name: "token", Value: "45c89b7fa77ab3d52b7eed083195c107"}}

	//初始化设备和产品定义
	items := []ItemConfig{{Key: "captureId", Name: "图片ID", Source: "data[10].captureId", DataType: model.ItemDataType{Type: Text}},
		{Key: "deviceName", Name: "设备名称", Source: "data[10].deviceName", DataType: model.ItemDataType{Type: Text}}}

	operationConfigs := []OperationConfig{}

	extract := `function extract(data){
					console.log("hahah length=",data.data.length);
					//return data.data[0];
					return data;
				}`
	calc := `function calculate(data){
					console.log("calculate:",data.MessageId,data.DeviceId);
					return data;
				}`
	functionConfigs := map[string]FunctionConfig{Function_Extract: {Key: Function_Extract, Name: "数据抽取", Function: extract},
		Function_Calc: {Key: Function_Calc, Name: "数据计算", Function: calc}}

	product := Product{Key: "p1", Name: "产品1",
		Items: items, OperationConfigs: operationConfigs, FunctionConfigs: functionConfigs}

	device = Device{Key: "test1", Name: "测试设备1", SourceId: "b6992eaf2fe2464da6189d7ea2dfdd1a", Product: product}

	httpDriver = HttpDriver{Gateway: gatewayConfig}
}

func ExecHttpTest(c chan PropertyMessage) {
	for {
		logrus.Info("gateway run")
		ExecHttp(c)
		period := httpDriver.GetCollectPeriod(model.API_GetProp)
		time.Sleep(time.Duration(period) * time.Second)
	}
}

func ExecHttp(c chan PropertyMessage) {
	start := time.Now() // 获取当前时间
	Init()
	var driver Driver = &httpDriver
	data, err := driver.FetchProp(device)
	//logrus.Debug("FetchData:", data)
	if err != nil {
		logrus.Error(err)
		return
	}
	data, err = driver.ExtracterProp(data, device)
	//logrus.Debug("Extracter:", data.(map[string]interface{})["deviceName"])
	if err != nil {
		logrus.Error(err)
		return
	}
	data, err = driver.TransformerProp(data, device)
	logrus.Info("Transformer:", data)
	if err != nil {
		logrus.Error(err)
		return
	}
	elapsed := time.Since(start)
	logrus.Info("ExecHttp执行完成耗时：", elapsed)
	tmp, ok := data.(PropertyMessage)
	if ok {
		c <- tmp
	}
}

func ExecCalc(data PropertyMessage) {
	start := time.Now() // 获取当前时间
	dataGateway := &service.DataGateway{Device: device}
	res, err := dataGateway.Calculater(data)
	if err != nil {
		logrus.Error(err)
		return
	}
	tmpP, ok := res.(PropertyMessage)
	if !ok {
		logrus.Error("calc error")
		return
	}
	dataGateway.LoaderMessage(tmpP)
	service.Public(tmpP, service.Router_prop)
	alarms, err := dataGateway.FilterAlarm(tmpP)
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, alarm := range alarms {
		tmpA, ok := alarm.(AlarmMessage)
		if !ok {
			logrus.Error("alarm is null")
		}
		service.Public(tmpA, service.Router_prop)
		dataGateway.LoaderAlarm(tmpA)
	}
	elapsed := time.Since(start)
	logrus.Info("ExecCalc执行完成耗时：", elapsed)

}

func TestHttp(t *testing.T) {
	t.Log("test http")
	logrus.SetLevel(logrus.DebugLevel)
	// 开始性能分析, 返回一个停止接口
	// stopper := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	// 在main()结束时停止性能分析
	// defer stopper.Stop()

	c := make(chan PropertyMessage)
	service.Connect()
	defer service.Close()
	go ExecHttpTest(c)
	// time.Sleep(time.Second)
	i := 1
	for data := range c {
		fmt.Println("data=", data.MessageId)
		fmt.Println("index=", i)
		ExecCalc(data)
		i++
		if i > 3 {
			break
		}
	}
	// 让程序至少运行1秒
	// time.Sleep(time.Second)
}
