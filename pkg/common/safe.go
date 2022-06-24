package common

import "log"

// Recover 恢复，带清理函数
func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups { // resource cleanup
		cleanup()
	}
	if r := recover(); r != nil {
		log.Printf("%v", r)
	}
}

// RunSafe 安全运行 goroutine
func RunSafe(fn func()) {
	defer Recover()
	fn()
}

// GoSafe 异步运行函数
func GoSafe(fn func()) {
	go RunSafe(fn)
}
