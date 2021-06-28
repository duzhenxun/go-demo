package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	t := time.Now()
	s:=strconv.FormatInt(t.UnixNano()/1e6,10)
	str := t.Format("2006-01-02 03:04:05")+"."+s[10:]
	fmt.Println(str)

}
