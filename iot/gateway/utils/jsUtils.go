package utils

import (
	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
)

func ExecJS(function string, functionName string, param interface{}) (interface{}, error) {
	return doExecJS(nil, function, functionName, param)
}

func ExecJSWithContext(function string, functionName string, context interface{}, param interface{}) (interface{}, error) {
	vm := otto.New()
	vm.Set("context", ToMap(context))
	return doExecJS(vm, function, functionName, param)
}

func doExecJS(vm *otto.Otto, function string, functionName string, param interface{}) (interface{}, error) {
	logrus.Debugf("doExecJS funtion name=%s,function\r\n%s", functionName, function)
	if vm == nil {
		vm = otto.New()
	}
	_, err := vm.Run(function)
	if err != nil {
		return nil, err
	}
	enc, err := vm.Call(functionName, nil, ToMap(param))
	if err != nil {
		return nil, err
	}
	result, err := enc.Export()
	return result, err
}
