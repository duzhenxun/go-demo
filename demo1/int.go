package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func main()  {
	var a int =100
	fmt.Println(reflect.TypeOf(a))

	b,_:=strconv.ParseFloat(fmt.Sprint(a),64)
	fmt.Println(reflect.TypeOf(b))

}
