package main

import "sync"

var count int

func main() {
	var wg sync.WaitGroup
	var l sync.Mutex
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Lock()
			count++
			l.Unlock()
		}()
	}
	wg.Wait()
	println(count)
}
