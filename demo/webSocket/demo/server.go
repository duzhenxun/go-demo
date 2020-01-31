package main

import (
	"github.com/gorilla/websocket"
	"go-demo/demo/webSocket/demo/impl"
	"net/http"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {

	wsUp := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsConn, err := wsUp.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	conn:= impl.NewConnection(wsConn)

	for {
		//接收消息，没有消息时 conn.ReadMessage 会阻塞，代码不会继续往下跑
		data, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}

		//发送消息
		if err := conn.WriteMessage(data); err != nil {
			conn.Close()
			return
		}
	}

}

func main() {

	http.HandleFunc("/ws", wsHandler)

	http.ListenAndServe(":7777", nil)
}
