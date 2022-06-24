package tcmalloc

import (
	"fmt"
	"syscall"
)

/*
 * Size classes:
 * Object大小列，一次申请的page数，一次移动的objects数(内存申请或回收)
 * 8, 16, ..., 64, (8)
 * 128, 192, ..., 2048, (31)
 * 2304, 2560, ..., 32768 (120)
 */

const (
	pageSize        = 4096
	maxNumPages     = 1 << 12
	smallObjectSize = 32768
	numClasses      = 159
)

func size2Class(size int) int {
	if size <= 64 {
		return (size - 1) / 8
	} else if size <= 2048 {
		return 8 + (size-64-1)/64
	} else {
		return 8 + 31 + (size-2048-1)/256
	}
}

func class2Size(class int) int {
	if class < 8 {
		return 8 + class*8
	} else if class < 8+31 {
		return 128 + (class-8)*64
	} else {
		return 2304 + (class-8-31)*256
	}
}

func mmapAnonymous(size int) (uintptr, error) {
	// MAP_ANON：会忽略参数 fd，映射区不与任何文件关联，而且映射区域无法和其他进程共享
	prot, flags, fd, offset := syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE, -1, 0
	addr, _, err := syscall.Syscall6(syscall.SYS_MMAP, uintptr(0), uintptr(size),
		uintptr(prot), uintptr(flags), uintptr(fd), uintptr(offset))
	if err != 0 {
		return 0, fmt.Errorf("mmap err: %s", err.Error())
	}
	return addr, nil
}
