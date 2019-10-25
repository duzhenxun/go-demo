package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func main()  {
	start :=time.Now()
	ch :=make(chan string)

	var urls = []string{"http://www.baidu.com",
		"http://www.qq.com",
		"http://www.58.com",
		"http://www.0532888.cn",
		"https://www.xin.com",
	}
	for _,url:=range urls[:len(urls)]{
		go fetch(url,ch)
	}

	for range urls[:len(urls)]{
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed \n",time.Since(start).Seconds())

}

func fetch(url string,ch chan<- string){
	start:=time.Now()
	res,err:=http.Get(url)
	if err!=nil{
		ch <- fmt.Sprint(err)
		return
	}
	nbytes,err:=io.Copy(ioutil.Discard,res.Body)
	if err!=nil{
		ch <- fmt.Sprintf("while reading %s:%v",url,err)
		return
	}
	secs:=time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s",secs,nbytes,url)

}