package dto

import (
	"mime/multipart"

	"koudai-box/iot/gateway/model"
)

//添加产品列表
type AddProductDataRequest struct {
	Name      string                `json:"name" form:"name" binding:"required"`
	Code      string                `json:"code" form:"code" binding:"required"`
	Image     *multipart.FileHeader `json:"image" form:"image" binding:"required"`
	Category  string                `json:"category" form:"category" binding:"required"`
	GatewayId int                   `json:"gatewayId" form:"gatewayId" binding:"required"`
	Desc      string                `json:"desc" form:"desc" binding:""`
}

//更新产品列表
type UpdateProductDataRequest struct {
	Id        int                   `json:"id" form:"id" binding:"required"`
	Name      string                `json:"name" form:"name" binding:"required"`
	Code      string                `json:"code" form:"code" binding:"required"`
	Image     *multipart.FileHeader `json:"image" form:"image" binding:""`
	Category  string                `json:"category" form:"category" binding:"required"`
	GatewayId int                   `json:"gatewayId" form:"gatewayId" binding:"required"`
	Desc      string                `json:"desc" form:"desc" binding:""`
}

//设置产品状态
type SetProductStateRequest struct {
	Id    int  `json:"id" form:"id" binding:"required"`
	State *int `json:"state" form:"state"  binding:"required"`
}

//查询列表
type QueryProductDataRequest struct {
	Search   string `json:"search" form:"search"`
	State    int    `json:"state" form:"state"`
	PageNo   int    `json:"pageNo" form:"pageNo"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}

//添加产品物模型
type AddProductItemDataRequest struct {
	Sort           int            `json:"sort"`
	ProductId      int            `json:"productId" form:"productId" binding:"required"` //产品ID
	Code           string         `json:"code" form:"code" binding:"required"`           //属性标识
	SourceCode     string         `json:"sourceCode" form:"sourceCode" binding:""`       //源属性标识
	NodeId         string         `json:"nodeId" form:"NodeId" binding:""`               //opc ua 节点
	Address        string         `json:"address" form:"address" binding:""`             //modbus，地址
	Quantity       int            `json:"quantity" form:"quantity" binding:""`           //modbus，quantity 数量
	OperaterType   string         `json:"operaterType" form:"OperaterType" binding:""`   //modbus，操作类型
	ZoomFactor     float32        `json:"zoomFactor" form:"ZoomFactor" binding:""`       //modbus，缩放因子：不能为0，默认为1
	ExchangeHL     bool           `json:"exchangeHL" form:"ExchangeHL" binding:""`       //modbus，交换寄存器内高低字节：true：互换，false：不互换
	ExchangeOrder  bool           `json:"exchangeOrder" form:"ExchangeOrder" binding:""` //modbus，交换寄存器顺序：true：互换，false：不互换
	Name           string         `json:"name" form:"name" binding:"required"`           //属性名称
	RW             string         `json:"rw" form:"rw" binding:"required"`               //读写标志
	DataType       model.DataType `json:"dataType" form:"dataType" binding:""`           //数据类型
	Unit           string         `json:"unit" form:"unit" binding:""`                   //单位
	Min            string         `json:"min" form:"min" binding:""`                     //最小值
	Max            string         `json:"max" form:"max" binding:""`
	Precision      int            `json:"precision"`                       //精度
	MaxLength      int            `json:"maxLength"`                       //最大长度
	BoolValue      string         `json:"boolValue"`                       //布尔值
	DateFormat     string         `json:"dateFormat"`                      //时间
	Dict           string         `json:"dict"`                            //枚举
	FileType       string         `json:"fileType"`                        // 文件类型
	PasswordLength string         `json:"passwordLength"`                  //密码长度
	Report         string         `json:"report" form:"report" binding:""` //上传方式
	Step           string         `json:"step" form:"step" binding:""`     //变化步长
	Desc           string         `json:"desc" form:"desc" binding:""`     //描述
}

//更新产品物模型
type UpdateProductItemDataRequest struct {
	Sort           int            `json:"sort"`
	ProductId      int            `json:"productId" form:"productId" binding:"required"` //产品ID
	Key            string         `json:"key" form:"key" binding:"required"`             //key
	Code           string         `json:"code" form:"code" binding:"required"`           //属性标识
	SourceCode     string         `json:"sourceCode" form:"sourceCode" binding:""`       //源属性标识
	NodeId         string         `json:"nodeId" form:"NodeId" binding:""`               //opc ua 节点
	Address        string         `json:"address" form:"address" binding:""`             //modbus，地址
	Quantity       int            `json:"quantity" form:"quantity" binding:""`           //modbus，quantity 数量
	OperaterType   string         `json:"operaterType" form:"OperaterType" binding:""`   //modbus，操作类型
	ZoomFactor     float32        `json:"zoomFactor" form:"ZoomFactor" binding:""`       //modbus，缩放因子：不能为0，默认为1
	ExchangeHL     bool           `json:"exchangeHL" form:"ExchangeHL" binding:""`       //modbus，交换寄存器内高低字节：true：互换，false：不互换
	ExchangeOrder  bool           `json:"exchangeOrder" form:"ExchangeOrder" binding:""` //modbus，交换寄存器顺序：true：互换，false：不互换
	Name           string         `json:"name" form:"name" binding:"required"`           //属性名称
	RW             string         `json:"rw" form:"rw" binding:"required"`               //读写标志
	DataType       model.DataType `json:"dataType" form:"dataType" binding:""`           //数据类型
	Unit           string         `json:"unit" form:"unit" binding:""`                   //单位
	Min            string         `json:"min" form:"min" binding:""`                     //最小值
	Max            string         `json:"max" form:"max" binding:""`                     //最大值
	Precision      int            `json:"precision"`                                     //精度
	MaxLength      int            `json:"maxLength"`                                     //最大长度
	BoolValue      string         `json:"boolValue"`                                     //布尔值
	DateFormat     string         `json:"dateFormat"`                                    //时间
	Dict           string         `json:"dict"`                                          //枚举
	FileType       string         `json:"fileType"`                                      //文件类型
	PasswordLength string         `json:"passwordLength"`                                //密码长度
	Report         string         `json:"report" form:"report" binding:""`               //上传方式
	Step           string         `json:"step" form:"step" binding:""`                   //变化步长
	Desc           string         `json:"desc" form:"desc" binding:""`                   //描述
}

/**物模型定义**/
type ProductItemConfig struct {
	Sort          int                `json:"sort"`
	Key           string             `json:"key"`
	Name          string             `json:"name"`
	Code          string             `json:"code"`
	SourceCode    string             `json:"sourceCode"`
	NodeId        string             `json:"nodeId"`        //opc ua 节点
	Address       string             `json:"address"`       //modbus，地址
	Quantity      int                `json:"quantity"`      //modbus，quantity 数量
	OperaterType  string             `json:"operaterType"`  //modbus，操作类型
	ZoomFactor    float32            `json:"zoomFactor"`    //modbus，缩放因子：不能为0，默认为1
	ExchangeHL    bool               `json:"exchangeHL"`    //modbus，交换寄存器内高低字节：true：互换，false：不互换
	ExchangeOrder bool               `json:"exchangeOrder"` //modbus，交换寄存器顺序：true：互换，false：不互换
	DataType      model.ItemDataType `json:"dataType"`      //数据类型
	Report        string             `json:"report"`
	Desc          string             `json:"desc"`
}

//添加产品操作物模型
type AddProductOperationDataRequest struct {
	ProductId int               `json:"productId" form:"productId" binding:"required"` //产品ID
	Code      string            `json:"code" form:"code" binding:"required"`           //属性标识
	Router    string            `json:"router" form:"router" binding:"required"`       //属性标识
	Name      string            `json:"name" form:"name" binding:"required"`           //属性名称
	Type      string            `json:"type" form:"type" binding:"required"`           //异步同步
	Desc      string            `json:"desc" form:"desc" binding:""`                   //描述
	Inputs    []model.Parameter `json:"inputs" form:"inputs" binding:""`               //入口参数
	Outputs   []model.Parameter `json:"outputs" form:"outputs" binding:""`             //出口参数
}

//添加产品操作物模型
type UpdateProductOperationDataRequest struct {
	ProductId int               `json:"productId" form:"productId" binding:"required"` //产品ID
	Key       string            `json:"key" form:"key" binding:"required"`             //
	Code      string            `json:"code" form:"code" binding:"required"`           //属性标识
	Router    string            `json:"router" form:"router" binding:"required"`       //属性标识
	Name      string            `json:"name" form:"name" binding:"required"`           //属性名称
	Type      string            `json:"type" form:"type" binding:"required"`           //异步同步
	Desc      string            `json:"desc" form:"desc" binding:""`                   //描述
	Inputs    []model.Parameter `json:"inputs" form:"inputs" binding:""`
	Outputs   []model.Parameter `json:"outputs" form:"outputs" binding:""` //入口参数
}

/**操作定义**/
type ProductApiConfig struct {
	Key             string            `json:"key"`
	Code            string            `json:"code"`
	Name            string            `json:"name"`
	Path            string            `json:"path"`
	Method          string            `json:"mthod"`           //方式，get、post
	Parameters      []model.Parameter `json:"parameters"`      //入参配置
	CollectType     string            `json:"collectType"`     //采集方式：定时、轮询
	CollectPeriod   int               `json:"collectPeriod"`   //采集周期:秒
	DataCombination string            `json:"dataCombination"` //数据类型，单条、数组
	Cron            string            `json:"cron"`            //时间表达式
}

//添加产品事件
type AddProductEventDataRequest struct {
	ProductId int               `json:"productId" form:"productId" binding:"required"` //产品ID
	Code      string            `json:"code" form:"code" binding:"required"`           //属性标识
	Name      string            `json:"name" form:"name" binding:"required"`           //属性名称
	Type      string            `json:"type" form:"type" binding:"required"`           //信息,告警,故障
	Desc      string            `json:"desc" form:"desc" binding:""`                   //描述
	Outputs   []model.Parameter `json:"outputs" form:"outputs" binding:""`             //入口参数
}

type UpdateProductEventDataRequest struct {
	ProductId int               `json:"productId" form:"productId" binding:"required"` //产品ID
	Key       string            `json:"key" form:"key" binding:"required"`             //
	Code      string            `json:"code" form:"code" binding:"required"`           //属性标识
	Name      string            `json:"name" form:"name" binding:"required"`           //属性名称
	Type      string            `json:"type" form:"type" binding:"required"`           //信息,告警,故障
	Desc      string            `json:"desc" form:"desc" binding:""`                   //描述
	Outputs   []model.Parameter `json:"outputs" form:"outputs" binding:""`             //入口参数
}

type ProductEventConfig struct {
	Key     string            `json:"key"`
	Code    string            `json:"code"`
	Name    string            `json:"name"`
	Type    string            `json:"type"`    //异步同步
	Outputs []model.Parameter `json:"outputs"` //接口配置
	Desc    string            `json:"desc"`
}

/**告警定义**/
type AddProductAlarmDataRequest struct {
	ProductId  int                     `json:"productId" form:"productId" binding:"required"` //产品ID
	Name       string                  `json:"name" form:"name" binding:"required"`
	Code       string                  `json:"code" form:"code" binding:"required"`
	Level      string                  `json:"level" form:"level" binding:"required"`   //级别
	Type       string                  `json:"type" form:"type" binding:"required"`     //类型，事件、告警
	Conditions []ProductCondition      `json:"conditions" form:"conditions" binding:""` //条件
	Operations []model.OperationConfig `json:"operations" form:"operations" binding:""` //操作
	Message    string                  `json:"message" form:"message" binding:""`
	Desc       string                  `json:"desc" form:"desc" binding:""`
	State      int                     `json:"state" form:"state" binding:""`
}

type UpdateProductAlarmDataRequest struct {
	ProductId  int                     `json:"productId" form:"productId" binding:"required"` //产品ID
	Key        string                  `json:"key" form:"key" binding:"required"`             //
	Name       string                  `json:"name" form:"name" binding:"required"`
	Code       string                  `json:"code" form:"code" binding:"required"`
	Level      string                  `json:"level" form:"level" binding:"required"`   //级别
	Type       string                  `json:"type" form:"type" binding:"required"`     //类型，事件、告警
	Conditions []ProductCondition      `json:"conditions" form:"conditions" binding:""` //条件
	Operations []model.OperationConfig `json:"operations" form:"operations" binding:""` //操作
	Message    string                  `json:"message" form:"message" binding:""`
	Desc       string                  `json:"desc" form:"desc" binding:""`
	State      int                     `json:"state" form:"state" binding:""`
}

type ProductAlarmConfig struct {
	Key        string                  `json:"key"`
	Name       string                  `json:"name"`
	Code       string                  `json:"code"`
	Level      string                  `json:"level"`      //级别
	Type       string                  `json:"type"`       //类型，事件、告警
	Conditions []ProductCondition      `json:"conditions"` //条件
	Operations []model.OperationConfig `json:"operations"` //操作
	Message    string                  `json:"messae"`
	CreateTime string                  `json:"createTime"`
	Desc       string                  `json:"desc"`
	State      int                     `json:"state"`
}

type ProductCondition model.Condition

//函数
type AddProductFunctionDataRequest struct {
	ProductId    int    `json:"productId" form:"productId" binding:"required"` //产品ID
	ExtractProp  string `json:"extractProp"  form:"extractProp"`
	ExtractEvent string `json:"extractEvent"  form:"extractEvent"`
	Calculate    string `json:"calculate"  form:"calculate"`
	Body         string `json:"body"  form:"body"`
}

type UpdateProductFunctionDataRequest struct {
	ProductId    int    `json:"productId" form:"productId" binding:"required"` //产品ID
	Key          string `json:"key" form:"key" binding:"required"`
	ExtractProp  string `json:"extractProp"  form:"extractProp"`
	ExtractEvent string `json:"extractEvent"  form:"extractEvent"`
	Calculate    string `json:"calculate"  form:"calculate"`
	Body         string `json:"body"  form:"body"`
}

type ProductFunctionConfig struct {
	Key          string `json:"key" form:"key"` //主键
	ExtractProp  string `json:"extractProp"  form:"extractProp"`
	ExtractEvent string `json:"extractEvent"  form:"extractEvent"`
	Calculate    string `json:"calculate"  form:"calculate"`
	Body         string `json:"body"  form:"body"`
}
