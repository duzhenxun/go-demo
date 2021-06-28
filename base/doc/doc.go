/*
简易计算器计算自定义包
*/
package demo

// 一种实现两个整数相加的函数，
// 返回值为两整数相加之和
func Add(a, b int) int {
	return a + b
}

// 一种实现两个整数相减的函数，
// 返回值为两整数相减之差
func Sub(a, b int) int {
	return a - b
}

// 一种实现两个整数相乘的函数，
// 返回值为两整数相乘之积
func Mul(a, b int) int {
	return a * b
}

//一种实现两个整数相除的函数，
//返回值为两整数相除之商
//
//奇才苛夺as http://www.baidu.com
//
func Div(a, b int) int {
	if b == 0 {
		panic("divide by zero")
	}

	return a / b
}
