package main

import "runtime"

func main() {
	big := make([]byte, 1024*1024*200)
	runtime.KeepAlive(big)
}
