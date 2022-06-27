package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"koudai-box/iot/gateway/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const Retry_Times = 5

var cache_ws_Clients = make(map[string]*WsClient)

type WsClient struct {
	Conns                 map[string]*websocket.Conn
	WsRequestChan         chan string
	WsResponseChan        chan string
	ConsecutivSendFailure int32
	RemoteWsUrl           string
	ResponseSignals       chan struct{}
	Shutdown              bool
	Gateway               *model.GatewayConfig
	// ReadLock sync.Mutex
}

// var WsRequestChan = make(chan string, 128)
func (that *WsClient) ReceiveServerCommand(api string, conn *websocket.Conn) {
	//wsIp := conn.RemoteAddr().String()
	//wsUrl := fmt.Sprintf("ws://%s", wsIp)
	defer func() {
		err := recover()
		if err != nil {
			logrus.Error(err)
		}
		logrus.Warnf("websocket[%s] panic, will restart after 5 seconds", that.Gateway.Key)
		that.abortAndReconnectServer(that.Gateway)
	}()
	for !that.Shutdown {
		i := 0
		for conn == nil {
			time.Sleep(1 * time.Second)
			if i == Retry_Times {
				break
			}
			continue
		}
		if conn == nil {
			continue
		}
		_, msg, err1 := conn.ReadMessage()
		if err1 != nil || len(msg) == 0 {
			logrus.Error(err1)
			continue
		}
		var dataMsg interface{}
		err := json.Unmarshal([]byte(msg), &dataMsg)
		if err != nil {
			logrus.Errorf("ws:gateway[%s]'s msg error,msg:%s", that.Gateway.Key, msg)
			continue
		}

		var pushMsg model.PushMsg
		pushMsg.Msg = dataMsg
		pushMsg.Type = api
		pushMsg.GatewayKey = that.Gateway.Key
		model.PushMsgChan <- pushMsg
	}
}

func (that *WsClient) abortAndReconnectServer(gateway *model.GatewayConfig) {
	StopWsClientConnect(gateway)
	time.Sleep(5 * time.Second)
	if !that.Shutdown {
		StartWsClientConnect(gateway)
	} else {
		logrus.Errorf("wsConnect[%s] has bean shutdown !", gateway.Key)
	}
}

func (that *WsClient) ConnectCheck() {
	for that.Conns != nil {
		if that.ConsecutivSendFailure > 15 { //连续发送错误15次以上， 重连
			atomic.StoreInt32(&that.ConsecutivSendFailure, 0)
			logrus.Warn("连续发送失败15次以上，重连")
			that.abortAndReconnectServer(that.Gateway)
			break
		} else {
			time.Sleep(2 * time.Second)
		}
	}
}

func StartWsClientConnect(gateway *model.GatewayConfig) error {
	times := 0
	for {
		err := initWsClientConnects(gateway)
		times++
		if err != nil {
			if times < Retry_Times {
				time.Sleep(time.Second * 10)
			} else {
				logrus.Errorf("start websocket[%s] timeout,retry %d times", gateway.Key, times)
				return fmt.Errorf("start websocket[%s] timeout,retry %d times", gateway.Key, times)
			}
		} else {
			logrus.Infof("++++++++++++++开启Websocket客户端[%s]应用[%s:%d]++++++++++++++", gateway.Key, gateway.Ip, gateway.Port)
			break
		}
	}
	return nil
}

func initWsClientConnects(gateway *model.GatewayConfig) error {
	if nil == gateway {
		logrus.Errorf("start websocket error,gateway is null.")
		return fmt.Errorf("start websocket error,gateway is null")
	}

	conns := make(map[string]*websocket.Conn)
	getPropApi, ok := gateway.ApiConfigs[model.API_GetProp]
	if !ok {
		logrus.Errorf("GetGateway[%s]'s API %s is null", gateway.Key, model.API_GetProp)
	} else {
		conn := initConn(gateway, getPropApi)
		if conn != nil {
			conns[getPropApi.Name] = conn
		}
	}
	getEventApi, ok := gateway.ApiConfigs[model.API_GetProp]
	if !ok {
		logrus.Errorf("GetGateway[%s]'s API %s is null", gateway.Key, model.API_GetProp)
	} else {
		conn := initConn(gateway, getEventApi)
		if conn != nil {
			conns[getEventApi.Name] = conn
		}
	}

	var ws WsClient
	ws.Conns = conns
	ws.RemoteWsUrl = fmt.Sprintf("ws://%s:%d", gateway.Ip, gateway.Port)
	ws.ResponseSignals = make(chan struct{}, 1)
	ws.Shutdown = false
	ws.Gateway = gateway
	for key, con := range ws.Conns {
		go ws.ReceiveServerCommand(key, con)
	}
	go ws.ConnectCheck()
	cache_ws_Clients[gateway.Key] = &ws
	return nil
}

func initConn(gateway *model.GatewayConfig, api model.ApiConfig) *websocket.Conn {
	url := fmt.Sprintf("ws://%s:%d/%s", gateway.Ip, gateway.Port, api.Path)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logrus.Errorf("websocket[%s]'s api[%s][%s] start connect error:\n%+v", gateway.Key, api.Name, url, err)
		return nil
	}
	return conn
}

func StopWsClientConnect(gateway *model.GatewayConfig) error {
	if nil == gateway {
		logrus.Errorf("start websocket error,gateway is null")
		return fmt.Errorf("start websocket error,gateway is null")
	}
	wsUrl := fmt.Sprintf("ws://%s:%d", gateway.Ip, gateway.Port)
	if client, ok := cache_ws_Clients[gateway.Key]; ok {
		logrus.Infoln("Stop ws connnected:", wsUrl)
		client := client
		client.Shutdown = true
		client.ResponseSignals <- struct{}{}
		for _, con := range client.Conns {
			con.Close()
		}
		client.Conns = nil
		delete(cache_ws_Clients, gateway.Key)
	}
	logrus.Debugf("++++++++++++++关闭Websocket客户端[%s]应用[%s:%d]++++++++++++++", gateway.Key, gateway.Ip, gateway.Port)
	return nil
}

func SendWsMessage(api model.ApiConfig, message string, getewayKey string) error {
	if len(api.Name) == 0 {
		logrus.Errorf("SendWsMessage error,gateway[%s]'s api is null.%+v", getewayKey, api)
		return fmt.Errorf("SendWsMessage error,gateway[%s]'s api is null.%+v", getewayKey, api)
	}
	ws, ok := cache_ws_Clients[getewayKey]
	if !ok {
		logrus.Errorf("SendWsMessage error,gateway[%s] is null.", getewayKey)
		return fmt.Errorf("SendWsMessage error,gateway[%s] is null", getewayKey)
	}
	conn, ok := ws.Conns[api.Name]
	if !ok {
		conn = initConn(ws.Gateway, api)
	}
	if nil == conn {
		logrus.Errorf("SendWsMessage error,gateway[%s]'s api[%s] conn is null", getewayKey, api.Name)
		return fmt.Errorf("SendWsMessage error,gateway[%s]'s api[%s] conn is null", getewayKey, api.Name)
	}
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return err
	}
	return nil
}

/**********************************服务端*****************************************/
var cache_ws_Servers = make(map[string]*model.GatewayConfig)
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitWSServer(gateway *model.GatewayConfig) {

	r := gin.Default()
	r.GET("/Props", Props)

	wsUrl := fmt.Sprintf("%s:%d", gateway.Ip, gateway.Port)
	r.Run(wsUrl)
	cache_ws_Servers[gateway.Key] = gateway
}

func Props(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close() //返回前关闭
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		//写入ws数据
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}
