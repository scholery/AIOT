package dto

type AddDeviceRequest struct {
	Code      string                 `json:"code" form:"code" binding:"required"`
	Name      string                 `json:"name" form:"name" binding:"required"`
	SourceId  string                 `json:"sourceId"`                                      //设备源标识，modbus：slaveid，opc ua：namespace。json：property key
	Geo       string                 `json:"geo"`                                           //经纬度
	Locale    string                 `json:"locale"`                                        //位置描述
	ProductId int                    `json:"productId" form:"productId" binding:"required"` //event、alarm
	Desc      string                 `json:"desc" form:"desc" binding:""`
	ExtProps  map[string]interface{} `json:"extProps" form:"extProps" binding:""` //扩展属性
}

type UpdateDeviceRequest struct {
	Id        int                    `json:"id" form:"id" binding:"required"`
	Code      string                 `json:"code" form:"code" binding:"required"`
	Name      string                 `json:"name" form:"name" binding:"required"`
	SourceId  string                 `json:"sourceId"`                                      //设备源标识，modbus：slaveid，opc ua：namespace。json：property key
	Geo       string                 `json:"geo"`                                           //经纬度
	Locale    string                 `json:"locale"`                                        //位置描述
	ProductId int                    `json:"productId" form:"productId" binding:"required"` //event、alarm
	Desc      string                 `json:"desc" form:"desc" binding:""`
	ExtProps  map[string]interface{} `json:"extProps" form:"extProps" binding:""` //扩展属性
}

type QueryDeviceDataRequest struct {
	ProductId      int    `json:"productId" form:"productId"`
	Search         string `json:"search" form:"search"`
	ActivateStatus int    `json:"activateStatus" form:"activateStatus"`
	RunningStatus  int    `json:"runningStatus" form:"runningStatus"`
	PageNo         int    `json:"pageNo" form:"pageNo"`
	PageSize       int    `json:"pageSize" form:"pageSize"`
}

//设置设备状态
type SetDeviceActivateRequest struct {
	Id             int  `json:"id" form:"id" binding:"required"`
	ActivateStatus *int `json:"activateStatus" form:"activateStatus"  binding:"required"`
}

type DeviceStateItem struct {
	Id             int  `json:"id"`
	ActivateStatus *int `json:"activateStatus"`
}

//批量设置设备状态
type SetDeviceStatesRequest struct {
	DeviceStates []DeviceStateItem `json:"deviceStates"`
}

type DeleteDeviceRequest struct {
	Ids []int `json:"ids" form:"ids" binding:"required"`
}

type SetBatchDeviceActivateRequest struct {
	Ids            []int `json:"ids" form:"ids" binding:"required"`
	ActivateStatus *int  `json:"activateStatus"`
}

//设置设备状态
type DevicePropertyRequest struct {
	Count int64  `json:"count" form:"count"`
	Begin string `json:"begin" form:"begin" `
	End   string `json:"end" form:"end" `
}
