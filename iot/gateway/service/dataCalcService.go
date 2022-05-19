package service

import (
	"errors"
	"strconv"
	"time"

	"koudai-box/iot/gateway/model"
	status "koudai-box/iot/gateway/status"
	"koudai-box/iot/gateway/utils"

	orm "koudai-box/iot/db"

	"github.com/sirupsen/logrus"
)

type DataGatewayApi interface {
	//数据计算接口
	Calculater(data interface{}) (interface{}, error)
	//告警过滤接口
	FilterAlarm(data model.PropertyMessage) (interface{}, error)
	//存储属性
	LoaderProperty(data model.PropertyMessage, push bool) (bool, error)
	//存储告警
	LoaderAlarm(data model.IotEventMessage, push bool) (bool, error)
	//存储事件
	LoadeEvent(data model.EventMessage, push bool) (bool, error)
	//数据存储接口
	Push(interface{}, string) bool
}

type DataGateway struct {
	Device *model.Device
}

/**
 *数据计算
 */
func (gateway *DataGateway) Calculater(data interface{}) (interface{}, error) {
	device := gateway.Device
	if len(device.Product.FunctionConfigs) == 0 {
		logrus.Errorf("Function_Calc is null,device key=%s", device.Key)
		return data, nil
	}
	function, ok := device.Product.FunctionConfigs[model.Function_Calc]
	if !ok || len(function.Function) == 0 {
		logrus.Errorf("Function_Calc is null,device key=%s", device.Key)
		return data, nil
	}
	logrus.Debugf("device[%s]'s Function_Calc funtion name=%s", gateway.Device.Key, model.Function_Calc)
	deviceStatus := status.GetDeviceStatus(device.Id)
	context := make(map[string]interface{})
	context["device"] = device
	context["status"] = deviceStatus
	return utils.ExecJSWithContext(function.Function, function.Key, context, data)
}

const js_calc = `
function calculate(data) {
	console.log("calculate context:",JSON.stringify(context))
	var tmp = data //JSON.parse(JSON.stringify(data))
	console.log("zhucz32 properties:",tmp.properties.X_B)
	console.log("zhucz32 Y_X:",tmp.properties.Y_X.value)
	console.log("zhucz32 Y_Y:",tmp.properties.Y_Y.value)
	console.log("zhucz32 Y_Z:",tmp.properties.Y_Z.value)
	var RealX = getX(tmp.properties.Y_X.value,tmp.properties.Y_Y.value,tmp.properties.Y_Z.value)
	var RealY = getY(tmp.properties.Y_X.value,tmp.properties.Y_Y.value,tmp.properties.Y_Z.value)
	var RealZ = getZ(tmp.properties.Y_X.value,tmp.properties.Y_Y.value,tmp.properties.Y_Z.value)
	console.log("zhucz32 RealX:",RealX)
	console.log("zhucz32 RealY:",RealY)
	console.log("zhucz32 RealZ:",RealZ)
	tmp.properties.X_B.value = RealX
	tmp.properties.Y_B.value = RealY
	tmp.properties.Z_B.value = RealZ
	console.log("zhucz33:",tmp.properties.X_B.value)
	return tmp;
}

// RealX = atan(AX/(sqrt(AY*AY+AZ*AZ)));
// RealY = atan(AY/(sqrt(AX*AX+AZ*AZ)));
// RealY = atan((sqrt(AX*AX+AY*AY))/AZ);
function getX(AX, AY, AZ) {
	console.log("zhucz32 getX:",AX,AY,AZ)
	var RealX = Math.atan(AX / Math.sqrt(Math.pow(AY, 2)+(Math.pow(AZ, 2))))
	console.log("zhucz32 getX RealX:",RealX)
	return RealX
}

function getY(AX, AY, AZ) {
	var RealY = Math.atan(AY / Math.sqrt(Math.pow(AX, 2)+(Math.pow(AZ, 2))))
	return RealY
}

function getZ(AX, AY, AZ) {
	var RealZ = Math.atan(Math.sqrt(Math.pow(AX, 2)+(Math.pow(AY, 2))) / AZ)
	return RealZ
}

`

/**
 *告警过滤
 */
func (gateway *DataGateway) FilterAlarm(data model.PropertyMessage) ([]interface{}, error) {
	logrus.Debugf("FilterAlarm,message id = %s", data.MessageId)
	device := gateway.Device
	if len(device.Product.AlarmConfigs) == 0 {
		logrus.Infof("AlarmConfigs is null,device key=%s", device.Key)
		return nil, errors.New("AlarmConfigs is null")
	}
	alarms := make([]interface{}, 0)
	for _, config := range device.Product.AlarmConfigs {
		if len(config.Conditions) == 0 {
			logrus.Infof("Conditions is null,device key=%s,alarm key=%s", device.Key, config.Key)
			continue
		}
		hasAlarm := false
		var poprs []model.PropertyItem
		for _, condition := range config.Conditions {
			match := utils.MatchContidion(data, condition)
			hasAlarm = hasAlarm || match
			if match {
				poprs = append(poprs, data.Properties[condition.Key])
			}
		}
		if hasAlarm {
			alarms = append(alarms, model.IotEventMessage{DeviceId: strconv.Itoa(device.Id), DeviceSign: device.Key, MessageId: utils.GetUUID(), Code: config.Code,
				Timestamp: time.Now().Unix(), Type: config.Type, Level: config.Level, Title: config.Name, Message: config.Message, Properties: poprs, Conditions: config.Conditions})
		}
	}
	return alarms, nil

}

/**
 *消息存储
 */
func (gateway *DataGateway) LoaderProperty(data model.PropertyMessage, push bool) (bool, error) {
	logrus.Debugf("save prop,message id = %s", data.MessageId)
	prop := orm.DeviceProperty{
		DeviceId:   data.DeviceId,
		MessageId:  data.MessageId,
		Timestamp:  data.Timestamp,
		Properties: utils.ToString(data.Properties),
		CreateTime: time.Now(),
		PushFlag:   0,
	}
	if _, err := orm.InsertProperty(prop); err != nil {
		return false, err
	}
	if push {
		return Push(data, model.Message_Type_Prop), nil
	}
	//属性推送
	return true, nil

}

/**
 *告警存储
 */
func (gateway *DataGateway) LoaderAlarm(data model.IotEventMessage, push bool) (bool, error) {
	logrus.Debugf("save alarm,message id = %s", data.MessageId)
	alarm := orm.Alarm{
		ProductId:   gateway.Device.Product.Id,
		ProductName: gateway.Device.Product.Name,
		DeviceId:    data.DeviceId,
		DeviceName:  gateway.Device.Name,
		DeviceSign:  gateway.Device.Key,
		MessageId:   data.MessageId,
		Timestamp:   data.Timestamp,
		Code:        data.Code,
		Title:       data.Title,
		Type:        data.Type,
		Level:       data.Level,
		Message:     data.Message,
		Properties:  utils.ToString(data.Properties),
		Conditions:  utils.ToString(data.Conditions),
		CreateTime:  time.Now(),
		PushFlag:    0,
	}
	if _, err := orm.InsertAlarm(alarm); err != nil {
		return false, err
	}
	if push {
		return Push(data, model.Message_Type_Iot_Event), nil
	}
	return true, nil
}

/**
 *事件存储
 */
func (gateway *DataGateway) LoadeEvent(data model.EventMessage, push bool) (bool, error) {
	logrus.Debugf("save event,message id = %s", data.MessageId)
	prop := orm.Event{
		ProductId:   gateway.Device.Product.Id,
		ProductName: gateway.Device.Product.Name,
		DeviceId:    data.DeviceId,
		DeviceName:  gateway.Device.Name,
		DeviceSign:  gateway.Device.Key,
		MessageId:   data.MessageId,
		Timestamp:   data.Timestamp,
		Type:        data.Type,
		Level:       "",
		Title:       data.Title,
		Message:     data.Message,
		Properties:  utils.ToString(data.Properties),
		CreateTime:  time.Now(),
		PushFlag:    0,
	}
	if _, err := orm.InsertEvent(prop); err != nil {
		return false, err
	}
	if push {
		return Push(data, model.Message_Type_Event), nil
	}
	return true, nil

}

/**
 *数据推送
 */
func Push(data interface{}, msgType string) bool {
	// return Public(data, router)
	var msgId string
	switch msgType {
	case model.Message_Type_Event:
		msgId = data.(model.EventMessage).MessageId
	case model.Message_Type_Iot_Event:
		msgId = data.(model.IotEventMessage).MessageId
	case model.Message_Type_Prop:
		msgId = data.(model.PropertyMessage).MessageId

	}
	model.PushOutMsgChan <- model.Message{
		SN:    utils.GetSN(),
		Type:  msgType,
		Msg:   data,
		MsgId: msgId,
	}
	return true
}
