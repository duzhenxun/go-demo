package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-demo/demo/webSocket2/sendMsg"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	content = make([]string, 0) //收到的内容
	mu      sync.Mutex
	addr    = flag.String("addr", ":8080", "register address")
)

func main() {
	flag.Parse()
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/api/send-msg", apiHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Println("启动失败", err)
	}

}
func wsHandler(w http.ResponseWriter, r *http.Request) {
	sendMsg, err := sendMsg.NewSendMsg(w, r)
	if err != nil {
		sendMsg.Close()
		fmt.Println(err)
	}
	for {
		if len(content) > 0 && sendMsg.IsClosed == false {
			mu.Lock()
			var msg = content
			content = make([]string, 0)
			mu.Unlock()
			for _, v := range msg {
				if err := sendMsg.Write([]byte(v)); err != nil {
					sendMsg.Close()
				}
			}
		}
		time.Sleep(time.Second * 1)
	}

}

//写消息接口
func apiHandler(w http.ResponseWriter, r *http.Request) {
	c := r.PostFormValue("content")
	mu.Lock()
	content = append(content, c)
	if len(content) > 10 {
		content = content[:10]
	}
	mu.Unlock()

	res, _ := json.Marshal(map[string]string{
		"code": "1",
		"msg":  fmt.Sprintf("ok,当前信息数量:%v", len(content)),
	})
	w.Write(res)
}
