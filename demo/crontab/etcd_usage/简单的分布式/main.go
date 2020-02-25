package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

//分布式集群下的乐观锁
//lease 实现锁过期
//OP操作
//txn事务 if else then

//1，上锁（创建租约，自动续租，拿着租约去抢占一个key）
//2，处理业务
//3，释放锁（取消自动续租，释放租约）
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

	//1，上锁（创建租约，自动续租，拿着租约去抢占一个key)
	lease := clientv3.NewLease(client)
	leaseGrantesp, _ := lease.Grant(context.TODO(), 5)
	//租约ID
	leaseId := leaseGrantesp.ID

	//函数退出后，自动续租会停止，练习时看代码了为直接一些，都写在main函数中
	ctx, cancelFunc := context.WithCancel(context.TODO())
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

	//5秒后会自动续租
	keepRespChan, _ := lease.KeepAlive(ctx, leaseId)

	//处理续约应答的协程
	go func() {
		for {
			select {
			case keepResp := <-keepRespChan:
				if keepResp == nil {
					fmt.Println("租约失效")
					goto END
				} else {
					fmt.Println("收到租约")
				}
			}
		}
	END:
	}()

	//进行抢key
	lockKey:="/cron/lock/ado"
	kv := clientv3.NewKV(client)
	txn := kv.Txn(context.TODO())
	//if 不存在key,then 设置它，else 抢锁失败
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey,"duzhenxun",clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))
	//提交事务
	txnResp,err:=txn.Commit()
	if err!=nil{
		return
	}
	if !txnResp.Succeeded{
		fmt.Println("锁被占用:",string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	//2,处理业务
	fmt.Println("正在处理业务中。。。。")
	time.Sleep(10*time.Second)
	fmt.Println("业务处理完成，释放锁")

	//3，释放锁(取消自动续租，释放租约)
	//defer里已处理

}
