package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

func main() {
	var (
		client *clientv3.Client
		err    error
	)

	if client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}); err != nil {
		fmt.Println(err)
	}

	//申请一个lease租约
	lease := clientv3.NewLease(client)
	//申请一个10秒的租约
	leaseGrantResp, err := lease.Grant(context.TODO(), 10)
	if err != nil {
		fmt.Println(err)
		return
	}

	leaseId := leaseGrantResp.ID
	keepRespChan, err := lease.KeepAlive(context.TODO(), leaseId)
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		for {
			select {
			case keepResp := <-keepRespChan:
				if keepResp == nil {
					log.Println("租约失效.服务器原因或其它的原因...")
					goto END
				} else {
					log.Println("收到续租应答", keepResp.ID)
				}
			}
		}
	END:
	}()

	//设置一个key,租约使用上面的id
	kv := clientv3.NewKV(client)
	lockKey := "/cron/lock/job1"
	putResp, err := kv.Put(context.TODO(), lockKey, "", clientv3.WithLease(leaseId))
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(lockKey+" 写入成功", putResp.Header.Revision)

	//测试代码，定期查看一下key是否过期
	for {
		getResp, err := kv.Get(context.TODO(), lockKey)
		if err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			log.Println(lockKey + " 过期了")
			break
		}
		log.Println(lockKey + " 还没有过期")
		time.Sleep(2 * time.Second)
	}
}
