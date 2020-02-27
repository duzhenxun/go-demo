package master

import (
	"github.com/coreos/etcd/clientv3"
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
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Microsecond,
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
