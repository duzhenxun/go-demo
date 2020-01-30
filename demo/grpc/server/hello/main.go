package main

/**
如果不方便运行客户端代码.直接使用grpcurl 也挺方便
➜  ~ grpcurl -plaintext 127.0.0.1:9080 list
grpc.reflection.v1alpha.ServerReflection
hello.HelloService
➜  ~ grpcurl -plaintext 127.0.0.1:9080 list hello.HelloService
hello.HelloService.Fun1
hello.HelloService.Fun2
➜  ~ grpcurl -plaintext 127.0.0.1:9080 describe hello.HelloService.Fun1
hello.HelloService.Fun1 is a method:
rpc Fun1 ( .hello.Request ) returns ( .hello.Response );
➜  ~ grpcurl -plaintext 127.0.0.1:9080 describe hello.Request
hello.Request is a message:
message Request {
  string name = 1;
}
➜  ~ grpcurl -plaintext 127.0.0.1:9080 describe hello.HelloService.Fun2
hello.HelloService.Fun2 is a method:
rpc Fun2 ( .hello.Request ) returns ( .hello.Response );
➜  ~ grpcurl -plaintext -d '{"name": "gopher"}' 127.0.0.1:9080 hello.HelloService.Fun1
{
  "message": "Token有误!"
}
➜  ~ grpcurl -plaintext -d '{"name": "gopher"}' 127.0.0.1:9080 hello.HelloService.Fun2
{
  "message": "fun2 hello gopher"
}
```
*/
import (
	"context"
	"errors"
	"flag"
	"fmt"
	hello2 "go-demo/demo/grpc/proto/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

var (
	port = flag.Int("port", 5000, "listening port")
)

func main() {
	//解析传入参数
	flag.Parse()

	//注册可用服务,服务中的fun1需要token验证,fun2可以直接访问
	grpcServer := grpc.NewServer()
	hello2.RegisterHelloServiceServer(grpcServer, &helloService{})

	//监听端口
	log.Printf("starting  service at %d", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	// 注册到concul中
	/*	if err = consul.Register("test.hello.fun1", "10.70.120.63", *port, "10.70.120.63:8500", time.Second*10, 15);err!=nil{
			panic(err)
		}
		if err = consul.Register("test.hello.fun2", "10.70.120.63", *port, "10.70.120.63:8500", time.Second*10, 15);err!=nil{
			panic(err)
		}*/

	grpcServer.Serve(lis)

}

// used to implement hello.HelloServiceServer.
type helloService struct {
}

//需要token认证
func (this *helloService) Fun1(ctx context.Context, in *hello2.Request) (*hello2.Response, error) {
	auth := Auth{}
	if err := auth.Check(ctx); err != nil {
		return &hello2.Response{Message: err.Error()}, nil
	}
	//设置时间防止客户端已断开,服务端还在傻傻的执行
	//https://book.eddycjy.com/golang/grpc/deadlines.html
	if ctx.Err() == context.Canceled {
		return nil, errors.New("客户端已断开")
	}
	fmt.Printf("fun1 name:%v\n", in.Name)
	return &hello2.Response{Message: "fun1 hello " + in.Name}, nil
}

//直接可以访问
func (this *helloService) Fun2(ctx context.Context, in *hello2.Request) (*hello2.Response, error) {

	fmt.Printf("fun2 name:%v\n", in.Name)
	return &hello2.Response{Message: "fun2 hello " + in.Name}, nil
}

type Auth struct {
	appKey    string
	appSecret string
}

func (a *Auth) Check(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)

	fmt.Println(md)

	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata.FromIncomingContext err")
	}
	var (
		appKey    string
		appSecret string
	)
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}

	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret() {
		return errors.New("Token is !")
	}

	return nil
}

func (a *Auth) GetAppKey() string {
	return "duzhenxun"
}

func (a *Auth) GetAppSecret() string {
	return "password"
}
