package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time
}

func main() {

	var (
		cronJob       *CronJob
		expr          *cronexpr.Expression //表达式
		now           time.Time
		scheduleTable map[string]*CronJob //任务表
	)
	scheduleTable = make(map[string]*CronJob)
	//当前时间
	now = time.Now()
	//Cron表达式
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	//任务注册到任务表中
	scheduleTable["job1"] = cronJob

	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	//任务注册到任务表中
	scheduleTable["job2"] = cronJob

	// 需要有1个调度协程，它定时检查所有的Cron任务，谁过期就执行谁
	go func() {
		var (
			jobName string
			cronJob *CronJob
			now     time.Time
		)
		//定时看任务表
		for {
			now = time.Now()
			for jobName, cronJob = range scheduleTable {
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					//启动一个协程，执行这个任务
					go func(jobName string) {
						fmt.Println("执行:", jobName)
					}(jobName)

					//计算下一次执行时间
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName, "下次执行时间:", cronJob.nextTime)
				}
			}

			//这里模拟一下休眠
			select {
			case <-time.NewTimer(time.Millisecond*100).C: //100毫秒可读，返回

			}
		}
	}()

	time.Sleep(time.Second*100)
}
