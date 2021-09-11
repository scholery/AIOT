package test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/goburrow/modbus"
)

const (
	tcpDevice = "localhost:5020"
)

func TestTCPClientAdvancedUsage(t *testing.T) {
	handler := modbus.NewTCPClientHandler(tcpDevice)
	handler.Timeout = 5 * time.Second
	handler.SlaveId = 5
	handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)
	handler.Connect()
	defer handler.Close()

	client := modbus.NewClient(handler)
	//results, err := client.ReadDiscreteInputs(1, 8)
	results, err := client.ReadHoldingRegisters(0, 1)
	if err != nil || results == nil {
		fmt.Println("ReadDiscreteInputs:", err)
	} else {
		fmt.Println("1results:", err)
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
