package main

import (
	"context"
	"flag"
	"fmt"
	"go-demo/demo/tracer/c/proto/hello"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

var (
	port = flag.Int("port", 10003, "listening port")
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

	//启动服务
	grpcServer.Serve(lis)

}

//hello服务
type helloServerService struct {
}

func (t *helloServerService) Fun1(ctx context.Context, request *hello.Request) (*hello.Response, error) {

	//调用三方服务
	resp, err := http.Get("http://127.0.0.1:10004/d")
	defer resp.Body.Close()
	fmt.Println(resp.Body, err)

	return &hello.Response{Message: "fun1 hello " + request.Name}, nil
}
func (t *helloServerService) Fun2(ctx context.Context, request *hello.Request) (*hello.Response, error) {
	fmt.Printf("fun2 name:%v \n", request.Name)
	return &hello.Response{Message: "fun2 hello " + request.Name}, nil
}
