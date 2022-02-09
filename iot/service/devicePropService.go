package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"koudai-box/iot/db"
	"koudai-box/iot/gateway/model"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func GetLastestProperty(deviceId string) (gin.H, error) {
	deviceProp, err := db.GetLastestProperty(deviceId)

	if err != nil {
		return nil, err
	}
	properties_map := make(map[string]model.PropertyItem)
	err = json.Unmarshal([]byte(deviceProp.Properties), &properties_map)
	if err != nil {
		return nil, errors.New("属性转换失败")
	}
	item := gin.H{
		"id":         deviceProp.Id,
		"properties": &properties_map,
		"timestamp":  time.Unix(deviceProp.Timestamp, 0).Local().Format("2006-01-02 15:04:05"),
	}
	return item, err
}

func CalcPredayAvg(deviceId int) *model.PropertyMessage {
	logrus.Infof("CalcPredayAvg device[%d]", deviceId)
	deviceProps := db.GetPredayProps(deviceId, time.Now())

	if deviceProps == nil {
		return nil
	}
	num := len(deviceProps)
	var zero map[string]*model.PropertyItem
	for n, prop := range deviceProps {
		properties_map := make(map[string]*model.PropertyItem)
		err := json.Unmarshal([]byte(prop.Properties), &properties_map)
		if err != nil {
			num -= 1
			continue
		}
		if n == 0 {
			zero = properties_map
			continue
		}
		for key, item := range properties_map {
			zero[key].Value = calcUtil(zero[key].Value, item.Value, "+", zero[key].Value)
		}
	}
	if num == 0 {
		return nil
	}
	tmp := make(map[string]model.PropertyItem)
	for key, item := range zero {
		item.Value = calcUtil(item.Value, num, "/", item.Value)
		tmp[key] = *item
	}

	return &model.PropertyMessage{
		Properties: tmp,
	}
}

func calcUtil(a, b interface{}, oper string, default_value interface{}) interface{} {
	t1, err := decimal.NewFromString(fmt.Sprintf("%+v", a))
	if err != nil {
		return default_value
	}
	t2, err := decimal.NewFromString(fmt.Sprintf("%+v", b))
	if err != nil {
		return default_value
	}
	var res decimal.Decimal
	if oper == "+" {
		res = t1.Add(t2)
	} else if oper == "/" {
		res = t1.Div(t2)
		numstr := t1.String()
		tmp := strings.Split(numstr, ".")
		if len(tmp) > 1 {
			bit := len(tmp[1])
			if bit > 0 {
				res = res.Round(int32(bit))
			}
		} else {
			res = res.Round(2)
		}
	}
	switch a.(type) {
	case int8, int16, int32, int, int64:
		t, _ := strconv.ParseInt(res.String(), 10, 64)
		return t
	case float32, float64:
		t, _ := strconv.ParseFloat(res.String(), 64)
		return t
	default:
		return default_value
	}
}
