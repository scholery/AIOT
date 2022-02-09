package test

import (
	"encoding/json"
	"fmt"
	"time"

	"koudai-box/iot/gateway/model"
	. "koudai-box/iot/gateway/model"

	"github.com/robertkrimen/otto"
)

type user struct {
	name     string
	salt     string
	password string
}

func Exec() {
	funstr := "var user = {'name':'zhucz'};function password_enc(model){console.log('model.key=',model.items[0].key);user.name = 'xxxxx';return model.items[0];}"
	err := password_enc(funstr)
	fmt.Println("err=", err)
}

func password_enc(js string) error {
	start := time.Now() // 获取当前时间
	vm := otto.New()
	_, err := vm.Run(js)
	modelStr := "{\"items\":[{\"key\":\"da\",\"source\":\"a\"},{\"key\":\"db\",\"source\":\"b\"}]}"
	var config Product
	err = json.Unmarshal([]byte(modelStr), &config)

	aaa, _ := vm.Get("user")
	bbb, _ := aaa.Object().Get("name")
	fmt.Println("enc1", bbb)
	enc, err := vm.Call("password_enc", nil, config)
	item1, _ := enc.Export()
	it := item1.(model.ItemConfig)
	fmt.Println("item=", item1)
	fmt.Println("it=", it.Key)
	ccc, _ := vm.Get("user")
	ddd, _ := ccc.Object().Get("name")
	fmt.Println("enc2", ddd)
	r := enc.Object()
	key, _ := enc.Object().Get("key")
	fmt.Println("enc.key=", key)
	fmt.Println("r=", r)
	elapsed := time.Since(start)
	fmt.Println("该函数执行完成耗时：", elapsed)
	return err
}
