package main

import (
	"net/http"
	"os/signal"
	"runtime"
	"syscall"

	"koudai-box/cache"
	"koudai-box/conf"
	"koudai-box/global"

	"koudai-box/iot/db"
	iot "koudai-box/iot/gateway/service"
	"koudai-box/iot/service"
	"koudai-box/iot/web"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)

	var server *http.Server
	signal.Notify(global.SystemExitChannel, syscall.SIGINT)
	signal.Notify(global.SystemExitChannel, syscall.SIGTERM)

	go func() {
		sig := <-global.SystemExitChannel
		logrus.Infof("收到停止指令 %s signal, threads:%d", sig.String(), runtime.NumGoroutine())
		if server != nil {
			server.Close()
		}
		destroy()
		logrus.Infof("服务停止，threads:%d", runtime.NumGoroutine())
	}()
	conf.InitConf("./conf/config.json")
	cache.Init()

	if err := initSystem(); err != nil {
		logrus.Error(err)
		return
	}
	iot.InitIot()
	server = web.Init(conf.GetConf().WebPort)
	server.ListenAndServe()

}

func initSystem() error {
	logrus.Infof("服务启动 init")
	db.DBInit()
	service.InitDictCache()
	service.InitListener()
	service.StartCron()
	return nil
}

func destroy() {
	logrus.Infof("收到停止指令 destroy")
	iot.StopPull()
	service.StopCron()
}
