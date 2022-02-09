package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"koudai-box/iot/gateway/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	// "github.com/gorilla/websocket"
)

var cache_ws_Servers = make(map[string]*WsServer)

type WsServer struct {
	IsRunning bool
	Gateway   *model.GatewayConfig
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartWsServerConnect(gateway *model.GatewayConfig) {
	r := gin.Default()
	r.GET("/event", ListenEvent)
	r.GET("/device", ListenDevice)
	r.Run(":" + strconv.Itoa(gateway.Port))
	cache_ws_Servers[strconv.Itoa(gateway.Port)] = &WsServer{IsRunning: true, Gateway: gateway}
	logrus.Infof("++++++++++++++开启Websocket服务端[%s]应用[%s:%d]++++++++++++++", gateway.Key, gateway.Ip, gateway.Port)
}

func StopWsServerConnect(gateway *model.GatewayConfig) {
	var ws = cache_ws_Servers[strconv.Itoa(gateway.Port)]
	if ws != nil {
		ws.IsRunning = false
		delete(cache_ws_Servers, strconv.Itoa(gateway.Port))
	}
	logrus.Infof("++++++++++++++关闭Websocket服务端[%s]应用[%s:%d]++++++++++++++", gateway.Key, gateway.Ip, gateway.Port)
}

func ListenDevice(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	var wsServer = cache_ws_Servers[ws.LocalAddr().String()]
	defer ws.Close() //返回前关闭
	for true {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}
		var dataMsg interface{}
		err = json.Unmarshal([]byte(message), &dataMsg)
		if err != nil {
			logrus.Errorf("ws:gateway[%s]'s msg error,msg:%s", wsServer.Gateway.Key, message)
			continue
		}
		var pushMsg model.PushMsg
		pushMsg.Msg = dataMsg
		pushMsg.GatewayKey = wsServer.Gateway.Key
		model.PushMsgChan <- pushMsg
	}
}

func ListenEvent(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	var wsServer = cache_ws_Servers[ws.LocalAddr().String()]
	defer ws.Close() //返回前关闭
	for true {
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}
		var dataMsg interface{}
		err = json.Unmarshal([]byte(message), &dataMsg)
		if err != nil {
			logrus.Errorf("ws:gateway[%s]'s msg error,msg:%s", wsServer.Gateway.Key, message)
			continue
		}
		var pushMsg model.PushMsg
		pushMsg.Msg = dataMsg
		pushMsg.GatewayKey = wsServer.Gateway.Key
		model.PushMsgChan <- pushMsg
	}
}
