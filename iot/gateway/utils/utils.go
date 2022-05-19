package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"koudai-box/module/svsdk"

	"koudai-box/iot/gateway/model"

	uuid "github.com/satori/go.uuid"
)

func GetUUID() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

func GetSN() string {
	return svsdk.GetSN()
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func PropCompareEQ(val1 interface{}, val2 interface{}, dataType model.ItemDataType) bool {
	noChange := false
	switch dataType.Type {
	case model.Int32, model.Int64, model.Long:
		v1 := fmt.Sprintf("%d", val1)
		v2 := fmt.Sprintf("%d", val2)
		t1, ok1 := strconv.Atoi(v1)
		t2, ok2 := strconv.Atoi(v2)
		s, err := strconv.Atoi(dataType.Step)
		if ok1 == nil && ok2 == nil && err == nil {
			noChange = Abs(t1-t2) <= s
		}
	case model.Float, model.Double:
		v1 := fmt.Sprintf("%f", val1)
		v2 := fmt.Sprintf("%f", val2)
		t1, ok1 := strconv.ParseFloat(v1, 64)
		t2, ok2 := strconv.ParseFloat(v2, 64)
		s, err := strconv.ParseFloat(dataType.Step, 64)
		if ok1 == nil && ok2 == nil && err == nil {
			noChange = math.Abs(float64(t1-t2)) <= s
		}
	}
	return !noChange
}

func MatchContidion(prop model.PropertyMessage, condition model.Condition) bool {
	match := false
	v, ok := prop.Properties[condition.Key]
	if !ok {
		return match
	}
	res := Compare(v.Value, condition.Value, condition.DataType)
	switch condition.Compare {
	case "=":
		match = res == 0
	case ">":
		match = res == 1
	case ">=":
		match = res >= 0
	case "<":
		match = res == -1
	case "<=":
		match = res <= 0
	}
	return match
}

func Compare(val1 interface{}, val2 interface{}, dataType model.DataType) int {
	res := 0
	switch dataType {
	case model.Int32, model.Int64, model.Long:
		v1 := fmt.Sprintf("%+v", val1)
		v2 := fmt.Sprintf("%+v", val2)
		t1, _ := strconv.Atoi(v1)
		t2, _ := strconv.Atoi(v2)
		if t1 > t2 {
			res = 1
		} else if t1 < t2 {
			res = -1
		}
	case model.Float, model.Double:
		v1 := fmt.Sprintf("%+v", val1)
		v2 := fmt.Sprintf("%+v", val2)
		t1, _ := strconv.ParseFloat(v1, 64)
		t2, _ := strconv.ParseFloat(v2, 64)
		if t1 > t2 {
			res = 1
		} else if t1 < t2 {
			res = -1
		}
	case model.Text:
		t1 := fmt.Sprintf("%+v", val1)
		t2 := fmt.Sprintf("%+v", val2)
		res = strings.Compare(t1, t2)
	default:
		t1 := fmt.Sprintf("%+v", val1)
		t2 := fmt.Sprintf("%+v", val2)
		res = strings.Compare(t1, t2)
	}
	return res
}

func getTplExpressions(str string) []string {
	reg_str := `\$\{.*?\}`
	re, _ := regexp.Compile(reg_str)
	all := re.FindAll([]byte(str), 2)
	keyArrays := make([]string, 0)
	for _, item := range all {
		item_str := string(item)
		if len(item_str) > 3 {
			item_str = item_str[2 : len(item_str)-1]
			keyArrays = append(keyArrays, item_str)
		}

	}
	return keyArrays
}

// 将tpl中的占位符 替换为真实值 ${data.0.att1}
func ParseTpl(tpl string, data map[string]interface{}) string {
	if len(tpl) < 4 {
		return tpl
	}
	expressions := getTplExpressions(tpl)
	for _, exp := range expressions {
		//fmt.Println("exp",exp)
		exp = strings.TrimSpace(exp)
		v, ok := data[exp]
		if !ok {
			continue
		}
		val := fmt.Sprintf("%+v", v)
		tpl = strings.Replace(tpl, "${"+exp+"}", val, -1)
	}
	return tpl
}
