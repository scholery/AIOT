package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"koudai-box/conf"

	"koudai-box/iot/db"
	"koudai-box/iot/gateway/model"

	"github.com/sirupsen/logrus"
)

var runLock sync.Mutex
var config *conf.Configuration
var pushedAlarm []string
var pushedEvents []string
var pushedPropss []string

func PushListener() {
	config = conf.GetConf()
	for {
		select {
		case msg := <-model.PushOutMsgChan:
			go PushMsg(msg)
		case msg := <-model.PushStatusChan:
			logrus.Debugf("SyncPushFlag:%+v", msg)
			go SyncPushFlag()
		}
	}
}

func PushMsg(msg model.Message) {
	logrus.Debugf("PushMsg:%+v", msg)
	runLock.Lock()
	defer runLock.Unlock()
	switch msg.Type {
	case model.Message_Type_Iot_Event:
		_, err := PushHttpData(config.IotAlarmPushHttptUrl, "POST", msg, nil)
		if err == nil {
			pushedAlarm = append(pushedAlarm, msg.MsgId)
		} else {
			logrus.Errorf("push alarm to http error:%+v", err)
		}
	case model.Message_Type_Event:
		_, err := PushHttpData(config.IotEventPushHttptUrl, "POST", msg, nil)
		if err == nil {
			pushedEvents = append(pushedEvents, msg.MsgId)
		} else {
			logrus.Errorf("push event to http error:%+v", err)
		}
	case model.Message_Type_Prop:
		_, err := PushHttpData(config.IotPropsPushHttptUrl, "POST", msg, nil)
		if err == nil {
			pushedPropss = append(pushedPropss, msg.MsgId)
		} else {
			logrus.Errorf("push props to http error:%+v", err)
		}
	case model.Message_Type_Device:
		_, err := PushHttpData(config.IotDeviceStatusPushHttptUrl, "POST", msg, nil)
		if err == nil {
			logrus.Errorf("push device to http error:%+v", err)
		}
	}

}

func PushHttpData(address string, method string, data interface{}, params []model.Parameter) (interface{}, error) {
	client := &http.Client{}
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
	runLock.Lock()
	defer runLock.Unlock()
	logrus.Debugf("alarm batch size:%d", len(pushedAlarm))
	if len(pushedAlarm) > 0 {
		num, ok := db.AlarmPushed(pushedAlarm)
		if ok {
			logrus.Debugf("alarm batch size sync:%d", num)
			pushedAlarm = []string{}
		}
	}

	logrus.Debugf("event batch size:%d", len(pushedEvents))
	if len(pushedEvents) > 0 {
		num, ok := db.EventPushed(pushedEvents)
		if ok {
			logrus.Debugf("event batch size sync:%d", num)
			pushedEvents = []string{}
		}
	}
	logrus.Debugf("props batch size:%d", len(pushedPropss))
	if len(pushedPropss) > 0 {
		num, ok := db.DevicdePropsPushed(pushedPropss)
		if ok {
			logrus.Debugf("props batch size sync:%d", num)
			pushedPropss = []string{}
		}
	}
}
