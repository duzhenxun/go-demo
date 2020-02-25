package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)
type result struct {
	err error
	output []byte
}
func main()  {
	var(
		ctx context.Context
		cancelFunc context.CancelFunc
		cmd *exec.Cmd
		resultChan chan *result
		res *result
	)
	//创建一个结果对列
	resultChan = make(chan *result,1000)
	//context chan byte
	//cnacelFnc: close(chan byte)
	ctx,cancelFunc  = context.WithCancel(context.TODO())
	go func() {
		var(
			output [] byte
			err error
		)
		//select {case <-ctx.Done}
		// 监听到 select中有<-ctx.Done 就杀掉当前的命令 pid , kill bash
		cmd = exec.CommandContext(ctx,"/bin/bash","-c","sleep 2;echo hello;")

		//执行任务，捕获输出
		output,err = cmd.CombinedOutput()
		//任务输出结果 传给main协程
		resultChan <- &result{
			err:    err,
			output: output,
		}
	}()

	time.Sleep(time.Second*1) //这里1秒后就要将content关闭，所以bash进程会被杀死
	cancelFunc()
	res = <-resultChan

	fmt.Println(res.err,string(res.output))
}
