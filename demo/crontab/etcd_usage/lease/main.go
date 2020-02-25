package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

func main()  {
	var(
		client *clientv3.Client
		err error
	)

	if client,err = clientv3.New(clientv3.Config{
		Endpoints:            []string{"localhost:2379"},
		DialTimeout:          5*time.Second,
	});err!=nil{
		fmt.Println(err)
	}

	//申请一个lease租约
	lease := clientv3.NewLease(client)
	//申请一个10秒的租约
	leaseGrantResp,err := lease.Grant(context.TODO(),10)
	if err!=nil{
		fmt.Println(err)
		return
	}

	//put一个kv,让它与租约联系起来。实现10秒后自动过期
	leaseId:=leaseGrantResp.ID

	//设置一个key,租约使用上面的id
	kv:=clientv3.NewKV(client)
	lockKey:="/cron/lock/job1"
	putResp,err := kv.Put(context.TODO(),lockKey,"",clientv3.WithLease(leaseId))
	if err!=nil{
		fmt.Println(err)
		return
	}
	log.Println(lockKey+" 写入成功",putResp.Header.Revision)

	//测试代码，定期查看一下key是否过期
	for{
		getResp,err:=kv.Get(context.TODO(),lockKey)
		if err!=nil{
			fmt.Println(err)
			return
		}
		if getResp.Count==0{
			log.Println(lockKey+" 过期了")
			break
		}
		log.Println(lockKey+" 还没有过期")
		time.Sleep(2*time.Second)
	}
}
