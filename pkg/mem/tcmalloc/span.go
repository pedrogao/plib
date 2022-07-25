package tcmalloc

import (
	"fmt"
	"log"
	"unsafe"
)

type spanState int

const (
	free spanState = iota + 1
	allocatedSmall
	allocatedLarge
)

type span struct {
	next *span
	prev *span

	pageStart uintptr
	numPages  int
	state     spanState

	objectClass int // 每个 span 都对应一个 object 大小
	freeObjects *objectList
}

func newSpanWith(pageStart uintptr, numPages int) *span {
	//addr, err := mmapAnonymous(int(unsafe.Sizeof(span{})))
	//if err != nil {
	//	panic(err)
	//}
	//
	//s := (*span)(unsafe.Pointer(addr))
	//s.pageStart = pageStart
	//s.numPages = numPages
	//s.state = free
	//s.objectClass = 0
	//s.freeObjects = newObjectList()

	s := &span{
		pageStart:   pageStart,
		numPages:    numPages,
		state:       free,
		objectClass: 0,
		freeObjects: newObjectList(),
	}
	return s
}

func newSpan(numPages int) (*span, error) {
	pageStart, err := mmapAnonymous(numPages * pageSize)
	if err != nil {
		return nil, fmt.Errorf("new span err: %s", err)
	}
	return newSpanWith(pageStart, numPages), nil
}

func (s *span) isConsistent() bool {
	// [1, Max]
	ret := 1 <= s.numPages && s.numPages <= maxNumPages

	if s.state == allocatedSmall {
		ret = ret && s.objectClass < numClasses
	}

	return ret
}

// 拆分 span
func (s *span) split(numPages int, pHeap *pageHeap) *span {
	// 先移除 s
	pHeap.getPageMap().removeSpan(s)
	// 减去 numPages 后
	s.numPages -= numPages
	// 再插入 span
	err := pHeap.getPageMap().insertSpan(s)
	if err != nil {
		log.Printf("insertAt span err: %s", err)
	}
	// 将新的 span 插入 map
	ns := newSpanWith(s.pageStart+uintptr(s.numPages*pageSize), numPages)
	err = pHeap.getPageMap().insertSpan(ns)
	if err != nil {
		log.Printf("insertAt span err: %s", err)
	}
	return ns
}

func (s *span) makeObjects(objectClass int) {
	if s == nil {
		log.Printf("try to make objects: %d, but span is nil", objectClass)
		return
	}
	s.state = allocatedSmall
	size := class2Size(objectClass)
	// +size 是为了一个 object 不越界
	for obj := s.pageStart; obj+uintptr(size) <= s.pageStart+uintptr(s.numPages*pageSize); obj = obj + uintptr(size) {
		o := (*object)(unsafe.Pointer(obj))
		s.freeObjects.push(o)
	}

	s.objectClass = objectClass
}

// 合并 span 到 s
// s + other => s
// other + s => s
func (s *span) mergeWith(other *span, isBefore bool) {
	// 如果 other 地址在前面，则更新 pageStart
	if isBefore {
		s.pageStart = other.pageStart
	}
	// 更新 numPages
	s.numPages += other.numPages
}

// coalesce 将 span 与前、后的 span 合并
func (s *span) coalesce(pHeap *pageHeap) {
	pHeap.getPageMap().removeSpan(s)

	pageNum := int(s.pageStart / pageSize)
	pageNumBefore := pageNum - 1
	pageNumAfter := pageNum + s.numPages

	before := pHeap.getPageMap().get(pageNumBefore)
	if before != 0 {
		beforeSpan := (*span)(unsafe.Pointer(before))
		s.mergeWith(beforeSpan, true)
		// beforeSpan 被合并了，从 map 和 list 中移除
		pHeap.getPageMap().removeSpan(beforeSpan)
		pHeap.getSpanList(beforeSpan).remove(beforeSpan)
	}

	// span 只有一页
	if pageNumBefore == pageNumAfter {
		return
	}

	after := pHeap.getPageMap().get(pageNumAfter)
	if after != 0 {
		afterSpan := (*span)(unsafe.Pointer(after))
		s.mergeWith(afterSpan, false)
		// afterSpan 被合并了，从 map 和 list 中移除
		pHeap.getPageMap().removeSpan(afterSpan)
		pHeap.getSpanList(afterSpan).remove(afterSpan)
	}
}
