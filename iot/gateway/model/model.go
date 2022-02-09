package model

//协议
const (
	Geteway_Protocol_HTTP_Client      = "http_client"
	Geteway_Protocol_HTTP_Server      = "http_server"
	Geteway_Protocol_WebSocket_Server = "Websocket_server"
	Geteway_Protocol_WebSocket_Client = "Websocket_client"
	Geteway_Protocol_ModbusTCP        = "Modbus_TCP"
	Geteway_Protocol_ModbusRTU        = "Modbus_RTU"
	Geteway_Protocol_MQTT             = "MQTT"
	Geteway_Protocol_MQTTSN           = "MQTT-SN"
	Geteway_Protocol_OPCUA            = "OPC_UA"
	Geteway_Protocol_BACnet_IP        = "BACnet_IP"
	Geteway_Protocol_CoAP             = "CoAP"
	Geteway_Protocol_LwM2M            = "LwM2M"
)
const (
	API_GetProp  = "获取属性" //取数据接口
	API_GetEvent = "获取事件" //写数据接口
)
const (
	Function_Extract_Prop  = "extractProp"  //数据抽取函数
	Function_Calc          = "calculate"    //计算函数
	Function_postBody      = "postBody"     //数据写入前处理函数
	Function_Extract_Event = "extractEvent" //事件抽取函数
)
const (
	DataReportType_Schedule = "schedule" //按时上报
	DataReportType_Change   = "change"   //变化上报
)
const (
	CollectType_Schedule = "schedule" //定时
	CollectType_Poll     = "poll"     //轮询
)
const (
	DataCombination_Single = "single" //单条
	DataCombination_Array  = "array"  //批量
)
const (
	STATUS_ALL       = -1 //状态：不生效，全部状态
	STATUS_ACTIVE    = 1  //状态：启用、激活
	STATUS_DISACTIVE = 0  //状态：停用、非激活
	STATUS_SUCCESS   = 1
	STATUS_ERROR     = 0
)

const (
	Msg_Type_Props  = "props"
	Msg_Type_Events = "events"
)

/**********************************消息通道************************************************/
type PropertyChan struct {
	PropertyMessage PropertyMessage
	Device          *Device
}

type EventChan struct {
	EventMessage EventMessage
	Device       *Device
}

type PushMsg struct {
	Msg        interface{}
	DeviceKey  string
	GatewayKey string
	Type       string //Msg_Type
}

type StatusMsg struct {
	DeviceId int
	Status   int
}

const Message_queen_size = 64

//设备属性消息通道
var PropMessChan = make(chan PropertyChan, Message_queen_size)
var EventMessChan = make(chan EventChan, Message_queen_size)
var PushMsgChan = make(chan PushMsg, Message_queen_size)

var CheckChan = make(chan struct{}, Message_queen_size)
var StatusMsgChan = make(chan StatusMsg, Message_queen_size)

var PushOutMsgChan = make(chan Message, Message_queen_size)
var PushStatusChan = make(chan struct{}, Message_queen_size)

/**********************************实体定义************************************************/
type Device struct {
	Id       int                    `json:"id"`
	Key      string                 `json:"key"`
	Name     string                 `json:"name"`
	SourceId string                 `json:"sourceId"`
	Geo      [2]float32             `json:"geo"`
	Product  *Product               `json:"product"`
	ExtProps map[string]interface{} `json:"extProps"` //扩展属性
	Desc     string                 `json:"desc"`
}

type Product struct {
	Id               int                       `json:"id"`
	Key              string                    `json:"key"`
	Name             string                    `json:"name"`
	Category         string                    `json:"category"` //分类
	Desc             string                    `json:"desc"`
	Items            []ItemConfig              `json:"items"`            //物模型
	EventConfigs     []string                  `json:"eventConfigs"`     //函数配置
	OperationConfigs []OperationConfig         `json:"operationConfigs"` //操作定义
	AlarmConfigs     []AlarmConfig             `json:"alarmSettings"`    //告警配置
	FunctionConfigs  map[string]FunctionConfig `json:"functionConfigs"`  //函数配置
	GatewayId        int                       `json:"gatewayId"`
}

/**物模型定义**/
type ItemConfig struct {
	Key            string       `json:"key"`
	Code           string       `json:"code"`
	Name           string       `json:"name"`
	OperaterType   string       `json:"operaterType"`  //modbus，操作类型
	Address        string       `json:"address"`       //modbus，地址
	Quantity       int          `json:"quantity"`      //modbus，quantity 数量
	Source         string       `json:"source"`        //源数据标识
	NodeId         string       `json:"nodeId"`        //opc ua 节点
	DataType       ItemDataType `json:"dataType"`      //数据类型
	ZoomFactor     float32      `json:"zoomFactor"`    //modbus，缩放因子：不能为0，默认为1
	ExchangeHL     bool         `json:"exchangeHL"`    //modbus，交换寄存器内高低字节：true：互换，false：不互换
	ExchangeOrder  bool         `json:"exchangeOrder"` //modbus，交换寄存器顺序：true：互换，false：不互换
	DataReportType string       `json:"dataReportType"`
	Desc           string       `json:"desc"`
}

type ItemDataType struct {
	RW             string            `json:"rw"`             //读写标志
	Type           DataType          `json:"type"`           //数据类型
	Unit           string            `json:"unit"`           //单位
	Min            string            `json:"min"`            //最小值
	Max            string            `json:"max"`            //最大值
	Step           string            `json:"step"`           //变化步长
	Precision      int               `json:"precision"`      //精度
	MaxLength      int               `json:"maxLength"`      //最大长度
	BoolValue      map[string]string `json:"boolValue"`      //布尔值
	DateFormat     string            `json:"dateFormat"`     //时间格式
	Dict           map[string]string `json:"dict"`           //枚举项
	FileType       string            `json:"fileType"`       // 0: url 1: base64 2: []byte
	PasswordLength string            `json:"passwordLength"` //密码长度
}

/**操作定义**/
type OperationConfig struct {
	Key     string      `json:"key"`
	Code    string      `json:"code"`
	Name    string      `json:"name"`
	Type    string      `json:"type"`    //异步同步
	Inputs  []Parameter `json:"inputs"`  //输入参数
	Outputs []Parameter `json:"outputs"` //输出参数
	Router  string      `json:"router"`  //接口配置
	Desc    string      `json:"desc"`
}

type Parameter struct {
	Key      string   `json:"key"`
	Name     string   `json:"name"`
	DataType DataType `json:"dataType"` //数据类型
	Value    string   `json:"value"`
}

type ApiConfig struct {
	Name            string      `json:"name"`
	Path            string      `json:"path"`
	Method          string      `json:"mthod"`           //方式，get、post
	Parameters      []Parameter `json:"parameters"`      //入参配置
	CollectType     string      `json:"collectType"`     //采集方式：定时、轮询
	CollectPeriod   int         `json:"collectPeriod"`   //采集周期:秒
	DataCombination string      `json:"dataCombination"` //数据类型，单条、数组
	Cron            string      `json:"cron"`            //时间表达式
}

/**告警定义**/
type AlarmConfig struct {
	Key        string            `json:"key"`
	Code       string            `json:"code"`
	Name       string            `json:"name"`
	Level      string            `json:"level"`      //级别
	Type       string            `json:"type"`       //类型，事件、告警
	Conditions []Condition       `json:"conditions"` //条件
	Operations []OperationConfig `json:"operations"` //操作
	Message    string            `json:"messae"`
}
type Condition struct {
	Key      string      `json:"key"`
	Name     string      `json:"name"`
	DataType DataType    `json:"dataType"` //数据类型
	Compare  string      `json:"compare"`
	Value    interface{} `json:"value"`
}
type FunctionConfig struct {
	Key      string `json:"key"`      //主键
	Name     string `json:"name"`     //名称
	Function string `json:"function"` //函数定义
}

type GatewayConfig struct {
	Id            int                  `json:"id"`
	Key           string               `json:"key"`
	Name          string               `json:"name"`
	Protocol      string               `json:"protocol"`      //协议，http、modbus、opc ua、mqtt、websocket、。。。
	Ip            string               `json:"ip"`            //ip
	Port          int                  `json:"port"`          //端口
	Parameters    []Parameter          `json:"parameters"`    //入参配置
	ApiConfigs    map[string]ApiConfig `json:"apiConfigs"`    //接口配置
	CollectType   string               `json:"collectType"`   //modbus\opcua\bacnet采集方式：定时、轮询
	CollectPeriod int                  `json:"collectPeriod"` //modbus\opcua\bacnet采集周期:秒
	Cron          string               `json:"cron"`          //modbus\opcua\bacnet时间表达式
	ModbusConfig  ModbusConfig         `json:"modbusConfig"`  //modbus
	Desc          string               `json:"desc"`
}

type ModbusConfig struct {
	Com      string `json:"com"`      //'串口'，modbus
	BaudRate int    `json:"baudRate"` //modbus，默认 9600
	DataBits int    `json:"dataBits"` //modbus，默认 8
	Parity   string `json:"parity"`   //modbus，默认 N
	StopBits int    `json:"stopBits"` //modbus，默认 1
}

type DeviceStatus struct {
	Id           int              `json:"id"`
	Key          string           `json:"key"`
	Status       int              `json:"status"` //状态
	Timestamp    int64            `json:"timestamp"`
	LastStatus   *PropertyMessage `json:"lastStatus"`   //最新
	ZeroStatus   *PropertyMessage `json:"zeroStatus"`   //初始值
	PreDayStatus *PropertyMessage `json:"preDayStatus"` //前一天平均
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
	Long                      // value --> 13
	GeoPoint                  // value --> 14
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
