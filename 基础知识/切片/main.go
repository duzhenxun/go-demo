package main

import (
	"fmt"
)

func main() {
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}

	s1 := arr[2:6] //2,3,4,5
	fmt.Println("s1", s1, len(s1), cap(s1))

	s2 := s1[1:3] //3,4
	fmt.Println("s2", s2, len(s2), cap(s2))

	s3 := s1[4:5] //(s1底层数组应是 [2,3,4,5,6,7,8],所有这里s3应是6)
	fmt.Println("s3", s3, len(s3), cap(s3))

	s4 := s3[:2] //（s3底层应是 6，7，8）这里打印后应该是 6，7
	fmt.Println("s4", s4, len(s4), cap(s4))

	fmt.Println("s4[:3]", s4[:3])
	s4 = append(s4, 9)  //发现追加后，发现底层会改变，由6，7，8变成了6，7，9
	fmt.Println("s4", s4, len(s4), cap(s4))

	fmt.Println("s4[:3]", s4[:3])

}
