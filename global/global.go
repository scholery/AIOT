package global

import "os"

// 系统退出信号chan
var SystemExitChannel = make(chan os.Signal)

const (
	TIME_TEMPLATE = "2006-01-02 15:04:05"
)
