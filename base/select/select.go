package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	var t1 = makeTask("adoJob", 1000)
	var t2 = makeTask("xs25Job", 1000)
	var allTask []string                  //因为我想只做一个消费端，将2个生产端生产出来的消费都扔到一起
	var tick = time.Tick(time.Second * 5) //每隔一段时间报告队列积压情况
	var workerCh = worker()

	for {
		var taskInfo string //具体任务
		var ch chan<- string
		if len(allTask) > 0 {
			taskInfo = allTask[0]
			ch = workerCh
		}
		select {
		case task := <-t1:
			allTask = append(allTask, task)
		case task := <-t2:
			allTask = append(allTask, task)
		case ch <- taskInfo: //任务详情写入到要消费工作中
			allTask = allTask[1:]
		case <-tick:
			log.Println("队列挤压数量", len(allTask))
		}
	}
}

//生产数据
func makeTask(queueName string, n int) chan string {
	ch := make(chan string, n)
	go func() {
		i := 1
		for {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000))) //假设生产任务占用时间
			ch <- fmt.Sprintf("%s,生产数据 %d", queueName, i)
			i++
		}
	}()
	return ch
}

//消费数据
func worker() chan<- string {
	ch := make(chan string)
	go func(tasks chan string) {
		for t := range tasks {
			time.Sleep(time.Second * 1) //假设我们每次消费任务需要花费1秒钟
			log.Printf("消费任务: %s \n", t)
		}
	}(ch)
	return ch
}
