package tcmalloc

import (
	"fmt"
	"log"
	"sync"
	"unsafe"
)

const (
	batchSize = 32 // 每次申请的个数，在 tcmalloc 中，每个 class 对应了不同的个数，这里简单处理为 32
)

type central struct {
	spans    [numClasses]*spanList
	pageHeap *pageHeap
	mu       sync.Mutex
}

func newCentral(pageHeap *pageHeap) *central {
	freeList := &central{
		spans:    [numClasses]*spanList{},
		pageHeap: pageHeap,
	}
	for i := 0; i < numClasses; i++ {
		freeList.spans[i] = newList()
	}

	return freeList
}

func (c *central) fill(class int) error {
	spanL := c.spans[class]
	numPages := class2Size(class) * batchSize / pageSize
	if numPages < 1 {
		numPages = 1 // 最小也得分配一页
	}
	freeSpan, err := c.pageHeap.getSpan(numPages) // 从 pageHeap 中获取空闲的 span
	if err != nil {
		return err
	}
	freeSpan.makeObjects(class)
	spanL.push(freeSpan)
	return nil
}

func (c *central) getObjects(class, batchSize int,
	list *objectList) (int, error) {
	c.mu.Lock()
	log.Printf("getObjects, class: %d, size: %d, list: %p\n", class, batchSize, list)
	spanL := c.spans[class]
	count := 0

	for iter := spanL.begin(); iter != nil; iter = iter.next {
		// 注意：一次请求 batchSize 个，但列表很有可能被分配完了
		// 因此实际的分配个数可能小于 batchSize 个
		for !iter.freeObjects.isEmpty() && count < batchSize {
			obj := iter.freeObjects.pop() // pop
			list.push(obj)
			count++
		}
	}
	// central 中也没有对象，需要分配
	if count == 0 {
		err := c.fill(class)
		if err != nil {
			return 0, err
		}
		c.mu.Unlock()
		return c.getObjects(class, batchSize, list)
	}

	c.mu.Unlock()
	return count, nil
}

func (c *central) insertObject(class int, object *object) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	spanL := c.spans[class]
	op := uintptr(unsafe.Pointer(object))
	var iter *span
	// 通过地址找到 object 对应的 list
	for iter = spanL.begin(); iter != nil; iter = iter.next {
		if iter.pageStart <= op &&
			op < iter.pageStart+uintptr(iter.numPages*pageSize) {
			break
		}
	}
	if iter == nil {
		return fmt.Errorf("can't find span")
	}
	// 将 object 回收进 list
	iter.freeObjects.push(object)

	// 如果此时 list 的对象总大小(length * size) 将要超过 pageNums * pageSize
	// 即 list 所容纳的大小超过了页数，则需将前、后页合并
	if (iter.freeObjects.length+1)*class2Size(iter.objectClass) >
		iter.numPages*pageSize {
		// 将 span 从 spanList 中移除
		spanL.remove(iter)
		// span 现在是游离态，合并后再插入进 list
		iter.coalesce(c.pageHeap)
		err := c.pageHeap.insertSpan(iter)
		if err != nil {
			return err
		}
	}
	return nil
}
