package main

import (
	//. "main/model"

	"main/service"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	service.Connect()
	defer service.Close()
	go service.StartPull()
	i := 1
	for {
		time.Sleep(time.Second)
		if i == 100 {
			service.StopPull()
			break
		}
		if i == 10 {
			logrus.Infof("runtime threads:%d", runtime.NumGoroutine())
		}
		i++
	}
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	//test.Exec()
	//TestTCPClientAdvancedUsage()
	//ExecOpc()
	//go ExecRun("AAAA")
	//go ExecRun("BBBB")
	//for true {
	//}
	//TestParseJson()
	// item := ItemConfig{Key: "1", RW: "r", Name: "Name"}
	// fmt.Println(item.Name)
	// c := make(chan PropertyMessage)
	// go test.ExecHttpTest(c)
	// // time.Sleep(time.Second)
	// i := 1
	// for data := range c {
	// 	fmt.Println("data=", data.MessageId)
	// 	fmt.Println("index=", i)
	// 	i++
	// 	if i > 10 {
	// 		break
	// 	}
	// }
	// for {

	// }
	// keys := strings.Split("a[1].b.c", ".")
	// keyTmp := keys[0]
	// index := strings.Index(keyTmp, "[")
	// index1 := strings.Index(keyTmp, "]")
	// fmt.Println(index, index1, len(keyTmp))
	// fmt.Println(keyTmp[index+1 : index1])

	// fmt.Println(len(keys))
	// fmt.Println(keys[0])
	// fmt.Println(strings.Join(keys[1:], "."))
	//g := DataGateway{Device: Device{}}
	//fmt.Println(g)
}
