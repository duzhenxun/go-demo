package master

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"go-demo/demo/crontab/crontab/common"
	"log"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

func InitJobMgr() (err error) {
	//配置
	config := clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints, //群集地址
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}

	client, err := clientv3.New(config)
	if err != nil {
		return
	}
	log.Println(client)
	//全局
	G_jobMgr = &JobMgr{
		client: client,
		kv:     clientv3.NewKV(client),
		lease:  clientv3.NewLease(client),
	}
	return
}

//保存
func (j *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	var (
		jobKey   string
		jobValue []byte
		putRest  *clientv3.PutResponse
	)

	jobKey = common.JobSaveDir + job.Name
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}

	//保存到etcd
	if putRest, err = j.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	//如果更新，返回旧新
	if putRest.PrevKv != nil {
		//对旧值进行反序列化
		if err = json.Unmarshal(putRest.PrevKv.Value, &oldJob); err != nil {
			err = nil
			return
		}
	}
	return
}

//删除
func (j *JobMgr) DeleteJob(name string) (oldJob *common.Job, err error) {
	var (
		delResp *clientv3.DeleteResponse
	)
	jobKey := common.JobSaveDir + name
	if delResp, err = j.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}
	//旧的值
	if len(delResp.PrevKvs) != 0 {
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJob); err != nil {
			err = nil
			return
		}
	}
	return
}

//列表
func (j *JobMgr) ListJobs()(jobList []*common.Job, err error) {
	var (
		dirKey string
		getRep *clientv3.GetResponse
	)
	dirKey = common.JobSaveDir
	if getRep, err = j.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); err != nil {
		return
	}
	jobList = make([]*common.Job, 0)
	for _, v := range getRep.Kvs {
		job := &common.Job{}
		if err = json.Unmarshal(v.Value, job); err != nil {
			err = nil
			continue
		}
		jobList = append(jobList, job)
	}
	return
}

//杀死任务
func(j *JobMgr)KillJb(name string) (err error){
	key:=common.JobKillerDir+name

	leaseResp,err:=j.lease.Grant(context.TODO(),1)
	if err!=nil {
		return
	}
	if _,err = j.kv.Put(context.TODO(),key,"",clientv3.WithLease(leaseResp.ID));err!=nil{
		return
	}
	return
}
