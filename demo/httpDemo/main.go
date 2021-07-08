package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/log", func(writer http.ResponseWriter, request *http.Request) {
		fileName := "http.log"
		f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
		defer f.Close()
		if err != nil {
			panic(err)
		}
		txt, _ := ioutil.ReadAll(request.Body)
		fmt.Println(string(txt))
		n, _ := f.Seek(0, 2)
		t := time.Now().Format("2006-01-02 15:04:05")
		f.WriteAt([]byte(t+"\n"+string(txt)+"\n"), n)

		writer.Write([]byte("ok"))
	})
	http.ListenAndServe(":5001", nil)
}