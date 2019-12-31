package main

import (
	"flag"
	"fmt"
	"math"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	host := flag.String("host", "kongadmin.com", "hostname")
	startPort := flag.Int("startPort", 80, "start port")
	endPort := flag.Int("endPort", 8888, "end port")
	timeout := flag.Int("timeout", 10000, "timeout 毫秒")
	process := flag.Int("process", 10, "process 个数")
	flag.Parse()

	ports := []int{}
	var wg sync.WaitGroup
	mutex := sync.Mutex{}

	fmt.Printf("%s 开放的端口:\n", *host)

	total := *endPort - *startPort + 1

	pageTotal := math.Ceil(float64(total)/float64(*process))
	all:=map[int][]int{}

	for i:=1;i<=int(pageTotal);i++{
		for j:=1;j<=*process;j++{
			tmp:=(*startPort+i-1)**process+j
			fmt.Println(tmp)
			if tmp<=total{
				all[i]=append(all[i],tmp)
			}
		}

	}


	fmt.Println(all,pageTotal)
	os.Exit(1)

	for port := *startPort; port <= *endPort; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			opened := isOpen(*host, p, *timeout)
			//fmt.Println(opened)
			if opened {
				fmt.Printf("%v\n", p)
				mutex.Lock()
				ports = append(ports, p)
				mutex.Unlock()
			}
		}(port)
	}
	wg.Wait()
	fmt.Printf("%s 开放的端口: %v\n", *host, ports)
}

func isOpen(host string, port int, timeout int) bool {

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Millisecond*time.Duration(timeout))
	if err == nil {
		_ = conn.Close()
		return true
	}
	return false
}
