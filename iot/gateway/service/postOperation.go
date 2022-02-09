package service

import (
	"fmt"
	"time"

	"koudai-box/iot/db"
	"koudai-box/iot/gateway/model"
	status "koudai-box/iot/gateway/status"
	"koudai-box/iot/gateway/utils"

	"github.com/sirupsen/logrus"
)

func PostOperation(deviceId int, operationCode string) error {
	device, ok := GetDeviceById(deviceId, model.STATUS_ACTIVE)
	if !ok {
		logrus.Errorf("PostOperation error,device is null or not start,deviceId[%d]", deviceId)
		return fmt.Errorf("PostOperation error,device is null or not start,deviceId[%d]", deviceId)
	}
	if nil == device.Product {
		logrus.Errorf("PostOperation error,product is null or not start,deviceId[%d]", deviceId)
		return fmt.Errorf("PostOperation error,product is null or not start,deviceId[%d]", deviceId)
	}
	product, ok := GetProduct(device.Product.Id, model.STATUS_ACTIVE)
	if !ok {
		logrus.Errorf("PostOperation error,product is null or not start,productId[%d]", device.Product.Id)
		return fmt.Errorf("PostOperation error,product is null or not start,productId[%d]", device.Product.Id)
	}
	var config *model.OperationConfig
	for _, op := range product.OperationConfigs {
		if operationCode == op.Code {
			config = &op
			break
		}
	}
	if nil == config {
		logrus.Errorf("PostOperation error,OperationConfig is null,OperationConfig[%s]", operationCode)
		return fmt.Errorf("PostOperation error,OperationConfig is null,OperationConfig[%s]", operationCode)
	}
	gatewayId := device.Product.GatewayId
	dri, ok := status.GetDriver(gatewayId)
	if !ok {
		logrus.Errorf("PostOperation error,gateway is null or not start,GatewayId[%d]", gatewayId)
		return fmt.Errorf("PostOperation error,gateway is null or not start,GatewayId[%d]", gatewayId)
	}
	props, _ := status.GetDeviceLastProp(deviceId)
	data := make(map[string]interface{})
	if len(config.Inputs) > 0 {
		for _, v := range config.Inputs {
			if len(v.Value) != 0 {
				data[v.Key] = v.Value
				continue
			}
			value, ok := props.Properties[v.Key]
			if ok {
				data[v.Key] = value
			}
		}
	}
	gateway := dri.GetGatewayConfig()
	var apiConfig *model.ApiConfig
	for _, api := range gateway.ApiConfigs {
		if api.Name == config.Router {
			apiConfig = &api
		}
	}
	if nil == apiConfig {
		logrus.Errorf("PostOperation error,apiConfig is not exist,router[%s]", config.Router)
		return fmt.Errorf("PostOperation error,apiConfig is not exist,router[%s]", config.Router)
	}
	resp, err := dri.PostOperation(*apiConfig, data, device)
	out := ""
	if nil != err {
		logrus.Errorf("PostOperation error,", deviceId)
		out = fmt.Sprintf("PostOperation error,%+v", err)
	} else {
		out = utils.ToString(resp)
	}
	db.InsertOperationRecord(db.OperationRecord{
		GatewayId:  gatewayId,
		ProductId:  device.Product.Id,
		DeviceId:   deviceId,
		Code:       config.Code,
		Name:       config.Name,
		CreateTime: time.Now(),
		Type:       config.Type,
		Desc:       config.Desc,
		Inputs:     utils.ToString(config.Inputs),
		Outputs:    out})
	return nil
}
