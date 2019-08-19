package main

import (
	"context"
	"fmt"
	"go-demo/grpc/proto/wms_truck_task"
	"google.golang.org/grpc"
)

func main()  {
	//addr:="feature1-service-uwms.fat.n.com:8000"  //this is ok~
	addr:="10.70.120.79:9080"
	conn,err:=grpc.Dial(addr,grpc.WithInsecure(),grpc.WithAuthority("wms"))

	if err!=nil{
		fmt.Println("conn fail")
		fmt.Println(err)
	}
	defer conn.Close()

	c:=wms_truck_task.NewTruckTaskServiceClient(conn)
	request,err:=c.GetTruckTask(context.Background(),&wms_truck_task.Query{SqlQuery:"id=456"})
	if err !=nil{
		fmt.Println("GetTruckTask fail")
		fmt.Println(err)
	}
	fmt.Println(request)
}