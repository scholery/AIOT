package test

import (
	"runtime"
	"testing"
	"time"

	"koudai-box/iot/gateway/service"
	iot "koudai-box/iot/gateway/service"

	"github.com/sirupsen/logrus"
)

func TestRun(t *testing.T) {
	logrus.SetLevel(logrus.InfoLevel)
	service.Connect()
	defer service.Close()
	go iot.StartPull()
	i := 1
	for {
		time.Sleep(time.Second)
		if i == 100 {
			iot.StopPull()
			break
		}
		if i == 10 {
			logrus.Infof("runtime threads:%d", runtime.NumGoroutine())
		}
		i++
	}
}

func run() {
	st := make(chan bool)
	data := make(chan int, 5)
	run := true
	logrus.Info("run")
	go Set(st, data)
	for run {
		select {
		case stop := <-st:
			run = stop
			logrus.Info("stop:", stop)
		case tmp := <-data:
			logrus.Info("data:", tmp)
		}
	}
	logrus.Info("run over")
}
func Set(st chan bool, data chan int) {
	i := 1
	for {
		if i%10 == 0 {
			logrus.Info("set:", i)
			//data <- i
		}
		if i > 100 {
			st <- false
			logrus.Info("st false:", i)
			return
		}
		i++
	}
}
