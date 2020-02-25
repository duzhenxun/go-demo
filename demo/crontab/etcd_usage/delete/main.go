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
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		delResp * clientv3.DeleteResponse
	)
	config = clientv3.Config{
		Endpoints:            []string{"localhost:2379"},
		DialTimeout:          5*time.Second,
	}

	if client,err=clientv3.New(config);err!=nil{
		log.Println(err.Error())
		return
	}
	defer client.Close()

	kv = clientv3.NewKV(client)
	if delResp,err = kv.Delete(context.TODO(),"/cron/job1",clientv3.WithPrevKV());err!=nil{
		fmt.Println(err)
		return
	}
	if len(delResp.PrevKvs)!=0{
		for k,v:=range delResp.PrevKvs{
			fmt.Println(k,string(v.Key),string(v.Value))
		}
	}

}