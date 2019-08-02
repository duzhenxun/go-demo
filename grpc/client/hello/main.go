package main

import (
	"context"
	"flag"
	"fmt"
	"go-demo/grpc/proto/hello"
	"google.golang.org/grpc"
	"time"
)

//go run client/hello/main.go -addr=127.0.0.1:5001 -name=asdfasdf

var addr = flag.String("addr", "127.0.0.1:5000", "register address")
var name = flag.String("name", "duzhenxun", "要发送的名称")

func main() {
	flag.Parse()
	auth := Authentication{
		appKey:    "duzhenxun",
		appSecret: "password",
	}

	//consul 连接
/*	consulName := consul.NewResolver("test.hello.fun2")
	conn, e := grpc.Dial("10.70.120.63:8500", grpc.WithInsecure(), grpc.WithBalancer(grpc.RoundRobin(consulName)), grpc.WithBlock(), grpc.WithPerRPCCredentials(&auth))
	if e != nil {
		panic(e)
	}*/

	//正常grpc常连
	conn, err := grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))
	if err != nil {
		fmt.Println(err.Error())
	}

	defer conn.Close()

	//使用服务
	client := hello.NewHelloServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	/*r, err := client.Fun1(ctx, &hello.Request{Name: *name})
	if err != nil {
		fmt.Printf("%v %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
	}
	fmt.Printf("%v %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Message)*/

	//无需认证
	r2, err := client.Fun2(ctx, &hello.Request{Name: *name})
	if err != nil {
		fmt.Printf("%v %s\n", time.Now().Format("2006-01-02 15:04:05"), err.Error())
	}
	fmt.Printf("%v %s\n", time.Now().Format("2006-01-02 15:04:05"), r2.Message)

	//循环请求
	ticker := time.NewTicker(time.Second * 2)
	for range ticker.C {
		client := hello.NewHelloServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		r, _ := client.Fun1(ctx, &hello.Request{Name: *name})
		fmt.Printf("%v %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Message)
	}

}

type Authentication struct {
	appKey    string
	appSecret string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (
	map[string]string, error,
) {
	return map[string]string{"app_key": a.appKey, "app_secret": a.appSecret}, nil
}
func (a *Authentication) RequireTransportSecurity() bool {
	return false
}
