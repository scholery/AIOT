package test

import "fmt"

func ExecRun(s string) {
	i := 1
	for true {
		fmt.Printf("这是第 %d 次执行：%s\n", i, s)
		i++
	}
}
