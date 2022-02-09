package service

import (
	"encoding/json"
	"errors"

	"koudai-box/iot/db"
	"koudai-box/iot/gateway/utils"
	"koudai-box/iot/web/dto"
)

//添加物模型
func AddProductEventService(request dto.AddProductEventDataRequest) error {
	product, err := db.QueryProductByID(request.ProductId)
	if err != nil {
		return err
	}
	configs := product.EventConfigs
	var operationConfigs []dto.ProductEventConfig
	err = json.Unmarshal([]byte(configs), &operationConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	//检测code是否项目
	for _, item := range operationConfigs {
		if item.Code == request.Code {
			return errors.New("属性标识不能相同")
		}
	}

	//组装数据
	operationConfig := dto.ProductEventConfig{
		Key:     utils.GetUUID(),
		Code:    request.Code,
		Name:    request.Name,
		Desc:    request.Desc,
		Type:    request.Type,
		Outputs: request.Outputs,
	}
	operationConfigs = append(operationConfigs, operationConfig)

	b, err := json.Marshal(operationConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	product.EventConfigs = string(b)
	return db.UpdateProduct(product)
}

//更新物模型
func UpdateProductEventService(request dto.UpdateProductEventDataRequest) error {
	product, err := db.QueryProductByID(request.ProductId)
	if err != nil {
		return err
	}
	items := product.EventConfigs
	var operationConfigs []dto.ProductEventConfig
	err = json.Unmarshal([]byte(items), &operationConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	//组装数据
	operationConfig := dto.ProductEventConfig{
		Key:     request.Key,
		Code:    request.Code,
		Name:    request.Name,
		Desc:    request.Desc,
		Type:    request.Type,
		Outputs: request.Outputs,
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
	product.EventConfigs = string(b)
	return db.UpdateProduct(product)
}

//查询
func QueryProductEventService(productId int) ([]dto.ProductEventConfigItem, error) {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return nil, err
	}
	items := product.EventConfigs
	var itemConfigs []dto.ProductEventConfig
	err = json.Unmarshal([]byte(items), &itemConfigs)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}

	list := make([]dto.ProductEventConfigItem, 0)
	for _, item := range itemConfigs {
		d := dto.ProductEventConfigItem{
			Key:     item.Key,
			Code:    item.Code,
			Name:    item.Name,
			Type:    item.Type,
			Desc:    item.Desc,
			Outputs: item.Outputs,
		}
		list = append(list, d)
	}
	return list, nil
}

//删除物模型
func DeleteProductEventSerivce(productId int, itemKey string) error {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return err
	}
	items := product.EventConfigs
	var itemConfigs []dto.ProductEventConfig
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
	newEventConfig := append(itemConfigs[:pos], itemConfigs[pos+1:]...)

	b, err := json.Marshal(newEventConfig)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	product.EventConfigs = string(b)
	return db.UpdateProduct(product)
}

func DetailProductEventSerivce(productId int, itemKey string) (*dto.ProductEventConfig, error) {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return nil, err
	}

	items := product.EventConfigs
	var itemConfigs []dto.ProductEventConfig
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
