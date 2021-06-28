package main

import (
	"fmt"
	"os"
)

func main()  {
	fmt.Println("请输入一个目录的路径:")
	var path string
	if n,err:=fmt.Scan(&path);err!=nil{
		fmt.Println(err.Error(),n)
	}
	openDir(path,0)

}

func openDir(path string,i int)  {
	i++
	f,err:=os.OpenFile(path,os.O_RDONLY,os.ModeDir)

	if err !=nil{
		fmt.Printf("Open file failed:%s.\n",err)
		return
	}
	defer f.Close()

	info,err:=f.Readdir(-1) //-1 读取目录中的所有目录
	for _,fileInfo:=range info{
		if fileInfo.IsDir(){
			fmt.Printf("%v--%s是一个目录\n",i,fileInfo.Name())
			openDir(path+"/"+fileInfo.Name(),i)
		}else{
			//fmt.Printf("%v--%s是一个文件\n",i,fileInfo.Name())
		}
	}
}