package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"go-demo/demo/consul/demo/proto/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var (
	port = flag.Int("port", 5000, "listening port")
)

func main() {
	flag.Parse()
	//监听端口
	log.Printf("starting  service at %d", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	hello.RegisterHelloServiceServer(grpcServer, &helloServerService{})

	// 反射
	reflection.Register(grpcServer)

	//注册服务到consul中
	consulAddress := "127.0.0.1:8500"
	// 注册到consul中,查看  curl http://127.0.0.1:8500/v1/catalog/service/test.hello.fun1
	client, err := consul.NewClient(&consul.Config{
		Address: consulAddress,
	})
	if err != nil {
		fmt.Println("consul client err:", err)
	}
	err = client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Kind:    "",
		ID:      "",
		Name:    "test.hello.fun1",
		Tags:    []string{"grpc"},
		Port:    5001,
		Address: "127.0.0.1",
	})
	if err != nil {
		fmt.Println("consul client err:", err)
	}

	//启动服务
	grpcServer.Serve(lis)

}

//hello服务
type helloServerService struct {
}

func (t *helloServerService) Fun1(ctx context.Context, request *hello.Request) (*hello.Response, error) {
	//设置时间防止客户端已断开,服务端还在傻傻的执行
	//https://book.eddycjy.com/golang/grpc/deadlines.html
	if ctx.Err() == context.Canceled {
		return nil, errors.New("客户端已断开")
	}
	fmt.Printf("fun1 name:%v \n", request.Name)
	return &hello.Response{Message: "fun1 hello " + request.Name}, nil
}
func (t *helloServerService) Fun2(ctx context.Context, request *hello.Request) (*hello.Response, error) {
	fmt.Printf("fun2 name:%v \n", request.Name)
	return &hello.Response{Message: "fun2 hello " + request.Name}, nil
}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
