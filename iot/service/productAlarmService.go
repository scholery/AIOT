package service

import (
	"encoding/json"
	"errors"
	"time"

	"koudai-box/iot/db"
	"koudai-box/iot/gateway/utils"
	"koudai-box/iot/web/dto"
)

//添加物模型
func AddProductAlarmService(request dto.AddProductAlarmDataRequest) error {
	product, err := db.QueryProductByID(request.ProductId)
	if err != nil {
		return err
	}
	configs := product.AlarmConfigs
	var operationConfigs []dto.ProductAlarmConfig
	err = json.Unmarshal([]byte(configs), &operationConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	//组装数据
	operationConfig := dto.ProductAlarmConfig{
		Key:        utils.GetUUID(),
		Level:      request.Level,
		Name:       request.Name,
		Code:       request.Code,
		Type:       request.Type,
		Conditions: request.Conditions,
		Operations: request.Operations,
		Message:    request.Message,
		Desc:       request.Desc,
		CreateTime: time.Now().Local().Format("2006-01-02 15:04:05"),
		State:      request.State,
	}

	operationConfigs = append(operationConfigs, operationConfig)

	b, err := json.Marshal(operationConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	product.AlarmConfigs = string(b)
	return db.UpdateProduct(product)
}

//更新物模型
func UpdateProductAlarmService(request dto.UpdateProductAlarmDataRequest) error {
	product, err := db.QueryProductByID(request.ProductId)
	if err != nil {
		return err
	}
	items := product.AlarmConfigs
	var operationConfigs []dto.ProductAlarmConfig
	err = json.Unmarshal([]byte(items), &operationConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	//组装数据
	operationConfig := dto.ProductAlarmConfig{
		Key:        request.Key,
		Level:      request.Level,
		Name:       request.Name,
		Code:       request.Code,
		Type:       request.Type,
		Conditions: request.Conditions,
		Operations: request.Operations,
		Message:    request.Message,
		Desc:       request.Desc,
		State:      request.State,
	}

	pos := -1
	for index, item := range operationConfigs {
		if item.Key == request.Key {
			pos = index
		}
	}
	if pos == -1 {
		return errors.New("没有找到操作")
	}
	//替换数据
	operationConfigs[pos] = operationConfig

	b, err := json.Marshal(operationConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	product.AlarmConfigs = string(b)
	return db.UpdateProduct(product)
}

//查询
func QueryProductAlarmService(productId int) ([]dto.ProductAlarmConfigItem, error) {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return nil, err
	}
	itemConfigs, err := ConvertAlarmConfig(product.AlarmConfigs)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}

	list := make([]dto.ProductAlarmConfigItem, 0)
	for _, item := range itemConfigs {
		d := dto.ProductAlarmConfigItem{
			Key:        item.Key,
			Level:      item.Level,
			Message:    item.Message,
			Name:       item.Name,
			Code:       item.Code,
			Type:       item.Type,
			CreateTime: item.CreateTime,
			Conditions: item.Conditions,
			State:      item.State,
		}
		list = append(list, d)
	}
	return list, nil
}

//删除物模型
func DeleteProductAlarmSerivce(productId int, itemKey string) error {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return err
	}
	items := product.AlarmConfigs
	var itemConfigs []dto.ProductAlarmConfig
	err = json.Unmarshal([]byte(items), &itemConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	//删除
	pos := -1
	for index, item := range itemConfigs {
		if item.Key == itemKey {
			pos = index
			break
		}
	}
	if pos == -1 {
		return errors.New("没有找到物模型")
	}
	newAlarmConfig := append(itemConfigs[:pos], itemConfigs[pos+1:]...)

	b, err := json.Marshal(newAlarmConfig)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	product.AlarmConfigs = string(b)
	return db.UpdateProduct(product)
}

func DetailProductAlarmSerivce(productId int, itemKey string) (*dto.ProductAlarmConfig, error) {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return nil, err
	}

	items := product.AlarmConfigs
	var itemConfigs []dto.ProductAlarmConfig
	err = json.Unmarshal([]byte(items), &itemConfigs)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}

	for _, item := range itemConfigs {
		if item.Key == itemKey {
			return &item, nil
		}
	}

	return nil, errors.New("没有找到信息")
}

func ConvertAlarmConfig(items string) ([]dto.ProductAlarmConfig, error) {
	if len(items) == 0 {
		return nil, errors.New("null")
	}
	var itemConfigs []dto.ProductAlarmConfig
	err := json.Unmarshal([]byte(items), &itemConfigs)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}

	return itemConfigs, nil
}
