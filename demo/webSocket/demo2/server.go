package main

import (
	"fmt"
	"go-demo/demo/webSocket/demo2/wsConn"
	"net/http"
)

func wsHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	fmt.Println(r.URL)
	fmt.Println(r.Cookies())
	conn, err := wsConn.NewWsConn(w, r)
	if err != nil {
		return
	}

	conn.Init()
}
func sendMsg(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	fmt.Println(query.Get("msg"))
	err := new(wsConn.WsConnection).ApiWriteMsg(1, []byte("收到了"))
	if err != nil {
		fmt.Println(err)
	}
	w.Write([]byte("收到了"))
}

func main() {

	http.HandleFunc("/ws", wsHandle)
	http.HandleFunc("/api", sendMsg)
	http.ListenAndServe(":7777", nil)
}
