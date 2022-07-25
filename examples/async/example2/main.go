package main

import (
	"fmt"
	"time"

	"github.com/pedrogao/plib/pkg/async"
)

func DoneAsync() int {
	fmt.Println("Warming up...")
	time.Sleep(1 * time.Second) // 模拟任务执行
	fmt.Println("Worker done")
	return 1
}

func main() {
	fmt.Println("Start...")

	future := async.Exec(func() any {
		return DoneAsync()
	})

	fmt.Println("Worker is running...")
	val := future.Await() // 阻塞
	fmt.Println(val)
}
