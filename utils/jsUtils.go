package utils

import (
	"github.com/robertkrimen/otto"
)

func ExecJS(function string, functionName string, param interface{}) (interface{}, error) {
	vm := otto.New()
	_, err := vm.Run(function)
	if err != nil {
		return nil, err
	}
	enc, err := vm.Call(functionName, nil, param)
	result, err := enc.Export()
	return result, err
}
