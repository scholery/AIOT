package service

import (
	"encoding/json"

	"koudai-box/iot/db"
	"koudai-box/iot/gateway/utils"
	"koudai-box/iot/web/dto"
)

//添加物模型
func AddProductFunctionService(request dto.AddProductFunctionDataRequest) error {
	product, err := db.QueryProductByID(request.ProductId)
	if err != nil {
		return err
	}

	//组装数据
	functionConfig := dto.ProductFunctionConfig{
		Key:          utils.GetUUID(),
		ExtractProp:  request.ExtractProp,
		ExtractEvent: request.ExtractEvent,
		Calculate:    request.Calculate,
		Body:         request.Body,
	}

	b, err := json.Marshal(functionConfig)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	product.FunctionConfigs = string(b)
	return db.UpdateProduct(product)
}

//更新物模型
func UpdateProductFunctionService(request dto.UpdateProductFunctionDataRequest) error {
	product, err := db.QueryProductByID(request.ProductId)
	if err != nil {
		return err
	}
	//组装数据
	functionConfigs := dto.ProductFunctionConfig{
		Key:          request.Key,
		ExtractProp:  request.ExtractProp,
		ExtractEvent: request.ExtractEvent,
		Calculate:    request.Calculate,
		Body:         request.Body,
	}

	b, err := json.Marshal(functionConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	product.FunctionConfigs = string(b)
	return db.UpdateProduct(product)
}

//查询
func QueryProductFunctionService(productId int) (dto.ProductFunctionConfigItem, error) {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return dto.ProductFunctionConfigItem{}, err
	}
	items := product.FunctionConfigs
	if len(items) == 0 {
		return dto.ProductFunctionConfigItem{}, nil
	}
	var item dto.ProductFunctionConfigItem
	err = json.Unmarshal([]byte(items), &item)
	if err != nil {
		logger.Errorln(err)
		return dto.ProductFunctionConfigItem{}, err
	}

	return item, nil
}

func DetailProductFunctionSerivce(productId int, itemKey string) (*dto.ProductFunctionConfig, error) {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return nil, err
	}

	var itemConfigs dto.ProductFunctionConfig
	err = json.Unmarshal([]byte(product.FunctionConfigs), &itemConfigs)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}

	return &itemConfigs, nil
}
