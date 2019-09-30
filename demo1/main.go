package main

import (
	"fmt"
	"sync"
	"time"
)

func say(s string) {
	time.Sleep(1000 * time.Microsecond)
	fmt.Println(s)
}

//a.com/a
//a.com/b

func main() {

	s := []int{7, 2, 8, -9, 4, 0}
	fmt.Println(s[:len(s)/2])
	fmt.Println(s[len(s)/2:])

	//c := make(chan string)
	var  tmp []map[string]string
	tmp = append(tmp, map[string]string{"name":"duzhenxun"},map[string]string{"name":"wang"},map[string]string{"name":"liu"})
	wg := sync.WaitGroup{}
	wg.Add(len(tmp))
	for _,v:=range tmp{
		go func(v map[string]string) {
			say(v["name"])
			wg.Done()
		}(v)
	}
	wg.Wait()


	fmt.Println(".......")

}
