package test

import (
	"encoding/json"
	"fmt"
	. "main/model"
	"testing"
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

func TestParseJson(t *testing.T) {
	data := "{\"a\":1,\"b\":2,\"c\":3}"
	model := "{\"items\":[{\"key\":\"da\",\"source\":\"a\"},{\"key\":\"db\",\"source\":\"b\"}]}"
	str, err := ParseJson(data, model)
	if err != nil {
		fmt.Println("ParseJson error:", err)
	}
	fmt.Println("ParseJson result:", str)
}
