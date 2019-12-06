package webSocket1

import (
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `data`
}

var (
	content = make([]string, 0)
	room    string
	mu      sync.Mutex
	addr    = flag.String("addr", ":8080", "register address")
)

func main() {
	flag.Parse()
	http.Handle("/ws", websocket.Handler(WsHandler))
	http.HandleFunc("/api", ApiHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

func WsHandler(w *websocket.Conn) {
	for {
		if len(content) == 0 {
			//continue
		}
		mu.Lock()
		var msg = content
		content = make([]string, 0)
		mu.Unlock()
		for _, v := range msg {
			if err := websocket.Message.Send(w, v); err != nil {
				fmt.Println(err)
			}
		}
		/*var reply string
		if err:=websocket.Message.Receive(w,&reply);err!=nil{
			fmt.Println("不能够接受消息 error==",err)
			break
		}*/
		time.Sleep(time.Second * 2)
	}
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	room = r.PostFormValue("room")
	content = append(content, r.PostFormValue("content"))
	if len(content) > 10 {
		content = content[:10]
	}
	fmt.Println(len(content))
	mu.Unlock()
	result := Result{}
	result.Code = 1
	result.Msg = "ok"
	res, _ := json.Marshal(result)
	w.Write(res)
}
