package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	
	fmt.Println(r.URL.Scheme)
	fmt.Println(r.URL.Path)
	fmt.Println(r.Form["url_long"])

	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprint(w, "远程办公视频会议")
}

func login(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.Header)
	if r.Method == "GET" {
		timeStart := time.Now()
		w.Write([]byte("hello"))
		t, _ := template.ParseFiles("login.html")
		time.Sleep(1e9 * 2)
		timeElapsed := time.Since(timeStart) / 1e9 //执行时间秒
		t.Execute(w, nil)
		fmt.Println(timeElapsed)
	} else {
		r.ParseForm()
		fmt.Println("method:", r.Method) //获取请求的方法
		if m, _ := regexp.MatchString("^[0-9]+$", r.Form.Get("username")); !m {
			fmt.Println("username不是数字:", strings.Join(r.Form["username"], ""))
		}
		//请求的是登陆数据，那么执行登陆的逻辑判断
		fmt.Println("username:", strings.Join(r.Form["username"], ""))
		fmt.Println("password:", r.Form["password"])
	}
}

func layui(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("layui.html")
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/login", login)

	err := http.ListenAndServe(":19090", nil)
	if err != nil {
		log.Fatal("listenAndServe:", err)
	}
}
