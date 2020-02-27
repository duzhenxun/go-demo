package common

import "encoding/json"

const (
	CodeSuccess int    = 1
	CodeError   int    = -1
	NoAuth      int    = 401
	JobSaveDir  string = "/cron/jobs/"
	JobKillerDir  string = "/cron/killer/"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *Result) GetCode() int {
	return r.Code
}

func (r *Result) SetCode(code int) *Result {
	r.Code = code
	return r
}

func (r *Result) GetMsg() string {
	return r.Msg
}
func (r *Result) SetMsg(msg string) *Result {
	r.Msg = msg
	return r
}
func (r *Result) GetData() interface{} {
	return r.Data
}
func (r *Result) SetData(data interface{}) *Result {
	r.Data = data
	return r
}

func (r *Result) ToJson() string {
	jsons, _ := json.Marshal(r)
	return string(jsons)
}
