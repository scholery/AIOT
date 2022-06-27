package listener

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"koudai-box/conf"

	"koudai-box/iot/db"
	"koudai-box/iot/gateway/model"
	status "koudai-box/iot/gateway/status"

	"github.com/sirupsen/logrus"
)

var runLock sync.Mutex
var config *conf.Configuration
var pushedAlarm [2]int64
var pushedAlarmF []string //失败
var pushedEvents [2]int64
var pushedEventsF []string //失败
var pushedProps [2]int64
var pushedPropsF []string //失败

func PushListener() {
	config = conf.GetConf()
	httpClient = createHTTPClient()
	for status.IsRunning() {
		select {
		case msg := <-model.PushOutMsgChan:
			go PushMsg(msg)
		case <-model.PushStatusChan:
			go SyncPushFlag()
		}
	}
	httpClient = nil
}

func PushMsg(msg model.Message) {
	logrus.Debugf("PushMsg:%+v", msg)
	switch msg.Type {
	case model.Message_Type_Iot_Event:
		_, err := PushHttpData(config.IotAlarmPushHttptUrl, "POST", msg, nil)
		data, ok := msg.Msg.(model.IotEventMessage)
		if err == nil && ok {
			runLock.Lock()
			if pushedAlarm[0] == 0 {
				pushedAlarm[0] = data.Timestamp
			} else {
				pushedAlarm[1] = data.Timestamp
			}
			defer runLock.Unlock()
		} else {
			runLock.Lock()
			pushedAlarmF = append(pushedAlarmF, msg.MsgId)
			defer runLock.Unlock()
			logrus.Errorf("push alarm to http error:%+v", err)
		}
	case model.Message_Type_Event:
		_, err := PushHttpData(config.IotEventPushHttptUrl, "POST", msg, nil)
		data, ok := msg.Msg.(model.EventMessage)
		if err == nil && ok {
			runLock.Lock()
			if pushedEvents[0] == 0 {
				pushedEvents[0] = data.Timestamp
			} else {
				pushedEvents[1] = data.Timestamp
			}
			defer runLock.Unlock()
		} else {
			runLock.Lock()
			pushedEventsF = append(pushedEventsF, msg.MsgId)
			defer runLock.Unlock()
			logrus.Errorf("push event to http error:%+v", err)
		}
	case model.Message_Type_Prop:
		_, err := PushHttpData(config.IotPropsPushHttptUrl, "POST", msg, nil)
		data, ok := msg.Msg.(model.PropertyMessage)
		if err == nil && ok {
			runLock.Lock()
			if pushedProps[0] == 0 {
				pushedProps[0] = data.Timestamp
			} else {
				pushedProps[1] = data.Timestamp
			}
			defer runLock.Unlock()
		} else {
			runLock.Lock()
			pushedPropsF = append(pushedPropsF, msg.MsgId)
			defer runLock.Unlock()
			logrus.Errorf("push props to http error:%+v", err)
		}
	case model.Message_Type_Device:
		_, err := PushHttpData(config.IotDeviceStatusPushHttptUrl, "POST", msg, nil)
		if err != nil {
			logrus.Errorf("push device to http error:%+v", err)
		}
	}

}

func PushHttpData(address string, method string, data interface{}, params []model.Parameter) (interface{}, error) {
	client := httpClient
	var buf *bytes.Buffer = nil
	if nil != data {
		body, err := json.Marshal(data)
		if err != nil {
			logrus.Errorf("api[%s]'s body parse err:", address, err)
			return nil, err
		}
		buf = bytes.NewBuffer(body)
	}

	req, _ := http.NewRequest(strings.ToUpper(method), address, buf)
	for _, param := range params {
		if len(param.Key) == 0 {
			continue
		}
		req.Header.Add(param.Key, param.Value)
	}
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	logrus.Debug("request:", req)
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("http err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		logrus.Errorf("api[%s]'s resp.StatusCode[%d]:", address, resp.StatusCode)
		return nil, fmt.Errorf("api[%s]'s resp.StatusCode[%d]", address, resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("body read err:", err)
		return nil, err
	}
	logrus.Debugf("api[%s]'s response body:%s", address, string(respBody))
	var retData interface{}
	err = json.Unmarshal(respBody, &retData)
	if err != nil {
		logrus.Error("err:", err)
		return nil, err
	}
	return retData, nil
}

func SyncPushFlag() {
	level, _ := logrus.ParseLevel(conf.GetConf().LogLevel)
	if level >= logrus.DebugLevel {
		debugMem()
	}
	runLock.Lock()
	defer runLock.Unlock()
	logrus.Infof("alarm batch pushed:%d-%d,fail:%+v", pushedAlarm[0], pushedAlarm[1], pushedAlarmF)
	if pushedAlarm[0] > 0 {
		num, ok := db.AlarmPushed(pushedAlarm, pushedAlarmF)
		if ok {
			logrus.Infof("alarm batch size sync:%d", num)
			pushedAlarmF = []string{}
			pushedAlarm[0] = 0
			pushedAlarm[1] = 0
		}
	}

	logrus.Infof("event batch pushed:%d-%d,fail:%+v", pushedEvents[0], pushedEvents[1], pushedEventsF)
	if pushedEvents[0] > 0 {
		num, ok := db.EventPushed(pushedEvents, pushedEventsF)
		if ok {
			logrus.Infof("event batch size sync:%d", num)
			pushedEventsF = []string{}
			pushedEvents[0] = 0
			pushedEvents[1] = 0
		}
	}
	logrus.Infof("props batch pushed:%d-%d,fail:%+v", pushedProps[0], pushedProps[1], pushedPropsF)
	if pushedProps[0] > 0 {
		num, ok := db.DevicdePropsPushed(pushedProps, pushedPropsF)
		if ok {
			logrus.Infof("props batch size sync:%d", num)
			pushedPropsF = []string{}
			pushedProps[0] = 0
			pushedProps[1] = 0
		}
	}
}

const (
	MaxIdleConns        int = 100
	MaxIdleConnsPerHost int = 30
	IdleConnTimeout     int = 90
)

var httpClient *http.Client

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
		},

		Timeout: 20 * time.Second,
	}
	return client
}

func debugMem() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	logrus.Infof(`TotalAlloc:%12d,Alloc:%9d,Sys:%9d,Mallocs:%9d,Frees:%9d,HeapAlloc:%9d,HeapSys:%9d,HeapIdle:%9d,HeapInuse:%9d,HeapReleased:%9d,HeapObjects:%9d,GCSys:%9d,OtherSys:%9d,NextGC:%9d,LastGC:%9d,PauseTotalNs:%9d,NumGC:%9d,NumForcedGC:%9d,GCCPUFraction:%24f`,
		ms.TotalAlloc,
		ms.Alloc,
		ms.Sys,
		ms.Mallocs,
		ms.Frees,
		ms.HeapAlloc,
		ms.HeapSys,
		ms.HeapIdle,
		ms.HeapInuse,
		ms.HeapReleased,
		ms.HeapObjects,
		ms.GCSys,
		ms.OtherSys,
		ms.NextGC,
		ms.LastGC,
		ms.PauseTotalNs,
		ms.NumGC,
		ms.NumForcedGC,
		ms.GCCPUFraction)
	debug.FreeOSMemory()
}
