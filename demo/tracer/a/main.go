package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/a", func(writer http.ResponseWriter, request *http.Request) {
		s := fmt.Sprintf("接口a,10001 ！ -- Time:%s<br>", time.Now())
		fmt.Fprintf(writer, s)

		//http 连接
		url := "http://127.0.0.1:10002/b"
		res, err := http.Get(url)
		fmt.Fprintf(writer, "请求B接口,res:%v,err:%v<br>", res, err)
	})
	http.ListenAndServe(":10001", nil)
}
