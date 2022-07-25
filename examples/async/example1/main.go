package main

import (
	"fmt"
	"time"
)

func DoneAsync() chan int {
	r := make(chan int)
	fmt.Println("Warming up...")

	go func() {
		time.Sleep(3 * time.Second) // 模拟任务执行
		r <- 1
		fmt.Println("Worker done")
	}()

	return r
}

func main() {
	fmt.Println("Start...")
	val := DoneAsync()
	fmt.Println("Worker is running...")
	fmt.Println(<-val) // 阻塞
}
