package main

import (
	"fmt"
	"unsafe"
)

func main()  {
	var i1 int = 1
	var i2 int8 = 1
	var i3 int16 = 1
	var i4 int32 = 1
	var i5 int64 = 1
	fmt.Println(unsafe.Sizeof(i1))
	fmt.Println(unsafe.Sizeof(i2))
	fmt.Println(unsafe.Sizeof(i3))
	fmt.Println(unsafe.Sizeof(i4))
	fmt.Println(unsafe.Sizeof(i5))
}
