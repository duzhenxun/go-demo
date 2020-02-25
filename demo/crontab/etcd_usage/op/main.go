package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		opResp clientv3.OpResponse
	)

	if client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}); err != nil {
		fmt.Println(err)
	}
	kv = clientv3.NewKV(client)

	//执行OP
	if opResp, err = kv.Do(context.TODO(), clientv3.OpPut("/cron/jobs/ado", "duzhenxun")); err != nil {
		fmt.Println(err)
		return
	}
	//写入的信息版本
	fmt.Println(opResp.Put().Header.Revision)

	//读取数据
	//kv.Get("/cron/jobs/ado")
	//OP操作
	if opResp, err = kv.Do(context.TODO(), clientv3.OpGet("/cron/jobs/ado")); err != nil {
		fmt.Println(err)
		return
	}
	//读取的到信息
	fmt.Println(string(opResp.Get().Kvs[0].Value))
}
