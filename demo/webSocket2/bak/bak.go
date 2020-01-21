package bak

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

func main1() {
	flag.Parse()
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/api/sendMsg", apiHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Println("启动失败！" + err.Error())
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	room = r.URL.Query().Get("room")

	// 完成ws协议的握手操作
	wsConnect, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	wsConn, err := impl.InitConnection(wsConnect)
	if err != nil {
		wsConn.Close()
	}


	for {
		//接收用户发来信息，回复信息
		close := wsConn.GetConnClosed()
		if len(content) > 0 && close != true {
			mu.Lock()
			var msg = content
			content = make([]string, 0)
			mu.Unlock()
			for _, v := range msg {
				if err := wsConn.WriteMessage([]byte(v)); err != nil {
					wsConn.Close()
				}
				time.Sleep(100 * time.Millisecond)
			}

		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

//api写数据
func apiHandler(w http.ResponseWriter, r *http.Request) {
	room := r.PostFormValue("room")
	c := r.PostFormValue("content")
	fmt.Println(room,c)
	mu.Lock()
	content = append(content, c)
	if len(content) > 10 {
		content = content[:10]
	}
	mu.Unlock()
	result := map[string]string{
		"code": "1", "msg": fmt.Sprintf("ok,当前信息数量:%v", len(content)),
	}
	res, _ := json.Marshal(result)
	w.Write(res)
}
