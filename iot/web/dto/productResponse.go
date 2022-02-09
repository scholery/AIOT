package dto

import "koudai-box/iot/gateway/model"

//列表结构
type ProductItem struct {
	Id                  int      `json:"id"`
	GatewayId           int      `json:"gatewayId"`
	Code                string   `json:"code"`
	Image               string   `json:"image"`
	Name                string   `json:"name"`
	Desc                string   `json:"desc"`
	Category            []string `json:"category"`
	State               int      `json:"state"`
	DeviceCount         int      `json:"deviceCount"`
	DisabledDeviceCount int      `json:"disabledDeviceCount"`
	OnlineDeviceCount   int      `json:"onlineDeviceCount"`
	OfflineDeviceCount  int      `json:"offlineDeviceCount"`
}

//列表结构
type ProductItemConfigItem struct {
	Key            string            `json:"key"`
	Code           string            `json:"code"`
	SourceCode     string            `json:"sourceCode"`
	NodeId         string            `json:"nodeId"`        //opc ua 节点
	Address        string            `json:"address"`       //modbus，地址
	Quantity       int               `json:"quantity"`      //modbus，quantity 数量
	OperaterType   string            `json:"operaterType"`  //modbus，操作类型
	ZoomFactor     float32           `json:"zoomFactor"`    //modbus，缩放因子：不能为0，默认为1
	ExchangeHL     bool              `json:"exchangeHL"`    //modbus，交换寄存器内高低字节：true：互换，false：不互换
	ExchangeOrder  bool              `json:"exchangeOrder"` //modbus，交换寄存器顺序：true：互换，false：不互换
	Name           string            `json:"name"`
	Min            string            `json:"min"`
	Max            string            `json:"max"`
	RW             string            `json:"rw"`
	Step           string            `json:"step"`
	DataType       model.DataType    `json:"dataType"` //数据类型
	Unit           string            `json:"unit"`     //单位
	Precision      int               `json:"precision"`
	MaxLength      int               `json:"maxLength"`
	BoolValue      map[string]string `json:"boolValue"`
	DateFormat     string            `json:"dateFormat"`
	Dict           map[string]string `json:"dict"`
	FileType       string            `json:"fileType"`
	PasswordLength string            `json:"passwordLength"`
	Desc           string            `json:"desc"`
	Report         string            `json:"report"`
}

type ProductOperationConfigItem struct {
	Key     string            `json:"key"`
	Code    string            `json:"code"`
	Name    string            `json:"name"`
	Router  string            `json:"router"`
	Type    string            `json:"type"` //数据类型
	Desc    string            `json:"desc"`
	Inputs  []model.Parameter `json:"inputs" `  //入口参数
	Outputs []model.Parameter `json:"outputs" ` //出口参数
}

type ProductEventConfigItem struct {
	Key     string            `json:"key"`
	Code    string            `json:"code"`
	Name    string            `json:"name"`
	Type    string            `json:"type"` //数据类型
	Desc    string            `json:"desc"`
	Outputs []model.Parameter `json:"outputs" ` //出口参数
}

type ProductAlarmConfigItem struct {
	Key        string             `json:"key"`
	Level      string             `json:"level"`
	Name       string             `json:"name"`
	Code       string             `json:"code"`
	Type       string             `json:"type"` //数据类型
	Message    string             `json:"message"`
	Conditions []ProductCondition `json:"conditions"`
	CreateTime string             `json:"createTime"`
	State      int                `json:"state"`
}

type ProductFunctionConfigItem struct {
	Key          string `json:"key"`
	ExtractProp  string `json:"extractProp"`
	ExtractEvent string `json:"extractEvent"`
	Calculate    string `json:"calculate"`
	Body         string `json:"body"`
}
