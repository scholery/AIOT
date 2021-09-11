package utils

import (
	"main/model"
	"math"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
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

func PropCompareEQ(val1 interface{}, val2 interface{}, dataType string, step string) bool {
	noChange := false
	switch dataType {
	case "int":
		v1, ok1 := val1.(int)
		v2, ok2 := val2.(int)
		s, err := strconv.Atoi(step)
		if ok1 && ok2 && err == nil {
			noChange = Abs(v1-v2) <= s
		}
	case "float32":
		v1, ok1 := val1.(float32)
		v2, ok2 := val2.(float32)
		s, err := strconv.ParseFloat(step, 32)
		if ok1 && ok2 && err == nil {
			noChange = math.Abs(float64(v1-v2)) <= s
		}
	case "float64":
		v1, ok1 := val1.(float64)
		v2, ok2 := val2.(float64)
		s, err := strconv.ParseFloat(step, 64)
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
	var v1, v2 float64
	switch dataType {
	case "int":
		t1, _ := val1.(int)
		t2, _ := val2.(int)
		v1 = float64(t1)
		v2 = float64(t2)
	case "float32":
		t1, _ := val1.(float32)
		t2, _ := val2.(float32)
		v1 = float64(t1)
		v2 = float64(t2)
	case "float64":
		t1, _ := val1.(float64)
		t2, _ := val2.(float64)
		v1 = float64(t1)
		v2 = float64(t2)
	}
	if v1 > v2 {
		res = 1
	} else if v1 < v2 {
		res = -1
	}
	return res
}
