package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func main() {
	fmt.Println("http://127.0.0.1:12345")
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		s := fmt.Sprintf("你好，世界！ -- Time:%s", time.Now())

		fmt.Fprintf(writer, s)
		ss := request.URL.Query()
		for k, v := range ss {
			fmt.Println("k=" + k + ",v=" + v[0])
		}
		//fmt.Fprintf(writer,ss)sfg

		data := url.Values{}
		data.Set("name", "duzhenxun")
		for k, v := range data {
			fmt.Println(k, v[0])
		}
		//fmt.Println(data.Get("name"))

	})
	http.ListenAndServe(":12345", nil)

}
