package utils

import (
	"main/model"
	"math"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func GetUUID() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
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
	case model.Int32:
		v1, ok1 := val1.(int)
		v2, ok2 := val2.(int)
		s, err := strconv.Atoi(dataType.Step)
		if ok1 && ok2 && err == nil {
			noChange = Abs(v1-v2) <= s
		}
	case model.Float:
		v1, ok1 := val1.(float64)
		v2, ok2 := val2.(float64)
		s, err := strconv.ParseFloat(dataType.Step, 32)
		if ok1 && ok2 && err == nil {
			noChange = math.Abs(float64(v1-v2)) <= s
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

func Compare(val1 interface{}, val2 interface{}, dataType string) int {
	res := 0
	switch dataType {
	case "int":
		t1, _ := val1.(int)
		t2, _ := val2.(int)
		if t1 > t2 {
			res = 1
		} else if t1 < t2 {
			res = -1
		}
	case "float32":
		t1, _ := val1.(float32)
		t2, _ := val2.(float32)
		if t1 > t2 {
			res = 1
		} else if t1 < t2 {
			res = -1
		}
	case "float64":
		t1, _ := val1.(float64)
		t2, _ := val2.(float64)
		if t1 > t2 {
			res = 1
		} else if t1 < t2 {
			res = -1
		}
	case "string":
		t1, _ := val1.(string)
		t2, _ := val2.(string)
		res = strings.Compare(t1, t2)
	}
	logrus.Info("Compare ", dataType, val1, val2, res)
	return res
}
