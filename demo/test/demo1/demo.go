package main

import (
	"fmt"
	"runtime"
)

//字符串长度
func strCount(str string) int {
	var l int
	l = len([]rune(str))
	return l
}

func strCountB(str string) int {
	var l int
	l = len([]rune(str))
	return l
}

func a() {
	str := map[string]string{
		"121": "11",
		"ddd": "111",
	}
	fmt.Println(str["999"])
	panic("aaaa")
}
func b() {
	defer func() {
		err := recover()
		switch err.(type) {
		case runtime.Error: // 运行时错误
			fmt.Println("runtime error:", err)
		default: // 非运行时错误
			fmt.Println("error:", err)
		}

		fmt.Println("我没事呵呵呵")
	}()
	a()
	fmt.Println("bbbbbbbb")
}
