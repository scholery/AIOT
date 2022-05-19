package modbusClient

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"koudai-box/iot/gateway/model"

	"github.com/sirupsen/logrus"
	modbus "github.com/things-go/go-modbus"
)

const (
	RETRY    = 3 //重试次数
	TIME_OUT = 1 //秒
)

//httpserver
var cache_servers map[string]modbus.Client = make(map[string]modbus.Client)

func ConnectModbus(gateway *model.GatewayConfig) error {
	switch gateway.Protocol {
	case model.Geteway_Protocol_ModbusRTU:
		return ConnectModbusRTU(gateway)
	case model.Geteway_Protocol_ModbusTCP:
		return ConnectModbusTCP(gateway)
	default:
		logrus.Errorf("ConnectModbus gateway[%s] modbus protocol error:%s", gateway.Key, gateway.Protocol)
	}
	return fmt.Errorf("ConnectModbus gateway[%s] modbus protocol error:%s", gateway.Key, gateway.Protocol)
}
func ConnectModbusRTU(gateway *model.GatewayConfig) error {
	if len(gateway.ModbusConfig.Com) == 0 {
		logrus.Errorf("gateway[%s]'s com is null", gateway.Key)
		return fmt.Errorf("gateway[%s]'s com is null", gateway.Key)
	}
	p := modbus.NewRTUClientProvider()
	p.Address = gateway.ModbusConfig.Com
	if gateway.ModbusConfig.BaudRate > 0 {
		p.BaudRate = gateway.ModbusConfig.BaudRate
	} else {
		p.BaudRate = 9600
	}
	if gateway.ModbusConfig.DataBits > 0 {
		p.DataBits = gateway.ModbusConfig.DataBits
	} else {
		p.DataBits = 8
	}
	if len(gateway.ModbusConfig.Parity) > 0 {
		p.Parity = gateway.ModbusConfig.Parity
	} else {
		p.Parity = "N"
	}
	if gateway.ModbusConfig.StopBits > 0 {
		p.StopBits = gateway.ModbusConfig.StopBits
	} else {
		p.StopBits = 1
	}
	p.Timeout = TIME_OUT * time.Second
	client := modbus.NewClient(p)
	client.LogMode(true)
	err := client.Connect()
	if err != nil {
		logrus.Errorf("ConnectModbusRTU gateway[%s] start error.dev[%s].%+v", gateway.Key, gateway.ModbusConfig.Com, err)
		return fmt.Errorf("ConnectModbusRTU gateway[%s] start error.dev[%s].%+v", gateway.Key, gateway.ModbusConfig.Com, err)
	}
	cache_servers[gateway.Key] = client
	logrus.Infof("++++++++++++++开启ModbusRTU客户端[%s][%s]++++++++++++++", gateway.Key, gateway.ModbusConfig.Com)
	return nil
}

func ConnectModbusTCP(gateway *model.GatewayConfig) error {
	url := fmt.Sprintf("%s:%d", gateway.Ip, gateway.Port)
	p := modbus.NewTCPClientProvider(url, modbus.WithEnableLogger())
	client := modbus.NewClient(p)
	err := client.Connect()
	if err != nil {
		logrus.Errorf("ConnectModbusTCP gateway[%s] start error.url[%s].%+v", gateway.Key, url, err)
		return fmt.Errorf("ConnectModbusTCP gateway[%s] start error.url[%s].%+v", gateway.Key, url, err)
	}
	cache_servers[gateway.Key] = client
	logrus.Infof("++++++++++++++开启ModbusTCP客户端[%s][%s]++++++++++++++", gateway.Key, url)
	return nil
}

func CloseModbus(gateway *model.GatewayConfig) error {
	client, ok := cache_servers[gateway.Key]
	if !ok {
		logrus.Errorf("modbus[%s]'s client is not exist.", gateway.Key)
		return fmt.Errorf("modbus[%s]'s client is not exist", gateway.Key)
	}
	client.Close()
	logrus.Infof("++++++++++++++关闭ModbusTCP/RTP客户端[%s]++++++++++++++", gateway.Key)
	return nil
}

func QueryValue(gateway *model.GatewayConfig, device *model.Device, item model.ItemConfig) (interface{}, error) {
	client, ok := cache_servers[gateway.Key]
	if !ok {
		err := ConnectModbus(gateway)
		client, ok = cache_servers[gateway.Key]
		if err != nil || !ok {
			logrus.Errorf("modbus[%s]'s client is not exist.", gateway.Key)
			return nil, fmt.Errorf("modbus[%s]'s client is not exist", gateway.Key)
		}
	}
	slaveId, err := strconv.ParseUint(device.SourceId, 16, 32)
	if err != nil {
		return nil, fmt.Errorf("modbus[%s]'s salveId[%s] is error,address[%x]", gateway.Key, device.SourceId, item.Address)
	}
	start, err := strconv.ParseUint(item.Address, 16, 32)
	if err != nil {
		return nil, fmt.Errorf("modbus[%s]'s salveId[%s-%x],address[%s] is error", gateway.Key, device.SourceId, slaveId, item.Address)
	}
	quantity := 1
	if item.Quantity > 0 {
		quantity = item.Quantity
	}
	value, err := client.ReadHoldingRegisters(byte(slaveId), uint16(start), uint16(quantity))
	if err != nil {
		//	错误重试
		for i := 1; i <= RETRY && err != nil && strings.Contains(err.Error(), "timeout"); i++ {
			time.Sleep(TIME_OUT * time.Second)
			logrus.Errorf("modbus[%s]'s salveId[%x],address[%x] timeout retry:%d", gateway.Key, slaveId, start, i)
			value, err = client.ReadHoldingRegisters(byte(slaveId), uint16(start), uint16(quantity))
		}
	}
	if err != nil {
		logrus.Errorf("modbus[%s]'s salveId[%x],address[%x] ReadHoldingRegisters err,%+v", gateway.Key, slaveId, start, err)
		return nil, fmt.Errorf("modbus[%s]'s salveId[%x],address[%x] ReadHoldingRegisters err,%+v", gateway.Key, slaveId, start, err)
	}
	logrus.Debugf("modbus[%s]'s salveId[%x],address[%x] quantity[%d],value[%+v]", gateway.Key, slaveId, start, quantity, value)
	if len(value) == 0 {
		logrus.Errorf("modbus[%s]'s salveId[%x],address[%x] values null.", gateway.Key, slaveId, start)
		return nil, fmt.Errorf("modbus[%s]'s salveId[%x],address[%x] value null", gateway.Key, slaveId, start)
	} else if len(value) == 1 {
		return value[0], nil
	} else {
		return value, nil
	}
}
