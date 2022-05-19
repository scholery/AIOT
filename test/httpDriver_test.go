package test

import (
	"fmt"
	"testing"
	"time"

	driver "koudai-box/iot/gateway/driver"
	"koudai-box/iot/gateway/driver/httpClient"
	model "koudai-box/iot/gateway/model"
	"koudai-box/iot/gateway/service"

	"github.com/sirupsen/logrus"
)

var httpDriver httpClient.HttpDriver
var device *model.Device

func Init() {
	//初始化网关
	var gatewayConfig model.GatewayConfig
	gatewayConfig.Key = "http-test"
	gatewayConfig.Name = "http测试"
	gatewayConfig.Ip = "124.160.72.210"
	gatewayConfig.Port = 8012
	gatewayConfig.Protocol = "http"
	gatewayConfig.ApiConfigs = map[string]model.ApiConfig{model.API_GetProp: {Method: "get", Name: "获取告警", Path: "/cm-admin/alarm/event/eventsRel/214", CollectType: model.CollectType_Schedule, CollectPeriod: 20, DataCombination: model.DataCombination_Single}}
	gatewayConfig.Parameters = []model.Parameter{{Key: "token", Name: "token", Value: "45c89b7fa77ab3d52b7eed083195c107"}}

	//初始化设备和产品定义
	items := []model.ItemConfig{{Key: "captureId", Name: "图片ID", Source: "data[10].captureId", DataType: model.ItemDataType{Type: model.Text}},
		{Key: "deviceName", Name: "设备名称", Source: "data[10].deviceName", DataType: model.ItemDataType{Type: model.Text}}}

	operationConfigs := []model.OperationConfig{}

	extract := `function extract(data){
					console.log("hahah length=",data.data.length);
					//return data.data[0];
					return data;
				}`
	calc := `function calculate(data){
					console.log("calculate:",data.MessageId,data.DeviceId);
					return data;
				}`
	functionConfigs := map[string]model.FunctionConfig{model.Function_Extract_Prop: {Key: model.Function_Extract_Prop, Name: "数据抽取", Function: extract},
		model.Function_Calc: {Key: model.Function_Calc, Name: "数据计算", Function: calc}}

	product := model.Product{Key: "p1", Name: "产品1",
		Items: items, OperationConfigs: operationConfigs, FunctionConfigs: functionConfigs}

	device = &model.Device{Key: "b6992eaf2fe2464da6189d7ea2dfdd1a", Name: "测试设备1", Product: &product}

	httpDriver = httpClient.HttpDriver{Gateway: &gatewayConfig}
}

func ExecHttpTest(c chan model.PropertyMessage) {
	for {
		logrus.Info("gateway run")
		ExecHttp(c)
		period := 20
		time.Sleep(time.Duration(period) * time.Second)
	}
}

func ExecHttp(c chan model.PropertyMessage) {
	start := time.Now() // 获取当前时间
	Init()
	var dri driver.Driver = &httpDriver
	data, err := dri.FetchProp(device, -1)
	//logrus.Debug("FetchData:", data)
	if err != nil {
		logrus.Error(err)
		return
	}
	data, err = dri.ExtracterProp(data, device.Product)
	//logrus.Debug("Extracter:", data.(map[string]interface{})["deviceName"])
	if err != nil {
		logrus.Error(err)
		return
	}
	data, err = dri.TransformerProp(data, device)
	logrus.Info("Transformer:", data)
	if err != nil {
		logrus.Error(err)
		return
	}
	elapsed := time.Since(start)
	logrus.Info("ExecHttp执行完成耗时：", elapsed)
	tmp, ok := data.(model.PropertyMessage)
	if ok {
		c <- tmp
	}
}

func ExecCalc(data model.PropertyMessage) {
	start := time.Now() // 获取当前时间
	dataGateway := &service.DataGateway{Device: device}
	res, err := dataGateway.Calculater(data)
	if err != nil {
		logrus.Error(err)
		return
	}
	tmpP, ok := res.(model.PropertyMessage)
	if !ok {
		logrus.Error("calc error")
		return
	}
	dataGateway.LoaderProperty(tmpP, true)
	alarms, err := dataGateway.FilterAlarm(tmpP)
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, alarm := range alarms {
		tmpA, ok := alarm.(model.IotEventMessage)
		if !ok {
			logrus.Error("alarm is null")
		}
		dataGateway.LoaderAlarm(tmpA, true)
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

	c := make(chan model.PropertyMessage)
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
