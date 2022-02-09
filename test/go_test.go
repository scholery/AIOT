package test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"koudai-box/iot/gateway/utils"

	"github.com/sirupsen/logrus"
)

func TestGo(t *testing.T) {
	start := time.Now() // 获取当前时间
	i := "2"
	v2 := fmt.Sprintf("%v", i)
	t2, _ := strconv.Atoi(v2)
	fmt.Printf("v2=%s,t2=%d\r\n", v2, t2)

	str := "我是${name},你是${you}"
	d := make(map[string]interface{})
	d["name"] = "zhucz"
	d["you"] = "kevin"
	d["me"] = "eeee"
	s := utils.ParseTpl(str, d)
	fmt.Printf("tpl=%s\r\n", s)
	elapsed := time.Since(start)
	logrus.Infof("tpl=%s", s)
	logrus.Infof("该函数执行完成耗时：%v", elapsed)
}
