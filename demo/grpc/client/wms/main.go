package main

import (
	"context"
	"fmt"
	wms_repository2 "go-demo/demo/grpc/proto/wms_repository"
	"google.golang.org/grpc"
)

func main() {
	addr := "10.70.30.121:9080"
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithAuthority("wms"))

	//addr := "service-uwms.fat.n.com:8000"
	// conn, err := grpc.Dial(addr, grpc.WithInsecure())

	defer conn.Close()
	client := wms_repository2.NewRepositoryServiceClient(conn)
	request, err := client.GetRepository(context.Background(), &wms_repository2.Query{SqlQuery: "id=456"})
	if err != nil {
		fmt.Println("fail:",err)
	}
	fmt.Println("ok",request)
}
