package main

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	pool := sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}

	processRequest := func(size int) {
		b := pool.Get().(*bytes.Buffer)    // Get 从队列里面拿，实际上是随机的
		time.Sleep(500 * time.Millisecond) // Simulate processing time
		b.Grow(size)                       // 扩容后再put

		fmt.Printf("buffer size: %d, cap: %d\n", b.Len(), b.Cap())
		fmt.Printf("address of buffer %p\n", b)

		pool.Put(b)
		time.Sleep(1 * time.Millisecond) // Simulate idle time
	}

	// Simulate a set of initial large writes.
	for i := 0; i < 10; i++ {
		go func() {
			processRequest(1 << 28) // 256MiB
		}()
	}

	time.Sleep(time.Second) // Let the initial set finish

	// Simulate an un-ending series of small writes.
	for i := 0; i < 10; i++ {
		go func() {
			for {
				processRequest(1 << 10) // 1KiB
			}
		}()
	}

	// Continually run a GC and track the allocated bytes.
	// 40 Cycles 才被完全 free
	var stats runtime.MemStats
	for i := 0; ; i++ {
		runtime.ReadMemStats(&stats)
		fmt.Printf("Cycle %d: %dMB\n", i, stats.Alloc/(1024*1024))
		time.Sleep(time.Second)
		runtime.GC()
	}
}
