package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main()  {
	http.HandleFunc("/list/", func(writer http.ResponseWriter, request *http.Request) {
		path:=request.URL.Path[len("/list"):]
		fmt.Println(path)
		file,err:=os.Open(path)
		if err!=nil{
			code:=http.StatusOK
			switch {
			case os.IsNotExist(err):
				code = http.StatusNotFound
			case os.IsPermission(err):
				code = http.StatusForbidden
			default:
				code =http.StatusInternalServerError
			}
			http.Error(writer,http.StatusText(code),http.StatusOK)
			return
		}
		defer file.Close()
		all,err:=ioutil.ReadAll(file)
		if err!=nil{
			panic(err)
		}
		writer.Write(all)
	})

	http.ListenAndServe(":8888",nil)
}
