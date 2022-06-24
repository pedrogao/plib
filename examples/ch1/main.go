package main

import (
	"fmt"
	"time"
)

/*
	@jobs
	@result
	refer https://www.cnblogs.com/bigdataZJ/p/go-channel-deadlock.html
*/
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func main() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 3; w++ {
		// consume job, and send to result
		go worker(w, jobs, results)
	}

	for j := 1; j <= 3; j++ {
		jobs <- j
	}
	close(jobs) // jobs can't receive

	for i := 0; i < 3; i++ { // 如果只消费三次，那么就不会死锁，因为就只有3个数据
		fmt.Println(<-results)
	}

	// for a := range results { // 输出完毕后，无法再接受，但也没人发送，所以会触发死锁
	//   fmt.Println(a)
	// }
}
