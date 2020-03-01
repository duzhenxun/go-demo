package worker

import (
	"fmt"
	"go-demo/demo/crontab/crontab/common"
	"log"
	"time"
)

type Scheduler struct {
	jobEventChan chan *common.JobEvent	//etcd任务事件队列
	jobPlanTable map[string]*common.JobSchedulePlan		//任务调度计划表
	jobExecutinTable map[string] *common.JobExecuteInfo //任务执行表
}

var (
	GScheduler *Scheduler
)

func InitScheduler() (err error) {
	GScheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.JobSchedulePlan),
		jobExecutinTable:make(map[string] *common.JobExecuteInfo),
	}
	go GScheduler.schedulerLoop()
	return

}

//调度协程
func (s *Scheduler) schedulerLoop() {
	var (
		jobEvent      *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
	)
	//初始化1秒
	scheduleAfter = s.TrySchedule()

	//调度的延时定时器
	scheduleTimer = time.NewTimer(scheduleAfter)

	//定时任务common.Job
	for {
		select {
		case jobEvent = <-s.jobEventChan: //取出任务变化事件
			s.handleJobEvent(jobEvent) //对内存中的维护
		case <-scheduleTimer.C: //最近的任务到期了
		}
		//调度一次任务
		scheduleAfter = s.TrySchedule()
		//重置调度器
		scheduleTimer.Reset(scheduleAfter)
	}
}

//处理任务事件
func (s *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		err             error
		jobSchedulePlan *common.JobSchedulePlan
		jobExisted      bool
	)
	switch jobEvent.EventType {
	case common.JobEventSave:
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		s.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JobEventDelete:
		//map存在值就将其删掉
		if jobSchedulePlan, jobExisted = s.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(s.jobPlanTable, jobEvent.Job.Name)
		}
	}
}
//执行任务
func(s *Scheduler)TryStartJob(jobPlan *common.JobSchedulePlan){
	//调度与执行是2件事情
	var(
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting bool
	)
	//执行任务可能很久,1分钟会调度多次，但只能执行1 次，防并发

	//如果任务正在执行，跳过本次调度
	if jobExecuteInfo,jobExecuting = s.jobExecutinTable[jobPlan.Job.Name];jobExecuting{
		return
	}
	//构建执行状态信息
	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)

	//保存执行状态
	s.jobExecutinTable[jobPlan.Job.Name] = jobExecuteInfo

	//执行任务，启动Shell命令
	log.Println("执行任务，启动Shell命令")
}

//计算任务状态
func (s *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchedulePlan
		now      time.Time
		nearTime *time.Time
	)
	fmt.Println(len(s.jobPlanTable))

	if len(s.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}

	now = time.Now()

	//1，遍历所有任务
	for _, jobPlan = range s.jobPlanTable {

		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//TODO 执行任务
			log.Println("执行任务",jobPlan.Job.Name)
			s.TryStartJob(jobPlan)

			jobPlan.NextTime = jobPlan.Expr.Next(now) //更新下次执行时间

			log.Println("更新下次执行时间",jobPlan.NextTime)
		}
		//统计最近一个要过期的任务时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	//下次调度间隔 （最近要执行的任务调度时间-当前时间）
	scheduleAfter = (*nearTime).Sub(now)
	return

}

//推送作务变化(其它地方推送进来)
func (s *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	s.jobEventChan <- jobEvent
}
