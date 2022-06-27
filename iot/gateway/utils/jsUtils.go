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
	contextMap := ToMap(context)
	vm.Set("context", contextMap)
	val, err := doExecJS(vm, function, functionName, param)
	contextMap = nil
	return val, err
}

func doExecJS(vm *otto.Otto, function string, functionName string, param interface{}) (interface{}, error) {
	logrus.Debugf("doExecJS funtion[%s],function[%s]", functionName, function)
	if vm == nil {
		vm = otto.New()
	}
	_, err := vm.Run(function)
	if err != nil {
		return nil, err
	}
	paramMap := ToMap(param)
	enc, err := vm.Call(functionName, nil, paramMap)
	paramMap = nil
	if err != nil {
		return nil, err
	}
	result, err := enc.Export()
	return result, err
}
