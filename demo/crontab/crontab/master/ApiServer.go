package master

import (
	"encoding/json"
	"go-demo/demo/crontab/crontab/common"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	G_apiServer *ApiServer
)

func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)

	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/jobs/save", handleJobSave)
	mux.HandleFunc("/jobs/del", handleJobDelete)
	mux.HandleFunc("/jobs/list", handleJobList)
	mux.HandleFunc("/jobs/kill", handleJobKill)

	//静态文件
	staticHandler := http.FileServer(http.Dir(GConfig.WebRoot))
	mux.Handle("/", http.StripPrefix("/", staticHandler))

	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(GConfig.ApiPort)); err != nil {
		return
	}
	//创建一个HTTP服务
	httpServer = &http.Server{
		Handler:      mux,
		ReadTimeout:  time.Duration(GConfig.ApiReadTimeout) * time.Second,
		WriteTimeout: time.Duration(GConfig.ApiWriteTimeout) * time.Second,
	}
	//单例
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	//启动服务
	go httpServer.Serve(listener)

	return
}

//保存任务
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		postJob  string
		jobFiled common.Job
		oldJob   *common.Job
		result   common.Result
	)
	//1解析post表单
	if err = r.ParseForm(); err != nil {
		result.SetCode(common.CodeError).SetMsg("表单解析失败")
		goto ERR
	}

	//取出表单字段
	postJob = r.PostForm.Get("job")
	if err = json.Unmarshal([]byte(postJob), &jobFiled); err != nil {
		result.SetCode(common.CodeError).SetMsg("表单 job 解析失败")
		goto ERR
	}

	//保存到etcd
	if oldJob, err = GJobmgr.SaveJob(&jobFiled); err != nil {
		result.SetCode(common.CodeError).SetMsg(err.Error())
		goto ERR
	}

	//返回正常应答
	result.SetCode(common.CodeSuccess).SetMsg("success").SetData(oldJob)
	w.Write([]byte(result.ToJson()))
	return
ERR:
	w.Write([]byte(result.ToJson()))
	return
}
//删除任务
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		name   string
		result common.Result
		oldJob *common.Job
	)

	if err = r.ParseForm(); err != nil {
		result.SetCode(common.CodeError).SetMsg("xxx")
		goto ERR
	}

	name = r.PostForm.Get("name")
	if name == "" {
		result.SetCode(common.CodeError).SetMsg("name 不能为空")
		goto ERR
	}

	if oldJob, err = GJobmgr.DeleteJob(name); err != nil {
		result.SetCode(common.CodeError).SetMsg(err.Error())
		goto ERR
	}

	result.SetCode(common.CodeSuccess).SetMsg("success").SetData(oldJob)
	w.Write([]byte(result.ToJson()))
	return
ERR:
	w.Write([]byte(result.ToJson()))
	return
}
//列表
func handleJobList(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		result  common.Result
		jobList []*common.Job
	)
	if jobList, err = GJobmgr.ListJobs(); err != nil {
		result.SetCode(common.CodeError).SetMsg(err.Error())
		goto ERR
	}
	result.SetCode(common.CodeSuccess).SetMsg("success").SetData(jobList)
	w.Write([]byte(result.ToJson()))
	return
ERR:
	w.Write([]byte(result.ToJson()))
	return
}
//杀死进程
func handleJobKill(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		result common.Result
		name   string
	)
	if err = r.ParseForm(); err != nil {
		result.SetCode(common.CodeError).SetMsg("xxx")
		goto ERR
	}

	name = r.PostForm.Get("name")
	if name == "" {
		result.SetCode(common.CodeError).SetMsg("name 不能为空")
		goto ERR
	}

	if err = GJobmgr.KillJb(name); err != nil {
		result.SetCode(common.CodeError).SetMsg(err.Error())
		goto ERR
	}
	result.SetCode(common.CodeSuccess).SetMsg("success")
	w.Write([]byte(result.ToJson()))
	return
ERR:
	w.Write([]byte(result.ToJson()))
}
