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
func AddProductItemService(request dto.AddProductItemDataRequest) error {
	product, err := db.QueryProductByID(request.ProductId)
	if err != nil {
		return err
	}
	items := product.Items
	var itemConfigs []dto.ProductItemConfig
	err = json.Unmarshal([]byte(items), &itemConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	//检测code是否项目
	for _, item := range itemConfigs {
		if item.Code == request.Code {
			return errors.New("属性标识不能相同")
		}
	}

	var boolValue map[string]string
	if request.BoolValue != "" {
		err = json.Unmarshal([]byte(request.BoolValue), &boolValue)
		if err != nil {
			logger.Errorln(err)
			return err
		}
	}

	var dict map[string]string
	if request.Dict != "" {
		err = json.Unmarshal([]byte(request.Dict), &dict)
		if err != nil {
			logger.Errorln(err)
			return err
		}
	}

	//组装数据
	itemDataType := model.ItemDataType{
		RW:             request.RW,
		Type:           request.DataType,
		Unit:           request.Unit,
		Min:            request.Min,
		Max:            request.Max,
		Step:           request.Step,
		Precision:      request.Precision,
		MaxLength:      request.MaxLength,
		BoolValue:      boolValue,
		DateFormat:     request.DateFormat,
		Dict:           dict,
		FileType:       request.FileType,
		PasswordLength: request.PasswordLength,
	}
	itemConfig := dto.ProductItemConfig{
		Key:           utils.GetUUID(),
		Name:          request.Name,
		Code:          request.Code,
		SourceCode:    request.SourceCode,
		NodeId:        request.NodeId,
		Address:       request.Address,
		Quantity:      request.Quantity,
		OperaterType:  request.OperaterType,
		ZoomFactor:    request.ZoomFactor,
		ExchangeHL:    request.ExchangeHL,
		ExchangeOrder: request.ExchangeOrder,
		DataType:      itemDataType,
		Report:        request.Report,
		Desc:          request.Desc,
	}
	itemConfigs = append(itemConfigs, itemConfig)

	b, err := json.Marshal(itemConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	product.Items = string(b)
	return db.UpdateProduct(product)
}

//更新物模型
func UpdateProductItemService(request dto.UpdateProductItemDataRequest) error {
	product, err := db.QueryProductByID(request.ProductId)
	if err != nil {
		return err
	}
	items := product.Items
	var itemConfigs []dto.ProductItemConfig
	err = json.Unmarshal([]byte(items), &itemConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	var boolValue map[string]string
	if request.BoolValue != "" {
		err = json.Unmarshal([]byte(request.BoolValue), &boolValue)
		if err != nil {
			logger.Errorln(err)
			return err
		}
	}

	var dict map[string]string
	if request.Dict != "" {
		err = json.Unmarshal([]byte(request.Dict), &dict)
		if err != nil {
			logger.Errorln(err)
			return err
		}
	}

	//组装数据
	itemDataType := model.ItemDataType{
		RW:             request.RW,
		Type:           request.DataType,
		Unit:           request.Unit,
		Min:            request.Min,
		Max:            request.Max,
		Step:           request.Step,
		Precision:      request.Precision,
		MaxLength:      request.MaxLength,
		BoolValue:      boolValue,
		DateFormat:     request.DateFormat,
		Dict:           dict,
		FileType:       request.FileType,
		PasswordLength: request.PasswordLength,
	}
	itemConfig := dto.ProductItemConfig{
		Key:           request.Key,
		Name:          request.Name,
		Code:          request.Code,
		SourceCode:    request.SourceCode,
		NodeId:        request.NodeId,
		Address:       request.Address,
		Quantity:      request.Quantity,
		OperaterType:  request.OperaterType,
		ZoomFactor:    request.ZoomFactor,
		ExchangeHL:    request.ExchangeHL,
		ExchangeOrder: request.ExchangeOrder,
		DataType:      itemDataType,
		Report:        request.Report,
		Desc:          request.Desc,
	}

	pos := -1
	for index, item := range itemConfigs {
		if item.Key == request.Key {
			pos = index
		}
	}
	if pos == -1 {
		return errors.New("没有找到物模型")
	}
	//替换数据
	itemConfigs[pos] = itemConfig

	b, err := json.Marshal(itemConfigs)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	product.Items = string(b)
	return db.UpdateProduct(product)
}

//查询
func QueryProductItemService(productId int) ([]dto.ProductItemConfigItem, error) {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return nil, err
	}
	itemConfigs, err := ConvertItemConfig(product.Items)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}

	list := make([]dto.ProductItemConfigItem, 0)
	for _, item := range itemConfigs {
		d := dto.ProductItemConfigItem{
			Key:            item.Key,
			Code:           item.Code,
			SourceCode:     item.SourceCode,
			NodeId:         item.NodeId,
			Address:        item.Address,
			Quantity:       item.Quantity,
			OperaterType:   item.OperaterType,
			ZoomFactor:     item.ZoomFactor,
			ExchangeHL:     item.ExchangeHL,
			ExchangeOrder:  item.ExchangeOrder,
			Name:           item.Name,
			Min:            item.DataType.Min,
			Max:            item.DataType.Max,
			Step:           item.DataType.Step,
			RW:             item.DataType.RW,
			DataType:       item.DataType.Type,
			Unit:           item.DataType.Unit,
			Precision:      item.DataType.Precision,
			MaxLength:      item.DataType.MaxLength,
			BoolValue:      item.DataType.BoolValue,
			DateFormat:     item.DataType.DateFormat,
			Dict:           item.DataType.Dict,
			FileType:       item.DataType.FileType,
			PasswordLength: item.DataType.PasswordLength,
			Desc:           item.Desc,
			Report:         item.Report,
		}
		list = append(list, d)
	}
	return list, nil
}

//删除物模型
func DeleteProductItemSerivce(productId int, itemKey string) error {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return err
	}
	items := product.Items
	var itemConfigs []dto.ProductItemConfig
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
	newItemConfig := append(itemConfigs[:pos], itemConfigs[pos+1:]...)

	b, err := json.Marshal(newItemConfig)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	product.Items = string(b)
	return db.UpdateProduct(product)
}

func DetailProductItemSerivce(productId int, itemKey string) (*dto.ProductItemConfig, error) {
	product, err := db.QueryProductByID(productId)
	if err != nil {
		return nil, err
	}

	items := product.Items
	var itemConfigs []dto.ProductItemConfig
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

func ConvertItemConfig(items string) ([]dto.ProductItemConfig, error) {
	if len(items) == 0 {
		return nil, errors.New("null")
	}
	var itemConfigs []dto.ProductItemConfig
	err := json.Unmarshal([]byte(items), &itemConfigs)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}

	return itemConfigs, nil
}
