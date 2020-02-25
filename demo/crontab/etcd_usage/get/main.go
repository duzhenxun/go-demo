package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main()  {
	var(
		client *clientv3.Client
		err error
		kv clientv3.KV
		getResp *clientv3.GetResponse
	)

	if client,err = clientv3.New(clientv3.Config{
		Endpoints:            []string{"localhost:2379"},
		DialTimeout:          5*time.Second,
	});err!=nil{
		fmt.Println(err)
	}

	kv = clientv3.NewKV(client)

	if getResp,err = kv.Get(context.TODO(),"/cron/jobs/",clientv3.WithPrefix());err!=nil{
		fmt.Println(err)
	}else{
		//总个数
		fmt.Println(getResp.Count)
		//分别打出所有的key,value
		for k,v:=range getResp.Kvs{
			fmt.Println(k,string(v.Key),string(v.Value))
		}
	}
}
