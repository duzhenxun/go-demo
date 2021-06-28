package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("info.txt")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()

	fmt.Println("read")
	scan := bufio.NewScanner(f)

	for scan.Scan(){
		lineText :=scan.Text()

		if lineText!=""{
			fmt.Println(lineText)
		}
		fmt.Println("======")
	}


}
