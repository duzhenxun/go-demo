package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/d", func(writer http.ResponseWriter, request *http.Request) {
		s := fmt.Sprint("接口a,10004 我是最终接口！<br>")
		fmt.Fprintf(writer, s)
	})
	http.ListenAndServe(":10004", nil)
}
