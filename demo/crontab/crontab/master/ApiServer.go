package master

import (
	"net"
	"net/http"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}
var(
	G_apiServer *ApiServer
)
func InitApiServer()(err error){
	var(
		mux *http.ServeMux
		listener net.Listener
		httpServer *http.Server
	)

	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save",handleJobSave)
	if listener,err = net.Listen("tcp",":8070");err!=nil{
		return
	}
	//创建一个HTTP服务
	httpServer = &http.Server{
		Handler:           mux,
		ReadTimeout:       5*time.Second,
		WriteTimeout:      5*time.Second,
	}
	//单例
	G_apiServer = &ApiServer{
		httpServer:httpServer,
	}

	//启动服务
	go httpServer.Serve(listener)

	return
}

func handleJobSave(w http.ResponseWriter,r *http.Request){

}