package mqttClient

import (
	"fmt"
	"time"

	"koudai-box/iot/gateway/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

const WAIT_Time = 3 * time.Second

var cache_ws_Clients = make(map[string]*MqttClient)

type MqttClient struct {
	Client   mqtt.Client
	ClientId string
	Url      string
	UserName string
	Password string
	Shutdown bool
	Gateway  *model.GatewayConfig
}

func StartMQTT(Gateway *model.GatewayConfig) error {
	if nil == Gateway {
		logrus.Errorf("Gateway[%s] is not exist. ", Gateway.Key)
		return fmt.Errorf("gateway[%s] is not exist", Gateway.Key)
	}
	client := MqttClient{
		Gateway:  Gateway,
		Url:      fmt.Sprintf("%s:%d", Gateway.Ip, Gateway.Port),
		ClientId: Gateway.Key,
	}
	cache_ws_Clients[Gateway.Key] = &client
	return client.StartConnect()
}

func StopMQTT(Gateway *model.GatewayConfig) error {
	client, ok := cache_ws_Clients[Gateway.Key]
	if !ok {
		logrus.Errorf("Gateway[%s] is not start. ", Gateway.Key)
		return fmt.Errorf("gateway[%s] is not start", Gateway.Key)
	}
	client.StopConnect()
	return nil
}

func (mqttClient *MqttClient) StartConnect() error {
	if mqttClient.Client != nil {
		mqttClient.StopConnect()
		time.Sleep(WAIT_Time)
	}
	for {
		opts := createClientOptions(mqttClient)
		mqttClient.Client = mqtt.NewClient(opts)
		token := mqttClient.Client.Connect()
		for !token.WaitTimeout(WAIT_Time) {
		}
		if err := token.Error(); err != nil {
			logrus.Errorln("mqtt conn error: ", err)
			time.Sleep(WAIT_Time)
			continue
		} else {
			go mqttClient.healthCheck()
			mqttClient.Shutdown = false
			logrus.Infoln("startConnect ok... mqUrl: ", mqttClient.Url)
			break
		}
	}
	gateway := mqttClient.Gateway
	logrus.Infof("++++++++++++++开启MQTT客户端[%s][%s]++++++++++++++", gateway.Key, mqttClient.Url)
	return nil
}

func (mqttClient *MqttClient) StopConnect() {
	mqttClient.Client.Disconnect(30 * 1000)
	mqttClient.Client = nil
	gateway := mqttClient.Gateway
	logrus.Infof("++++++++++++++关闭MQTT客户端[%s][%s]++++++++++++++", gateway.Key, mqttClient.Url)
}

func (mqttClient *MqttClient) healthCheck() {
	for !mqttClient.Shutdown {
		logrus.Infoln("start healthCheck........., mqUrl:", mqttClient.Url)
		if !mqttClient.Client.IsConnectionOpen() {
			logrus.Infoln("start reConnect........., mqUrl:", mqttClient.Url)
			mqttClient.StartConnect()
		}
		time.Sleep(10 * time.Second)
	}
}

func createClientOptions(mqttClient *MqttClient) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttClient.Url)
	opts.SetUsername(mqttClient.UserName)
	opts.SetPassword(mqttClient.Password)
	opts.SetClientID(mqttClient.ClientId)
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(false)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetResumeSubs(true)
	opts.SetOnConnectHandler(mqttClient.connHandle)
	return opts
}

func (mqttClient *MqttClient) connHandle(cli mqtt.Client) {
	subAllChan(&cli, mqttClient)
}

func subAllChan(cli *mqtt.Client, mqttClient *MqttClient) {
	for _, value := range mqttClient.Gateway.ApiConfigs {
		subTopic(value.Path, &mqttClient.Client, mqttClient)
	}
}

func subTopic(path string, cli *mqtt.Client, mqttClient *MqttClient) {
	listen(path, cli, mqttClient)
}

func listen(topic string, cli *mqtt.Client, mqttClient *MqttClient) {
	(*cli).Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		logrus.Println("message received111", string(msg.Payload()))
		logrus.Println(string(msg.Payload()) == "")
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
		var pushMsg model.PushMsg
		pushMsg.Msg = string(msg.Payload())
		pushMsg.GatewayKey = mqttClient.Gateway.Key
		model.PushMsgChan <- pushMsg
	})
}
