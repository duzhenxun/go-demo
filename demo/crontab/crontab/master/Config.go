package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//配置
type Config struct {
	//api接口相关配置
	ApiPort         int `json:"api_port"`
	ApiReadTimeout  int `json:"api_read_timeout"`
	ApiWriteTimeout int `json:"api_write_timeout"`

	//ETCD相关配置
	EtcdEndpoints   []string `json:"etcd_endpoints"`
	EtcdDialTimeout int      `json:"etcd_dial_timeout"`
}

var (
	G_config *Config
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	G_config = &conf
	fmt.Println(G_config)
	return
}
