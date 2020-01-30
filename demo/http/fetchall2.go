package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {
	var urls = []string{"http://www.baidu.com",
		"https://www.qq.com",
		"https://lf.58.com",
		"http://www.0532888.cn",
		"https://www.xin.com",
	}
	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(len(urls))
	for _, val := range urls {
		go func(url string) {
			start := time.Now()
			res, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				return
			}
			nbytes, err := io.Copy(ioutil.Discard, res.Body)
			if err != nil {
				fmt.Printf("while reading %s:%v\n", url, err)
				return
			}
			fmt.Printf("%.2fs %7d %s\n", time.Since(start).Seconds(), nbytes, url)
			wg.Done()
		}(val)
	}

	wg.Wait()

	fmt.Printf("%.2fs elapsed \n", time.Since(start).Seconds())

}
