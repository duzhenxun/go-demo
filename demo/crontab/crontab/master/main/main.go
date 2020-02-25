package main

import (
	"fmt"
	"go-demo/demo/crontab/crontab/master"
	"runtime"
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	var (
		err error
	)
	//初始化线程
	initEnv()

	//启动api http服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}
	return
ERR:
	fmt.Println(err)
}
