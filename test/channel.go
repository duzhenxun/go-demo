package main

import "fmt"

func main() {

	c := make(chan int, 3)

	fmt.Println(len(c))

	c <- 1
	c <- 2

	fmt.Println(len(c))
	len:=len(c)
	for i := 0; i < len; i++ {
		//fmt.Println(<-c)
		v, ok := <-c
		fmt.Println(v, ok)
	}

}
