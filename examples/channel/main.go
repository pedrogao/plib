package main

import (
	"fmt"
	"time"
)

func main() {
	pipe := make(chan int, 10)
	arr := []int{10, 30, 20, 40, 10, 30, 30, 60} // 230

	for i := 0; i < 8; i++ {
		// 安全模式下执行函数
		go func(s int) {
			fmt.Printf("received: %d\n", s)

			for j := 0; j < s; j++ {
				pipe <- j
			}
		}(arr[i])
	}

	var total []int
	b := false
	for {
		select {
		case i := <-pipe:
			total = append(total, i)
			// 判断是否完成
			if len(total) >= 230 {
				b = true
				break
			}
		case <-time.After(time.Second * 10):
			b = true
			break
		}
		if b {
			break
		}
	}

	// close(pipe) // 无法在此关闭，因为无法判断数据是否发送完毕

	fmt.Printf("%d\n", len(total))
}
