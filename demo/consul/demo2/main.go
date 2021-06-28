package main

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
)

func main() {
	//regSer()
	findSer()
}

func regSer() {
	config := consul.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := consul.NewClient(config)
	if err != nil {
		fmt.Println("consul client err:", err)
	}

	/*reg := new(consul.AgentServiceRegistration)
	reg.Name = "demo2.test.actapi"
	reg.Port = 443
	reg.Address = "testqaactapi.busi.inke.cn"*/
	// 增加consul健康检查回调函数
	/*	check := new(consul.AgentServiceCheck)
		check.HTTP = fmt.Sprintf("https://%s:%d", reg.Address, reg.Port)
		check.Timeout = "5s"
		check.Interval = "5s"
		check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
		reg.Check = check*/

	client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      "demo2.php.actapi_1",
		Name:    "demo2.php.actapi",
		Tags:    []string{"php"},
		Port:    443,
		Address: "testqaactapi.busi.inke.cn",
	})

	client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      "demo2.php.actapi_2",
		Name:    "demo2.php.actapi",
		Tags:    []string{"php"},
		Port:    80,
		Address: "127.0.0.1",
	})
}

func findSer() {
	client, err := consul.NewClient(&consul.Config{Address: "127.0.0.1:8500"})
	if err != nil {
		fmt.Println("consul client err:", err)
	}
	//srv, _, _ := client.Agent().Service("demo2.php.actapi", nil)
	//fmt.Println(srv.Address, srv.Port)

	//取出健康的
	srv, _, _ := client.Health().Service("test.hello.fun1", "", true, nil)
	for k, v := range srv {
		fmt.Println(k, v.Service.Address, v.Service.Port)
	}

}
