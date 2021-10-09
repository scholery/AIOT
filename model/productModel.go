package model

const API_GetProp string = "getProp" //取数据接口
const API_SetProp string = "setProp" //写数据接口

const Function_Extract string = "extract"   //数据抽取函数
const Function_Calc string = "calculate"    //计算函数
const Function_postBody string = "postBody" //数据写入前处理函数

const DataReportType_Schedule string = "schedule" //按时上报
const DataReportType_Change string = "change"     //变化上报

type Product struct {
	Key              string                    `json:"key"`
	Name             string                    `json:"name"`
	Category         string                    `json:"category"`        //分类
	CollectPeriod    int                       `json:"collectPeriod"`   //采集周期
	DataCombination  string                    `json:"dataCombination"` //数据类型，单条、数组
	Desc             string                    `json:"desc"`
	Items            []ItemConfig              `json:"items"`            //物模型
	OperationConfigs []OperationConfig         `json:"operationConfigs"` //操作定义
	AlarmConfigs     []AlarmConfig             `json:"alarmSettings"`    //告警配置
	FunctionConfigs  map[string]FunctionConfig `json:"functionConfigs"`  //函数配置
}

/**物模型定义**/
type ItemConfig struct {
	Key            string       `json:"key"`
	Name           string       `json:"name"`
	OperaterType   string       `json:"operaterType"`   //modbus，操作类型
	Address        string       `json:"address"`        //modbus，地址
	Source         string       `json:"source"`         //源数据标识
	NodeId         string       `json:"nodeId"`         //opc ua 节点
	DataType       ItemDataType `json:"dataType"`       //数据类型
	ZoomFactor     string       `json:"zoomFactor"`     //modbus，缩放因子：不能为0，默认为1
	ExchangeHL     bool         `json:"exchangeHL"`     //modbus，交换寄存器内高低字节：true：互换，false：不互换
	ExchangeOrder  string       `json:"exchangeOrder"`  //modbus，交换寄存器顺序：true：互换，false：不互换
	DataReportType string       `json:"dataReportType"` //modbus，数据上报方式，可选按时上报和变更上报：change、schedule
	Desc           string       `json:"desc"`
}

type ItemDataType struct {
	RW   string            `json:"rw"` //读写标志
	Type DataType          `json:"type"`
	Unit string            `json:"unit"`
	Min  string            `json:"min"`
	Max  string            `json:"max"`
	Step string            `json:"step"`
	Dict map[string]string `json:"dict"`
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
	Type       string            `json:"type"`
	Conditions []Condition       `json:"conditions"`
	Operations []OperationConfig `json:"operations"`
	Message    string            `json:"messae"`
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
	Int32     DataType = iota // value --> 0
	Int64                     // value --> 1
	Float                     // value --> 2
	Double                    // value --> 3
	Enum                      // value --> 4
	Bool                      // value --> 5
	Text                      // value --> 6
	Date                      // value --> 7
	Timestamp                 // value --> 8
	Struct                    // value --> 9
	Array                     // value --> 10
	File                      // value --> 11
	Password                  // value --> 12
)

// int32：32位整型。需定义取值范围、步长和单位符号。
// int64：32位整型。需定义取值范围、步长和单位符号。
// float：单精度浮点型。需定义取值范围、步长和单位符号。
// double：双精度浮点型。需定义取值范围、步长和单位符号。
// enum：枚举型。定义枚举项的参数值和参数描述，例如：1表示加热模式、2表示制冷模式。
// bool：布尔型。采用0或1来定义布尔值，例如：0表示关、1表示开。
// text：字符串。需定义字符串的数据长度，最长支持10240字节。
// date：时间戳。格式为String类型的UTC时间戳，单位：毫秒。
// struct：JSON对象。定义一个JSON结构体，新增JSON参数项，例如：定义灯的颜色是由Red、Green、Blue三个参数组成的结构体。不支持结构体嵌套。
// array：数组。需声明数组内的元素类型、数组元素个数。元素类型可选择int32、float、double、text或struct，需确保同一个数组元素类型相同。元素个数，限制1~512个。
// file（文件，支持URL[地址]/base64[base64编码]/binary[二进制]）
// password（密码）"
