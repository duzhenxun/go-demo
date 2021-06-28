package main

import "fmt"

func main() {
	//float保留N位小数
	f_score := 3.1415
	s_score := fmt.Sprintf("%.2f", f_score)
	fmt.Println(s_score)
}
