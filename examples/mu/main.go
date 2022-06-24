package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var total []int
	var mu sync.Mutex

	arr := []int{10, 30, 29, 39, 14, 34, 23, 59}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		// 安全模式下执行函数
		go func(s int) {
			fmt.Printf("received: %d\n", s)
			defer func() {
				if err := recover(); err != nil {
					log.Printf("%s", err)
				}
			}()
			defer wg.Done()

			for j := 0; j < s; j++ {
				mu.Lock()
				total = append(total, j)
				mu.Unlock()
			}
		}(arr[i])
	}

	wg.Wait()

	fmt.Printf("%v\n", total)

	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	fmt.Printf("%v\n", total)
}
