package test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/goburrow/modbus"
	modbus2 "github.com/things-go/go-modbus"
)

const (
	tcpDevice = "localhost:502"
	com       = "com3"
)

func TestTCPClientAdvancedUsage(t *testing.T) {
	handler := modbus.NewTCPClientHandler(tcpDevice)
	handler.Timeout = 5 * time.Second
	handler.SlaveId = 1
	handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)
	handler.Connect()
	defer handler.Close()

	client := modbus.NewClient(handler)
	//results, err := client.ReadDiscreteInputs(1, 8)
	results, err := client.ReadHoldingRegisters(5, 1)
	if err != nil || results == nil {
		t.Errorf("ReadDiscreteInputs:%v", err)
	} else {
		t.Logf("1results:%v", results)
	}
	// results, err = client.WriteMultipleRegisters(1, 2, []byte{0, 3, 0, 4})
	// if err != nil || results == nil {
	// 	fmt.Println("WriteMultipleRegisters:", err)
	// } else {
	// 	fmt.Println("2results:", err)
	// }
	// results, err = client.WriteMultipleCoils(5, 10, []byte{4, 3})
	// if err != nil || results == nil {
	// 	fmt.Println("WriteMultipleCoils:", err)
	// } else {
	// 	fmt.Println("3results:", err)
	// }
}

func TestTCPClient2(t *testing.T) {
	p := modbus2.NewTCPClientProvider(tcpDevice, modbus2.WithEnableLogger())
	client := modbus2.NewClient(p)
	err := client.Connect()
	if err != nil {
		t.Errorf("start error.%v", err)
		return
	}
	value, err := client.ReadHoldingRegisters(1, uint16(5), uint16(1))
	if err != nil {
		t.Errorf("err,%v", err)
	} else {
		t.Logf("values:%v", value)
		for i, v := range value {
			fmt.Printf("%d=%#v\n", i+1, int32(v))
		}
	}
}

func TestRTUClient2(t *testing.T) {
	p := modbus2.NewRTUClientProvider()
	p.Address = com
	p.BaudRate = 9600
	p.DataBits = 8
	p.StopBits = 1
	p.Timeout = 1000 * time.Millisecond
	client := modbus2.NewClient(p)
	client.LogMode(true)
	err := client.Connect()
	if err != nil {
		t.Errorf("start error.%v", err)
		return
	}
	value, err := client.ReadHoldingRegisters(1, uint16(5), uint16(1))
	if err != nil {
		t.Errorf("err,%v", err)
	} else {
		t.Logf("values:%v", value)
		for i, v := range value {
			fmt.Printf("%d=%#v\n", i+1, int32(v))
		}
	}
}
