package main

import (
	"fmt"
)

func main() {
	//要求截取前4个字符串，最后是"xs25"
	//var str = "xs25.cn"
	//fmt.Println("xs25.cn截取前4个字符串的结果是：",str[0:4])
	//运行后的结果：xs25.cn截取前4个字符串的结果是： xs25
	//要求截取前4个字符串，最后是"小手25"
	var str = "小手25是什么"
	fmt.Println(len(str))
	fmt.Println(len([]byte(str)))
	fmt.Println(len([]uint8(str)))


	fmt.Println("小手25是什么 ,截取前4个字符串的结果是：", str[:4])

	fmt.Println("小手25是什么 ,截取前4个字符串的结果是：", str[:8])

	s := []int32(str)
	fmt.Println(len(s))
	fmt.Println("===============")
	fmt.Println(string(s[:5]))

	ss := []string{"ado", "adoduzhenxun", "小手25是什么", "小手25小手舞动起来", "杜zhenxun"}

	maxLenStr := ""
	for i := 0; i < len(ss); i++ {
		var repeat bool
		tmpArr := map[int32]int{}
		for k, v := range []rune(ss[i]) {
			if tmpArr[v] != 0 && len(tmpArr) > 0 {
				repeat = true //有重复
				break
			}
			tmpArr[v] = k
		}
		//没有重复找最长的
		if !repeat && len(ss[i])>len(maxLenStr){
			maxLenStr = ss[i]
		}
	}
	fmt.Println("无重复最长的是：",maxLenStr)

	//结果：小手25是什么

	/*	fmt.Println(utf8.RuneCountInString(s))
		fmt.Println("======")
		fmt.Println()
		b:=[]byte(s)
		for len(b)>0{
			ch,size:=utf8.DecodeRune(b)
			b = b[size:]
			fmt.Printf("%c \n",ch)
		}*/

	/*
		for _,b:=range []byte(s){
			fmt.Printf("%X \n",b)
		}
		for k,v:=range s{
			fmt.Printf("%d %X \n",k,v)
		}*/
	/*
		for i,ch:=range[]rune(s){
			fmt.Printf("(%d %c)",i,ch)
		}*/
}

func lengthOf(s string) int {
	diffArr := make(map[byte]int)
	for i, ch := range []byte(s) {
		if l, ok := diffArr[ch]; ok {
			fmt.Println("=====", l)
		}
		fmt.Println(i, ch)
		diffArr[ch] = i
	}
	return 0
}
