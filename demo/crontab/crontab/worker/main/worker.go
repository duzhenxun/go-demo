package main

import (
	"flag"
	"fmt"
	"go-demo/demo/crontab/crontab/worker"
	"runtime"
	"time"
)
var (
	confFile string
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initArgs() {
	//worker -config ./conf.json
	flag.StringVar(&confFile, "config", "./conf.json", "conf.json")
	flag.Parse()
}

func main()  {
	var (
		err error
	)
	//初如化参数
	initArgs()

	//初始化线程
	initEnv()

	//加载配置
	if err = worker.InitConfig(confFile); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(worker.GConfig)

	//启动调度器
	if err = worker.InitScheduler();err!=nil{
		return
	}

	//任务管理器
	if err = worker.InitJobMgr(); err != nil {
		fmt.Println(err)
		return
	}

	for{
		time.Sleep(1*time.Second)
	}

	return
}
