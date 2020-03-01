package worker

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go-demo/demo/crontab/crontab/common"
	"log"
	"strings"
	"time"
)

type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	GJobmgr *JobMgr
)

func InitJobMgr() (err error) {
	//配置
	config := clientv3.Config{
		Endpoints:   GConfig.EtcdEndpoints, //群集地址
		DialTimeout: time.Duration(GConfig.EtcdDialTimeout) * time.Millisecond,
	}

	client, err := clientv3.New(config)
	if err != nil {
		log.Println(err)
		return
	}
	//全局
	GJobmgr = &JobMgr{
		client:  client,
		kv:      clientv3.NewKV(client),
		lease:   clientv3.NewLease(client),
		watcher: clientv3.NewWatcher(client),
	}
	// 启动任务监听
	GJobmgr.watchJobs()
	return
}

//监听变化
func (j *JobMgr) watchJobs() (err error) {
	var (
		getResp  *clientv3.GetResponse
		job      *common.Job
		jobEvent *common.JobEvent
	)
	//获取所有任务
	if getResp, err = j.kv.Get(context.TODO(), common.JobSaveDir, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, v := range getResp.Kvs {
		//反序列化
		if job, err = common.UnpackJob(v.Value); err == nil {

			jobEvent = &common.JobEvent{
				EventType: common.JobEventSave,
				Job:       job,
			}

			//TODO: 把这个任务同步给调度协程
			log.Println(jobEvent)
			GScheduler.PushJobEvent(jobEvent)
		}
	}

	//监听任务变化
	go func() {
		watchChan := j.watcher.Watch(context.TODO(), common.JobSaveDir, clientv3.WithRev(getResp.Header.Revision+1), clientv3.WithPrefix())
		//处理监听事件
		for watchResp := range watchChan {
			for _, watchEvent := range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					//TODO 推一个更新事件
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					jobEvent = &common.JobEvent{
						EventType: common.JobEventSave,
						Job:       job,
					}

				case mvccpb.DELETE:
					//TODO 推一个删除事件
					//Delete /cron/jobs/ado
					jobName := strings.TrimPrefix(string(watchEvent.Kv.Key), common.JobSaveDir)

					job = &common.Job{
						Name: jobName,
					}
					//构造事件
					jobEvent = &common.JobEvent{
						EventType: common.JobEventDelete,
						Job:       job,
					}
				}
				log.Println(jobEvent)
				GScheduler.PushJobEvent(jobEvent)
			}
		}
	}()

	return
}
