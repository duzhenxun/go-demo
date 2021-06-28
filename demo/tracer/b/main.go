package main

import (
	"context"
	"fmt"
	"go-demo/demo/tracer/b/proto/hello"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/b", func(writer http.ResponseWriter, request *http.Request) {
		s := fmt.Sprintf("接口b,10002 ！ -- Time:%s<br>", time.Now())
		fmt.Fprintf(writer, s)
		//调用DGrpc
		ctx := context.Background()

		//调用服务
		var conn *grpc.ClientConn
		conn, _ = grpc.Dial("127.0.0.1:10003", grpc.WithInsecure(), grpc.WithAuthority("s1"))
		defer conn.Close()
		client := hello.NewHelloServiceClient(conn)
		r, _ := client.Fun1(ctx, &hello.Request{Name: "duzhenxun"})
		fmt.Println(r.Message)
	})
	http.ListenAndServe(":10002", nil)
}
