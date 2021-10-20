package utils

import (
	"fmt"
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

func PropCompareEQ(val1 interface{}, val2 interface{}, dataType model.ItemDataType) bool {
	noChange := false
	switch dataType.Type {
	case model.Int32:
		v1 := fmt.Sprintf("%d", val1)
		v2 := fmt.Sprintf("%d", val2)
		t1, ok1 := strconv.Atoi(v1)
		t2, ok2 := strconv.Atoi(v2)
		s, err := strconv.Atoi(dataType.Step)
		if ok1 == nil && ok2 == nil && err == nil {
			noChange = Abs(t1-t2) <= s
		}
	case model.Float:
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
	case model.Int32:
		v1 := fmt.Sprintf("%d", val1)
		v2 := fmt.Sprintf("%d", val2)
		t1, _ := strconv.Atoi(v1)
		t2, _ := strconv.Atoi(v2)
		if t1 > t2 {
			res = 1
		} else if t1 < t2 {
			res = -1
		}
	case model.Float:
		v1 := fmt.Sprintf("%f", val1)
		v2 := fmt.Sprintf("%f", val2)
		t1, _ := strconv.ParseFloat(v1, 64)
		t2, _ := strconv.ParseFloat(v2, 64)
		if t1 > t2 {
			res = 1
		} else if t1 < t2 {
			res = -1
		}
	case model.Text:
		t1 := fmt.Sprintf("%v", val1)
		t2 := fmt.Sprintf("%v", val2)
		res = strings.Compare(t1, t2)
	}
	return res
}
