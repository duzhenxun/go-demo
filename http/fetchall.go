package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func main()  {
	start :=time.Now()
	ch :=make(chan string)
	urls,_:=Extract("http://www.0532888.cn")
	/*var urls = []string{"http://www.baidu.com",
		"http://www.qq.com",
		"http://www.58.com",
		"http://www.0532888.cn",
		"https://www.xin.com",
	}*/
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


func Extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

//!-Extract

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}