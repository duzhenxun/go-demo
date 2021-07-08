package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"go-demo/demo/curl/curlWeb/wsConn"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

var (
	wConn *wsConn.WsConnection
)

func main() {
	var (
		port = flag.Int("port", 25252, "端口号")
		uri  = flag.String("uri", "https://testqaactapi.busi.inke.cn/h5/", "地址")
		h    = flag.Bool("h", false, "帮助信息")
	)

	flag.Parse()
	//帮助信息
	if *h == true {
		str := "actHelp version: 活动小助手 v1.0.2(by:DuZhenxun)\n Usage: actHelp [-h] [-p 端口号] \n\nOptions:\n"
		fmt.Fprintf(os.Stderr, str)
		flag.PrintDefaults()
		return
	}
	addr := ":" + strconv.Itoa(*port)

	//===== 打开浏览器 ========//
	serverUri := *uri + "?p=" + strconv.Itoa(*port)
	//系统信息
	osInfo := map[string]interface{}{}
	osInfo["version"] = "1.0.2"
	osInfo["os"] = runtime.GOOS
	osInfo["cpu"] = runtime.NumCPU()
	addrs, _ := net.InterfaceAddrs()
	osInfo["addr"] = fmt.Sprint(addrs)
	osInfo["time"] = time.Now().Unix()
	osInfos, _ := json.Marshal(osInfo)
	osBase := base64.StdEncoding.EncodeToString(osInfos)
	token := hmacSha256(osBase, "dzx")

	serverUri += "Z00X" + base64.StdEncoding.EncodeToString([]byte("&os_base||"+osBase+"&token||"+token))
	//fmt.Println(serverUri)
	openErr := Open(serverUri)
	if openErr != nil {
		fmt.Println(openErr)
	}
	//绑定路由地址
	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/down", downHandle)
	http.HandleFunc("/ws", wsHandle)
	log.Println("http://127.0.0.1" + addr + " 服务已启动...")
	http.ListenAndServe(addr, nil)
}
func indexHandle(w http.ResponseWriter, r *http.Request) {
	s := "活动小助手 v1.0.2 (by:Duzhenxun)"
	w.Write([]byte(s))
}

//下载
func downHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	resp := map[string]interface{}{
		"code": 200,
		"msg":  "ok",
	}
	decoder := json.NewDecoder(r.Body)
	var params map[string]string
	decoder.Decode(&params)
	if _, ok := params["file_name"]; !ok {
		resp["code"] = 201
		resp["msg"] = "缺少字段 file_name"
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}
	if _, ok := params["url"]; !ok {
		resp["code"] = 201
		resp["msg"] = "缺少字段 url"
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}
	if params["file_name"] == "" || params["url"] == "" {
		resp["code"] = 201
		resp["msg"] = "值不能为空"
		b, _ := json.Marshal(resp)
		w.Write(b)
		return
	}

	go start(params["file_name"], params["url"])

	b, _ := json.Marshal(resp)
	w.Write(b)
	return
}

//ws服务
func wsHandle(w http.ResponseWriter, r *http.Request) {
	wsUp := websocket.Upgrader{
		HandshakeTimeout: time.Second * 5,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}
	wsSocket, err := wsUp.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	wConn = wsConn.New(wsSocket)
	for {
		data, err := wConn.ReadMessage()
		fmt.Println(data)
		if err != nil {
			wConn.Close()
			return
		}
		if err := wConn.WriteMessage(data.MessageType, data.Data); err != nil {
			wConn.Close()
			return
		}
	}
}
