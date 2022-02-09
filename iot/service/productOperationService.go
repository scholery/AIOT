package service

import (
	"encoding/json"
	"errors"

	"koudai-box/iot/db"
	"koudai-box/iot/gateway/model"
	"koudai-box/iot/gateway/utils"
	"koudai-box/iot/web/dto"
)

//添加物模型
func AddProductOperationService(request dto.AddProductOperationDataRequest) error {
	product, err := db.QueryProductByID(request.ProductId)
	if err != nil {
		return err
	}
	configs := product.OperationConfigs
	var operationConfigs []model.OperationConfig
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
	operationConfig := model.OperationConfig{
		Key:     utils.GetUUID(),
		Code:    request.Code,
		Name:    request.Name,
		Router:  request.Router,
		Desc:    request.Desc,
		Type:    request.Type,
		Inputs:  request.Inputs,
		Outputs: request.Outputs,
	}
	operationConfigs = append(operationConfigs, operationConfig)

	b, err := json.Marshal(operationConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	product.OperationConfigs = string(b)
	return db.UpdateProduct(product)
}

//更新物模型
func UpdateProductOperationService(request dto.UpdateProductOperationDataRequest) error {
	product, err := db.QueryProductByID(request.ProductId)
	if err != nil {
		return err
	}
	items := product.OperationConfigs
	var operationConfigs []model.OperationConfig
	err = json.Unmarshal([]byte(items), &operationConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	//组装数据
	operationConfig := model.OperationConfig{
		Key:     request.Key,
		Code:    request.Code,
		Name:    request.Name,
		Desc:    request.Desc,
		Type:    request.Type,
		Router:  request.Router,
		Inputs:  request.Inputs,
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
	product.OperationConfigs = string(b)
	return db.UpdateProduct(product)
}

//查询
func QueryProductOperationService(productId int) ([]dto.ProductOperationConfigItem, error) {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return nil, err
	}
	items := product.OperationConfigs
	var itemConfigs []model.OperationConfig
	err = json.Unmarshal([]byte(items), &itemConfigs)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}

	list := make([]dto.ProductOperationConfigItem, 0)
	for _, item := range itemConfigs {
		d := dto.ProductOperationConfigItem{
			Key:     item.Key,
			Code:    item.Code,
			Name:    item.Name,
			Router:  item.Router,
			Type:    item.Type,
			Desc:    item.Desc,
			Inputs:  item.Inputs,
			Outputs: item.Outputs,
		}
		list = append(list, d)
	}
	return list, nil
}

//删除物模型
func DeleteProductOperationSerivce(productId int, itemKey string) error {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return err
	}
	items := product.OperationConfigs
	var itemConfigs []model.OperationConfig
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
	newOperationConfig := append(itemConfigs[:pos], itemConfigs[pos+1:]...)

	b, err := json.Marshal(newOperationConfig)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	product.OperationConfigs = string(b)
	return db.UpdateProduct(product)
}

func DetailProductOperationSerivce(productId int, itemKey string) (*model.OperationConfig, error) {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return nil, err
	}

	items := product.OperationConfigs
	var itemConfigs []model.OperationConfig
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
