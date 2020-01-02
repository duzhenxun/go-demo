package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go run main.go -startPort 1 -endPort 50000 -process 1000 -host kongadmin.com -timeout 100
func main() {
	ip := flag.String("ip", "127.0.0.1", "ip地址如： 192.168.0.1-255 或直接输入域名 kongadmin.com")
	startPort := flag.Int("startPort", 80, "开始端口号如：-start 80 ")
	endPort := flag.Int("endPort", 88, "结束端口号如：-end 80 ")
	timeout := flag.Int("timeout", 200, "超时时长(毫秒)如: -timeout 200")
	process := flag.Int("process", 100, "进程数如：-process 10")
	path := flag.String("path", "log", "进程数如：-path log")
	flag.Parse()

	//创建目录
	os.Mkdir(*path, os.ModePerm)

	fmt.Printf("========== %v 开始执行任务，ip:%v,start:%v,end:%v,timeout:%v,process:%v \n", time.Now().Format("2016-01-02 15:04:05"), *ip, *startPort, *endPort, *timeout, *process)
	start := time.Now()
	scanIP := ScanIp{}
	hosts, err := scanIP.getAllIp(*ip)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//扫所有的ip
	fileName:=*path+"/"+*ip+"_port.txt"

	for i := 0; i < len(hosts); i++ {
		ports:=scanIP.getHostOpenPort(hosts[i], *startPort, *endPort, *timeout, *process)
		if len(ports)>0{
			f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
			if err!=nil{
				fmt.Println(err)
				f.Close()
				continue
			}
			var str = fmt.Sprintf("%v ip:%v,开放端口:%v \n",time.Now().Format("2016-01-02 15:04:05"),hosts[i],ports)
			if _,err:=f.WriteString(str);err!=nil{
				f.Close()
				continue
			}
		}
	}

	fmt.Printf("========== %v 所有操作完成，总执行时长：%.2fs \n", time.Now().Format("2016-01-02 15:04:05"), time.Since(start).Seconds())
}














//ip 扫描
type ScanIp struct {
}

func (s *ScanIp) getHostOpenPort(host string, startPort int, endPort int, timeout int, process int) ([]int) {
	var (
		total     int
		pageCount int
		pageTotal int
		ports     []int
		wg        sync.WaitGroup
		mutex     sync.Mutex
	)

	total = endPort - startPort + 1
	if total < process {
		pageCount = total
	} else {
		pageCount = process
	}
	pageTotal = int(math.Ceil(float64(total) / float64(process)))
	//fmt.Printf("\n%v 【%v】端口号（%v-%v），总数:%v 个，总协程:%v 个，每个协程处理:%v 个，超时时间:%v毫秒 \n", time.Now().Format("2006-01-02 15:04:05"), host, startPort, endPort, total, pageCount, pageTotal, timeout)
	start := time.Now()
	all := map[int][]int{}
	for i := 1; i <= pageCount; i++ {
		for j := 0; j < pageTotal; j++ {
			tmp := (i-1)*pageTotal + j + startPort
			if tmp <= endPort {
				all[i] = append(all[i], tmp)
			}
		}
	}

	for k, v := range all {
		wg.Add(1)
		go func(value []int, key int) {
			defer wg.Done()
			//start := time.Now()
			var tmpPorts []int
			for i := 0; i < len(value); i++ {
				opened := s.isOpen(host, value[i], timeout)
				if opened {
					tmpPorts = append(tmpPorts, value[i])
					//fmt.Printf("任务%v，执行时长:%.2fs，【端口开放】:%v \n", key, time.Since(start).Seconds(), value[i])
				}else{
					//fmt.Printf("任务%v，执行时长:%.2fs，【端口关闭】:%v \n", key, time.Since(start).Seconds(), value[i])

				}
			}
			mutex.Lock()
			ports = append(ports, tmpPorts...)
			mutex.Unlock()
			if len(tmpPorts) > 0 {
				//fmt.Printf("协程%v 执行完成，总执行时长： %.2fs，开放端口： %v \n", key, time.Since(start).Seconds(), tmpPorts)
			}
		}(v, k)
	}
	wg.Wait()
	fmt.Printf("%v【%v】总执行时长%.2fs , 开放的端口:%v\n", time.Now().Format("2006-01-02 15:04:05"), host, time.Since(start).Seconds(), ports)
	return ports
}

//获取所有ip
func (s *ScanIp) getAllIp(ip string) ([]string, error) {
	var hosts []string
	ips := strings.Split(ip, "-")
	firstIp, err := net.ResolveIPAddr("ip", ips[0])

	if err != nil {
		return hosts, errors.New(ips[0] + " 解析失败" + err.Error())
	}
	if net.ParseIP(firstIp.String()) == nil {
		return hosts, errors.New(ips[0] + " ip地址有误~")
	}

	ips[0] = firstIp.String()
	if len(ips) == 2 {
		ipArr := strings.Split(ips[0], ".")
		startIp, _ := strconv.Atoi(ipArr[3])
		endIp, err := strconv.Atoi(ips[1])
		if err != nil || endIp < startIp {
			endIp = startIp
		}
		if endIp > 255 {
			endIp = 255
		}
		totalIp := endIp - startIp + 1
		for i := 0; i < totalIp; i++ {
			hosts = append(hosts, fmt.Sprintf("%s.%s.%s.%d", ipArr[0], ipArr[1], ipArr[2], startIp+i))
		}
	} else {
		hosts = append(hosts, ips[0])

	}
	return hosts, nil
	//return []string{"10.70.66.214","10.70.66.215"}, nil
}

//查看端口号是否打开
func (s *ScanIp) isOpen(host string, port int, timeout int) bool {
	time.Sleep(time.Microsecond * 100)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Millisecond*time.Duration(timeout))
	if err == nil {
		_ = conn.Close()
		return true
	}else{
		return false
	}

}
