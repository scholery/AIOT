package global

import "os"

// 系统退出信号chan
var SystemExitChannel = make(chan os.Signal)
