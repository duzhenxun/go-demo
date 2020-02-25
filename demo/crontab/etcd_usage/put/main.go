package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)
//连接etcd,设置key,获取key
func main(){
	var(
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		putResp * clientv3.PutResponse
		getResp * clientv3.GetResponse
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
	//创建一个key
	job1:="/cron/job1"
	if putResp,err = kv.Put(context.TODO(),job1,"job1,ado",clientv3.WithPrevKV());err!=nil{
		fmt.Println(err)
	}else{
		fmt.Println(putResp.Header.Revision)
		//获取上次一的值
		//fmt.Println(string(putResp.PrevKv.Value))
	}

	//获取key 值
	if getResp,err=kv.Get(context.TODO(),job1);err!=nil{
		fmt.Println(err)
	}else{
		fmt.Println(job1+" 的值是:",getResp.Kvs)
	}


}
