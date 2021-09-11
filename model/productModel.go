package model

const API_GetProp string = "getProp"
const API_SetProp string = "setProp"

const Function_Extract string = "extract"
const Function_Calc string = "calculate"
const Function_postBody string = "postBody"

type Product struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Category string `json:"category"`
	//GatewayConfig    GatewayConfig             `json:"gatewayConfig"`
	CollectPeriod    int                       `json:"collectPeriod"`
	DataCombination  string                    `json:"dataCombination"`
	Desc             string                    `json:"desc"`
	Items            []ItemConfig              `json:"items"`
	OperationConfigs []OperationConfig         `json:"operationConfigs"`
	AlarmConfigs     []AlarmConfig             `json:"alarmSettings"`
	FunctionConfigs  map[string]FunctionConfig `json:"functionConfigs"`
}
type ModelModbus struct {
	Items []ItemConfigModbus `json:"items"`
}
type ModelOPCUA struct {
	Items []ItemConfigOPC `json:"items"`
}

/**物模型定义**/
type ItemConfigModbus struct {
	Key            string `json:"key"`
	Name           string `json:"name"`
	OperaterType   string `json:"operaterType"`
	Address        string `json:"address"`
	DataType       string `json:"dataType"`
	Unit           string `json:"unit"`
	Scale          string `json:"scale"`
	ZoomFactor     string `json:"zoomFactor"`
	ExchangeHL     bool   `json:"exchangeHL"`
	ExchangeOrder  string `json:"exchangeOrder"`
	DataReportType string `json:"dataReportType"`
	Desc           string `json:"desc"`
}
type ItemConfigOPC struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	NodeId   string `json:"nodeId"`
	DataType string `json:"dataType"`
	Unit     string `json:"unit"`
	Scale    string `json:"scale"`
	StepSize string `json:"stepSize"`
	RW       string `json:"rw"`
	Desc     string `json:"desc"`
}
type ItemConfig struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Source   string `json:"source"`
	DataType string `json:"dataType"`
	Unit     string `json:"unit"`
	Scale    string `json:"scale"`
	StepSize string `json:"stepSize"`
	RW       string `json:"rw"`
	Desc     string `json:"desc"`
}

/**操作定义**/
type OperationConfig struct {
	Key       string      `json:"key"`
	Name      string      `json:"name"`
	Inputs    []Parameter `json:"inputs"`
	ApiConfig ApiConfig   `json:"apiConfig"`
	Desc      string      `json:"desc"`
}

type Parameter struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ApiConfig struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Method string `json:"mthod"`
}

/**告警定义**/
type AlarmConfig struct {
	Key        string            `json:"key"`
	Name       string            `json:"name"`
	Level      string            `json:"level"`
	Conditions []Condition       `json:"conditions"`
	Operations []OperationConfig `json:"operations"`
	Message    string             `json:"messae"`
}
type Condition struct {
	Key      string      `json:"key"`
	Name     string      `json:"name"`
	DataType string      `json:"dataType"`
	Compare  string      `json:"compare"`
	Value    interface{} `json:"vaue"`
}
type FunctionConfig struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Function string `json:"function"`
}
type GatewayConfig struct {
	Key        string               `json:"key"`
	Name       string               `json:"name"`
	Protocol   string               `json:"protocol"`
	Ip         string               `json:"ip"`
	Port       int                  `json:"port"`
	Parameters []Parameter          `json:"parameters"`
	ApiConfigs map[string]ApiConfig `json:"apiConfigs"`
	Desc       string               `json:"desc"`
}
type DataType int

// iota 初始化后会自动递增
const (
	Running    DataType = iota // value --> 0
	Stopped                    // value --> 1
	Rebooting                  // value --> 2
	Terminated                 // value --> 3
)

// int32：32位整型。需定义取值范围、步长和单位符号。
// float：单精度浮点型。需定义取值范围、步长和单位符号。
// double：双精度浮点型。需定义取值范围、步长和单位符号。
// enum：枚举型。定义枚举项的参数值和参数描述，例如：1表示加热模式、2表示制冷模式。
// bool：布尔型。采用0或1来定义布尔值，例如：0表示关、1表示开。
// text：字符串。需定义字符串的数据长度，最长支持10240字节。
// date：时间戳。格式为String类型的UTC时间戳，单位：毫秒。
// struct：JSON对象。定义一个JSON结构体，新增JSON参数项，例如：定义灯的颜色是由Red、Green、Blue三个参数组成的结构体。不支持结构体嵌套。
// array：数组。需声明数组内的元素类型、数组元素个数。元素类型可选择int32、float、double、text或struct，需确保同一个数组元素类型相同。元素个数，限制1~512个。
