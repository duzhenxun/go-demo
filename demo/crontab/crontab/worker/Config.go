package worker

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	//ETCD相关配置
	EtcdEndpoints   []string `json:"etcd_endpoints"`
	EtcdDialTimeout int      `json:"etcd_dial_timeout"`
}
var(
	G_config *Config
)
func InitConfig(filename string) (err error)  {
	var(
		content [] byte
		conf Config
	)

	if content,err = ioutil.ReadFile(filename);err!=nil{
		return
	}

	if err = json.Unmarshal(content,&conf);err!=nil{
		return
	}
	return
}