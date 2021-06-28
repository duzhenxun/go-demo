package demo

import "time"
// 字符串长度
func strCount(str string) int  {
	var l int
	l = len([]rune(str))
	return l
}
// 计算函数
func algorithm(a int,b int) int {
	time.Sleep(time.Millisecond*100)
	return a+b
}

type Demo struct {

}
// 注释
//		e.g. t.Test1(123)
func(d *Demo) Test1(v int)int{
	return v+v
}

// 说啥返回啥
func(d *Demo) AdoSay(s string) string{
	return s
}

