package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"reflect"
)

var (
	confPath string
)

func init() {
	flag.StringVar(&confPath, "d", "./", " set comet config file path")
}

func main() {

	pwd,_:=os.Getwd()
	confPath:=pwd+"/conf/"
	fmt.Println(confPath)


	//读取yaml
	config := viper.New()
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath(confPath)

	if err := config.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}
	log.Println(config.Get("Author"))
	alise := config.Get("Information.Alise")
	log.Println(alise)
	log.Println(alise.([]interface{})[0].(string))
	for k, v := range alise.([]interface{}) {
		log.Println(k, "====", v)
	}

	//读取toml
	config2 := viper.New()
	config2.SetConfigFile(confPath+"test.toml")
	config2.ReadInConfig()
	fmt.Println(config2.AllKeys())
	fmt.Println(config2.Get("zookeeper"))
	fmt.Println(config2.Get("zookeeper.host"))
	fmt.Println(config2.AllKeys())
	rpcLogicAddrs := config2.Get("rpcLogicAddrs")
	fmt.Println(rpcLogicAddrs)
	fmt.Println(reflect.TypeOf(rpcLogicAddrs))

	var data []map[string]interface{}
	if tmp,err:=json.Marshal(rpcLogicAddrs);err==nil{
		json.Unmarshal(tmp,&data)
	}
	if len(data)>0{
		fmt.Println(data[0])
		fmt.Println(data[0]["addr"])
	}

	//路径问题
	fmt.Println("========================")



}
