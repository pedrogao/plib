package tcmalloc

import (
	"log"
	"sync"
	"unsafe"
)

/*
 page_map 用来保存 class 与 span 之间的关系：当 malloc 时，计算出 class 然后拿到对应的 span
 page_heap 用来保存全部的 span、pages
*/

var (
	globalPageHeap *pageHeap
	global         *central
	globalPageMap  *pageMap
	threadCache    *cache
	globalLock     sync.Mutex
)

func init() {
	globalPageMap = newPageMap()
	globalPageHeap = newPageHeap(globalPageMap)
	global = newCentral(globalPageHeap)
	threadCache = newCache(global)
}

// Malloc 分配内存
func Malloc(size int) uintptr {
	if size <= 0 {
		return 0
	}
	// 小对象
	if size <= smallObjectSize {
		// 注意：原则上对于线程缓存是不需要加锁的，但是 go 没有提供 _thread cache
		// 之类的API，为了保证线程安全，因此加上全局锁，避免数据竞争
		// globalLock.Lock()
		// defer globalLock.Unlock()
		// 计算 size 对应的 class
		class := size2Class(size)
		obj, err := threadCache.getObject(class)
		if err != nil {
			log.Printf("malloc err: %s", err)
			return 0
		}
		if obj == nil {
			return 0
		}
		return uintptr(unsafe.Pointer(obj))
	}
	// 大对象
	numPages := (size-1)/pageSize + 1
	globalLock.Lock()
	defer globalLock.Unlock()

	s, err := globalPageHeap.getSpan(numPages)
	if err != nil {
		log.Printf("malloc err: %s", err)
		return 0
	}
	s.state = allocatedLarge
	return s.pageStart
}

// Free 释放内存
func Free(p uintptr) {
	if p == 0 {
		return
	}
	// 先获取对应的 span
	pageNum := p / pageSize
	addr := globalPageMap.get(int(pageNum))
	if addr == 0 {
		return
	}
	s := (*span)(unsafe.Pointer(addr))
	// 无需再次释放
	if s.state == free {
		return
	}
	// 小对象，回收至缓存
	if s.state == allocatedSmall {
		threadCache.insert(s.objectClass, (*object)(unsafe.Pointer(p)))
		return
	}
	// 大对象直接放回 pageHeap
	globalLock.Lock()
	defer globalLock.Unlock()

	s.coalesce(globalPageHeap)
	err := globalPageHeap.insertSpan(s)
	if err != nil {
		log.Printf("free err: %s", err)
	}
}
