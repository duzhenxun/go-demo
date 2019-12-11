package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func main() {

	filePath,_:=os.Getwd()
	filePath=filePath+"/demo"
	fileTmp:=filePath+"/http.go"

	file, _ := os.Open(fileTmp)
	fmt.Println(file.Name())
	defer file.Close()

	fileInfo,_ := os.Stat(fileTmp)
	fmt.Println(fileInfo.Name())
	fmt.Println(fileInfo.Size())
	fmt.Println(fileInfo.IsDir())
	fmt.Println(fileInfo.Mode())
	fmt.Println(fileInfo.ModTime())
	sys:=fileInfo.Sys()
	fmt.Println(sys)
	stat:=sys.(*syscall.Stat_t)
	fmt.Println(time.Unix(stat.Atimespec.Unix())) //上次访问时间

	file2,_:=os.Create(filePath+"/test.txt")
	defer file2.Close()
	file2.Chmod(fileInfo.Mode()|os.ModeSticky)






	fmt.Println(os.Getwd())
}
