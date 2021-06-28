package main

import (
	"fmt"
	"reflect"
)

func main() {
	type AddUserNoble struct {
		DivisionCode int    `json:"division_code" schema:"division_code,required" label:"贵族身份code码：1000玄铁; 1001.青铜；1002.白银；1003.黄金；1004.铂金；1005.钻石；1006.星耀；1007.王者"`
		TransID      string `json:"trans_id" schema:"trans_id" label:"唯一单号" `
		UID          int    `json:"uid" schema:"uid,required" label:"用户id" `
		Duration     int    `json:"duration" schema:"duration,required" label:"时长(天)" `
	}
	n := 12345

	fmt.Println(reflect.ValueOf(n).Kind())

	params := AddUserNoble{}

	v := reflect.ValueOf(params)
	fmt.Println(v.Kind())
	fmt.Println(v.NumField())

	value := 2
	value2 := "2"
	array := []int{1, 2, 3}
	array2 := []string{"1", "2"}
	fmt.Println(value, value2, array, array2)
	fmt.Println("=========")
	targetValue := reflect.ValueOf(array2)
	for i := 0; i < targetValue.Len(); i++ {
		/*if targetValue.Index(i).Interface() == value {
			fmt.Println(targetValue.Index(i).Interface())
		}*/
		if reflect.DeepEqual(value2, targetValue.Index(i).Interface()) == true {
			fmt.Println(targetValue.Index(i).Interface())
		}
	}

	/*switch reflect.ValueOf(value).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == value {
				fmt.Println(targetValue.Index(i).Interface())
			}
		}
		break
	case reflect.String:

	}*/

}
