package test

import (
	"fmt"
	"time"
)

func Chann(ch chan int, stopCh chan bool) {
	var i int
	i = 10
	for j := 0; j < i; j++ {
		ch <- j
		time.Sleep(time.Second)
	}
	stopCh <- true
}

func main1() {

	ch := make(chan int)
	c := 0
	stopCh := make(chan bool)

	go Chann(ch, stopCh)

	for {
		select {
		case c = <-ch:
			fmt.Println("Receive", c)
			fmt.Println("channel")
		case s := <-ch:
			fmt.Println("Receive", s)
		case _ = <-stopCh:
			goto end
		}
	}
end:
	test()
}

func test() {
	for m := 1; m < 10; m++ {
		/*    fmt.Printf("第%d次：\n",m) */
		for n := 1; n <= m; n++ {
			fmt.Printf("%dx%d=%d ", n, m, m*n)
		}
		fmt.Println("")
	}
}
