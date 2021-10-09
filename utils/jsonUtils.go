package utils

import (
	"main/model"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetPropertyItem(item model.ItemConfig, value interface{}) model.PropertyItem {

	return model.PropertyItem{Key: item.Key, Name: item.Name, Value: value, DataType: item.DataType}
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
