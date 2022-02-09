package test

import (
	"encoding/json"

	. "koudai-box/iot/gateway/model"
)

func ParseJson(data string, modelStr string) (string, error) {
	var config Product
	err := json.Unmarshal([]byte(modelStr), &config)

	if err != nil {
		return "", err
	}
	dataMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(data), &dataMap)
	if err != nil {
		return "", err
	}
	dataTmp := make(map[string]interface{})
	for _, item := range config.Items {
		dataTmp[item.Key] = dataMap[item.Source]
	}
	str, er := json.Marshal(dataTmp)
	return string(str), er
}
