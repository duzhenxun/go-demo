package main

import (
	"fmt"
	"time"
)

func chanDemo(a int, b int) int {
	time.Sleep(time.Millisecond * 100)
	return a + b
}
func main() {
	type wsMessage struct {
		messageType int
		data        []byte
	}
	ch := make(chan int, 100)
	//b := make(chan wsMessage, 3)
	go func() {
		for {
			select {
			case ch <- 0:
			case ch <- 1:
			}
			time.Sleep(time.Second * 1)
		}
	}()

	for {
		select {
		case i := <-ch:
			fmt.Println("当", i)
		default:
			fmt.Println("无数据")
		}

	}

}
