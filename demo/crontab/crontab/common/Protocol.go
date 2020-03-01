package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"time"
)
//定时任务
type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cron_expr"`
}

//任务调度计划
type JobSchedulePlan struct {
	Job *Job	//要调度的任务
	Expr *cronexpr.Expression //解析好的cronexpr表达式
	NextTime time.Time //下次调度时间
}
//返序列号
func UnpackJob(value []byte) (job *Job, err error) {
	if err = json.Unmarshal(value, &job); err != nil {
		return
	}
	return
}

//变化事件
type JobEvent struct {
	EventType int //1save,2delete
	Job       *Job
}

//任务执行状态
type JobExecuteInfo struct {
	Job *Job
	PlanTime time.Time  //理论调度时间
	RealTime time.Time   //实际调度时间
}

func BuildJobEvent(eventType int, job *Job) *JobEvent {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

//构造任务表达式
func BuildJobSchedulePlan(job * Job)(jobSchedulePlan *JobSchedulePlan,err error){
	
	var(
		expr *cronexpr.Expression
	)
	if expr,err = cronexpr.Parse(job.CronExpr);err!=nil{
		return 
	}
	jobSchedulePlan = &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
	
	return 
}
//构造执行状态信息
func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulePlan)(jobExecuteInfo *JobExecuteInfo){
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulePlan.Job,
		PlanTime: jobSchedulePlan.NextTime,
		RealTime: time.Now(),
	}
	return
}