package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	PrintMemUsage()

	pid := os.Getpid()
	fmt.Printf("process: %d\n", pid)
	time.Sleep(time.Second * 1)

	var overall [][]byte
	for i := 0; i < 4; i++ {
		a := make([]byte, 0, 1024*1024*50)
		overall = append(overall, a)

		PrintMemUsage()
		time.Sleep(time.Second)
	}

	overall = nil
	PrintMemUsage()

	runtime.GC()
	PrintMemUsage()
}

// PrintMemUsage print mem usage
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
