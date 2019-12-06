package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"go-demo/demo/webSocket2/impl"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	addr    = flag.String("addr", ":8080", "register address")
	content = make([]string, 0) //收到的内容
	room    string              //房间号
	mu      sync.Mutex
)

func main() {
	flag.Parse()
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/api", apiHandler)
	if err:=http.ListenAndServe(*addr, nil);err!=nil{
		log.Println(err.Error())
	}else{
		log.Println(*addr)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 完成ws协议的握手操作
	wsConnect, err := upgrader.Upgrade(w, r, nil);
	if err != nil {
		return
	}
	wsConn, err := impl.InitConnection(wsConnect)
	if err != nil {
		wsConn.Close()
	}
	for {
		//接收用户发来信息，回复信息
		/*data, err := wsConn.ReadMessage()
		if err != nil {
			wsConn.Close()
		}
		fmt.Println(data)*/
		//wsConn.WriteMessage(data)

		//fmt.Println(len(content))
		if len(content) > 0 {
			mu.Lock()
			var msg = content
			content = make([]string, 0)
			mu.Unlock()
			for _, v := range msg {
				time.Sleep(1*time.Second)
				if err := wsConn.WriteMessage([]byte(v)); err != nil {
					wsConn.Close()
				}
			}
		} else {
			//wsConn.WriteMessage([]byte("还没有信息"))
			time.Sleep(1*time.Second)
		}
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	room = r.PostFormValue("room")
	content = append(content, r.PostFormValue("content"))
	if len(content) > 5 {
		content = content[:5]
	}
	mu.Unlock()
	result := map[string]string{
		"code": "1", "msg": fmt.Sprintf("ok,当前信息数量:%v",len(content)),
	}
	res, _ := json.Marshal(result)
	w.Write(res)
}
