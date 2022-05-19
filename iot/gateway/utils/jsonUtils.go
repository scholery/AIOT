package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"koudai-box/iot/gateway/model"

	"github.com/sirupsen/logrus"
)

//数据转换接口
func Transformer2DeviceProp(data interface{}, device *model.Device) (interface{}, error) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("transformer:data format error")
	}
	if len(device.Product.Items) == 0 {
		return nil, errors.New("product model item is empty")
	}
	dataTmp := make(map[string]model.PropertyItem)
	for _, item := range device.Product.Items {
		dataTmp[item.Code] = GetPropertyItem(item, GetMapValue(dataMap, item.Source))
	}

	return model.PropertyMessage{DeviceId: strconv.Itoa(device.Id), DeviceSign: device.Key, MessageId: GetUUID(),
		Timestamp: time.Now().Unix(), Properties: dataTmp}, nil
}

func GetPropertyItem(item model.ItemConfig, value interface{}) model.PropertyItem {
	var val float64
	if nil != value && item.ZoomFactor != 0 {
		t1, err := strconv.ParseFloat(fmt.Sprintf("%+v", value), 64)
		if err == nil {
			val = t1 * float64(item.ZoomFactor)
		}
	}
	if val != 0 {
		bit := 0
		numstr := fmt.Sprint(val)
		tmp := strings.Split(numstr, ".")
		if len(tmp) > 1 {
			bit = len(tmp[1])
		}
		if bit > 0 && bit > item.DataType.Precision {
			v1 := strconv.FormatFloat(val, 'f', item.DataType.Precision, 64)
			t1, err := strconv.ParseFloat(v1, 64)
			if err == nil {
				val = t1
			}
		}
		if len(item.DataType.Min) > 0 {
			min, err := strconv.ParseFloat(item.DataType.Min, 64)
			if err == nil && min > val {
				val = min
			}
		}
		if len(item.DataType.Max) > 0 {
			max, err := strconv.ParseFloat(item.DataType.Max, 64)
			if err == nil && max < val {
				val = max
			}
		}
	}
	if val != 0 {
		return model.PropertyItem{Sort: item.Sort, Code: item.Code, Name: item.Name, Value: val, DataType: item.DataType}
	}
	return model.PropertyItem{Sort: item.Sort, Code: item.Code, Name: item.Name, Value: value, DataType: item.DataType}
}
func GetMapValue(dataMap map[string]interface{}, key string) interface{} {
	keys := strings.Split(key, ".")
	if len(keys) == 0 {
		return ""
	}
	size := len(keys)
	no, keyTmp := IsArray(keys[0])
	//单个
	if no < 0 && size == 1 {
		return dataMap[keyTmp]
	} else if no < 0 && size > 1 {
		data, ok := dataMap[keyTmp]
		if !ok {
			logrus.Errorf("parse item error,not found, key=%s ", keyTmp)
			return ""
		}
		tmp, ok := data.(map[string]interface{})
		if !ok {
			logrus.Errorf("parse item error,type is not obj, key=%s ", keyTmp)
			return data
		}
		return GetMapValue(tmp, strings.Join(keys[1:], "."))
	}
	//数组
	data, ok := dataMap[keyTmp]
	if !ok {
		logrus.Errorf("parse array item error,not found, key=%s ", keyTmp)
		return data
	}
	tmp, ok := data.([]interface{})
	if !ok {
		logrus.Errorf("parse array item error,type is not array, key=%s ", keyTmp)
		return ""
	}
	if size == 1 {
		//解析基本类型数组
		if len(tmp) > no {
			return tmp[no]
		} else {
			logrus.Errorf("parse array item error,index out of range, key=%s ", key)
			return ""
		}
	} else {
		//解析对象数组
		if len(tmp) <= no {
			logrus.Errorf("parse array item error,index out of range, key=%s", key)
			return ""
		}
		tmp1, ok := tmp[no].(map[string]interface{})
		if !ok {
			logrus.Errorf("parse array item error,type is not obj, key=%s ", keyTmp)
			return ""
		}
		return GetMapValue(tmp1, strings.Join(keys[1:], "."))
	}
}

func IsArray(key string) (int, string) {
	keyTmp := key
	if !strings.Contains(key, "[") || !strings.Contains(key, "]") {
		return -1, keyTmp
	}
	index := strings.Index(key, "[")
	index1 := strings.Index(key, "]")
	no := -1
	if index1 > index && index1 < len(key) {
		t := key[index+1 : index1]
		t2, err := strconv.Atoi(t)
		if err == nil {
			no = t2
			keyTmp = key[0:index]
		}
	}
	return no, keyTmp
}

func ToString(data interface{}) string {
	if data == nil {
		return ""
	}
	byte, err := json.Marshal(data)
	if err != nil {
		logrus.Error("to json error:", data)
		return ""
	}
	ret := string(byte)
	if ret == "null" {
		return ""
	}
	return ret
}

func ToMap(data interface{}) map[string]interface{} {
	str, ok := data.(string)
	var by []byte
	if !ok {
		bb, err := json.Marshal(data)
		if err != nil {
			logrus.Error("to json error:", data)
			return nil
		}
		by = bb
	} else {
		by = []byte(str)
	}

	m := make(map[string]interface{})
	json.Unmarshal(by, &m)
	return m
}
