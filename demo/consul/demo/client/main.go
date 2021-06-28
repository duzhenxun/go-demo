package main

import (
	"context"
	"flag"
	"fmt"
	"go-demo/demo/consul"
	"go-demo/demo/consul/demo/proto/hello"
	"google.golang.org/grpc"
	"strings"
	"time"
)

var (
	addr = flag.String("addr", "127.0.0.1:5000", "server address")
	name = flag.String("name", "duzhenxun", "要测试的string")
)

func main() {
	flag.Parse()
	//创建ctx
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	//调用服务
	var conn *grpc.ClientConn
	s := *addr
	p := s[(strings.Index(s, ":") + 1):]
	if p == "8500" {
		//consul 连接
		consulName := consul.NewResolver("test.hello.fun1")
		conn, _ = grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithBalancer(grpc.RoundRobin(consulName)), grpc.WithBlock())
	} else {
		//正常grpc
		conn, _ = grpc.Dial(*addr, grpc.WithInsecure(), grpc.WithAuthority("s1"))
	}

	defer conn.Close()
	client := hello.NewHelloServiceClient(conn)
	r, _ := client.Fun1(ctx, &hello.Request{Name: *name})
	fmt.Println(r.Message)
}
