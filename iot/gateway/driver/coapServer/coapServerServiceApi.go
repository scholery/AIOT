package coapserver

import (
	"encoding/json"
	"log"
	"net"

	"koudai-box/iot/gateway/model"

	"github.com/dustin/go-coap"
	"github.com/sirupsen/logrus"
)

func RegisterURL() error {
	mux := coap.NewServeMux()

	mux.Handle("/gateway/props", coap.FuncHandler(batchDeviceProps))

	mux.Handle("/gateway/props", coap.FuncHandler(singleDeviceProps))

	mux.Handle("/gateway/events", coap.FuncHandler(batchRecieveEvents))

	mux.Handle("/gateway/deviceKey", coap.FuncHandler(recieveEvents))

	err := coap.ListenAndServe("udp", ":5683", mux)
	if err != nil {
		return err
	}
	return nil
}

//单条推送
func batchDeviceProps(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	log.Printf("Got message in batchDeviceProps: path=%q: %#v from %+v", m.Path(), m, a)
	var request interface{}
	json.Unmarshal(m.Payload, request)
	var pushMsg model.PushMsg
	pushMsg.Msg = request
	pushMsg.Type = model.Msg_Type_Props
	pushMsg.GatewayKey = getPathParam(m)["gatewayKey"]
	model.PushMsgChan <- pushMsg
	if m.IsConfirmable() {
		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   []byte("消息已接收，正在处理"),
		}
		res.SetOption(coap.ContentFormat, coap.TextPlain)

		return res
	}
	return nil
}

//多条推送
func singleDeviceProps(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	log.Printf("Got message in singleDeviceProps: path=%q: %#v from %+v", m.Path(), m, a)
	var request interface{}
	json.Unmarshal(m.Payload, request)
	var pushMsg model.PushMsg
	pushMsg.Msg = request
	pushMsg.Type = model.Msg_Type_Props
	pushMsg.GatewayKey = getPathParam(m)["gatewayKey"]
	model.PushMsgChan <- pushMsg
	if m.IsConfirmable() {
		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   []byte("消息已接收，正在处理"),
		}
		res.SetOption(coap.ContentFormat, coap.TextPlain)

		return res
	}
	return nil
}

//单条推送
func batchRecieveEvents(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	log.Printf("Got message in singleDeviceProps: path=%q: %#v from %+v", m.Path(), m, a)
	var request interface{}
	json.Unmarshal(m.Payload, request)
	parm := getPathParam(m)
	var pushMsg model.PushMsg
	pushMsg.Msg = request
	pushMsg.Type = model.Msg_Type_Events
	pushMsg.GatewayKey = parm["gatewayKey"]
	pushMsg.DeviceKey = parm["deviceKey"]
	model.PushMsgChan <- pushMsg
	if m.IsConfirmable() {
		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   []byte("消息已接收，正在处理"),
		}
		res.SetOption(coap.ContentFormat, coap.TextPlain)

		return res
	}
	return nil
}

//单条推送
func recieveEvents(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	log.Printf("Got message in singleDeviceProps: path=%q: %#v from %+v", m.Path(), m, a)
	var request interface{}
	json.Unmarshal(m.Payload, request)
	parm := getPathParam(m)
	logrus.Debugf("gateway[%s] deviceKey[%s]'s events:", parm["gatewayKey"], parm["deviceKey"], request)
	var pushMsg model.PushMsg
	pushMsg.Msg = request
	pushMsg.Type = model.Msg_Type_Events
	pushMsg.GatewayKey = parm["gatewayKey"]
	pushMsg.DeviceKey = parm["deviceKey"]
	model.PushMsgChan <- pushMsg
	if m.IsConfirmable() {
		res := &coap.Message{
			Type:      coap.Acknowledgement,
			Code:      coap.Content,
			MessageID: m.MessageID,
			Token:     m.Token,
			Payload:   []byte("消息已接收，正在处理"),
		}
		res.SetOption(coap.ContentFormat, coap.TextPlain)

		return res
	}
	return nil
}

// func (m *coap.Message) extractOptions() {
// 	sss := m.Options(15)
// 	for _, v := range sss {
// 		b, _ := v.(string)
// 		log.Printf(b)
// 	}
// }

func getPathParam(m *coap.Message) map[string]string {
	tmp := make(map[string]string)
	tmp["gatewaykey"] = m.Path()[0]
	return tmp
}
