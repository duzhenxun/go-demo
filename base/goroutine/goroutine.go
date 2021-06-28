package main

import (
	"time"
)

func main()  {
	var a [10] int
	for i:=0;i<10;i++{
		go func(k int) {
			for{
				a[k]++
				//runtime.Gosched()//交出控制权
			}
		}(i)
	}
	time.Sleep(time.Second*10)
	//fmt.Println(a)
}