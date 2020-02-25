package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"log"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		getResp *clientv3.GetResponse
		watchRespChan clientv3.WatchChan
	)
	config = clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		log.Println(err.Error())
		return
	}
	defer client.Close()
	kv = clientv3.NewKV(client)
	go func() {
		for{
			kv.Put(context.TODO(), "/cron/jobs/ado", "watch ado")
			kv.Delete(context.TODO(),"/cron/jobs/ado")
			time.Sleep(1*time.Second)
		}

	}()

	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/ado"); err != nil {
		fmt.Println(err)
		return
	}

	if len(getResp.Kvs) != 0 {
		fmt.Println(string(getResp.Kvs[0].Value))
	}

	//当前etcd集群事务ID，单调递增的
	watchStartRevision := getResp.Header.Revision + 1

	//创建一个监听器
	watcher := clientv3.NewWatcher(client)
	//返回一个chan

	//watchRespChan = watcher.Watch(context.TODO(), "/cron/jobs/ado", clientv3.WithRev(watchStartRevision))

	//TODO 这里加一个测试代码，
	//===== 模拟5秒后关闭watch监听 START
	ctx,cancelFunc:=context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})
	watchRespChan = watcher.Watch(ctx,"/cron/jobs/ado", clientv3.WithRev(watchStartRevision))
	//===== 模拟5秒后关闭watch监听 END

	//循环chan中的数据
	for watchResp := range watchRespChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为：", string(event.Kv.Value), "revsion:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}
}
