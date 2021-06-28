package main

import (
	"fmt"
	"time"
)

func main() {

	//这里模拟一下休眠
	testA()
}

func testA() {
	startTime := time.Now().Unix()
	timeLimit := int64(10)
	for {
		nowTime := time.Now().Unix()
		if (nowTime - startTime) > timeLimit {
			fmt.Println("end")
			break
		}
		fmt.Println(nowTime)
		time.Sleep(time.Second)
	}
}
