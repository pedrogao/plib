package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/pedrogao/plib/pkg/mem"
	"github.com/struCoder/pidusage"
)

func main() {
	PrintMemUsage()

	pid := os.Getpid()
	fmt.Printf("process: %d\n", pid)
	time.Sleep(time.Second * 5)

	var overall []*mem.Area
	// 1216 KB
	for i := 0; i < 10; i++ {
		area, err := mem.Init(1024 * 1024 * 50)
		if err != nil {
			log.Fatalln(err)
		}

		overall = append(overall, area)

		PrintMemUsage()
		time.Sleep(time.Second)
	}

	printStat()

	time.Sleep(time.Second * 5)

	for i := 0; i < 10; i++ {
		err := overall[i].Release()
		if err != nil {
			log.Fatalln(err)
		}
		PrintMemUsage()
	}

	time.Sleep(time.Second * 5)

	runtime.GC()
	PrintMemUsage()
}

func printStat() {
	sysInfo, err := pidusage.GetStat(os.Getpid())
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("CPU = %v ", sysInfo.CPU)
	fmt.Printf("Mem = %v \n", sysInfo.CPU)
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
