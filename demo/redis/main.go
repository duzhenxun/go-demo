package main

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/go-redis/redis"
	"time"
)

func main() {
	var addr = "127.0.0.1:6379"
	var password = ""

	//baseRedis(addr,password)
	pool := PoolInitRedis(addr, password)
	time.Sleep(time.Second * 5) //redis中查看有多少连接

	c := pool.Get()
	c2:=pool.Get()
	c3:=pool.Get()
	c4:=pool.Get()
	c5:=pool.Get()
	fmt.Println(c,c2,c3,c4,c5)
	//time.Sleep(time.Second * 5)//redis一共有多少个连接？？
/*	c.Close()
	c2.Close()
	c3.Close()
	c4.Close()
	c5.Close()*/
	//time.Sleep(time.Second*5) //redis一共有多少个连接？？

/*	//下次是怎么取出来的？？
	b1:=pool.Get()
	b2:=pool.Get()
	b3:=pool.Get()
	fmt.Println(b1,b2,b3) //当前都是哪些连接？
	time.Sleep(time.Second*5)
	b1.Close()
	b2.Close()
	b3.Close()*/

	for{
		fmt.Println("主程序运行中....")
		time.Sleep(time.Second*1) //redis一共有多少个连接？？
	}

}


//连接池
// redis pool
func PoolInitRedis(server string, password string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     2,
		IdleTimeout: 30 * time.Second,
		MaxActive:   3,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}
}


//基本写法
func baseRedis(addr string, password string) {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	p, err := c.Ping().Result()
	if err != nil {
		fmt.Println("redis kill")
	}
	fmt.Println(p)

	c.Do("SET", "key", "duzhenxun")
	rs := c.Do("GET", "key").Val()
	fmt.Println(rs)

	time.Sleep(time.Second * 3)
	c.Close()
	time.Sleep(time.Second * 2)
}

