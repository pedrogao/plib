package tcmalloc

import (
	"log"
)

const (
	cacheGCLimit = 1 << 20 // 1M
)

type cache struct {
	totalSize   int
	freeObjects [numClasses]*objectList
	central     *central
}

func newCache(central *central) *cache {
	ch := &cache{
		totalSize:   0,
		freeObjects: [numClasses]*objectList{},
		central:     central,
	}

	for i := 0; i < numClasses; i++ {
		ch.freeObjects[i] = newObjectList()
	}

	return ch
}

func (c *cache) fill(class int) error {
	freeObject := c.freeObjects[class]
	count, err := c.central.getObjects(class, batchSize, freeObject) // 从 central 中获取新的空对象
	if err != nil {
		return err
	}
	c.totalSize += class2Size(class) * count
	return nil
}

func (c *cache) getObject(class int) (*object, error) {
	freeObject := c.freeObjects[class]
	// 空闲列表为空，则填充列表
	if freeObject.isEmpty() {
		if err := c.fill(class); err != nil {
			return nil, err
		}
	}
	// 更新当前缓存中的可分配内存大小
	c.totalSize -= class2Size(class)
	return freeObject.pop(), nil
}

func (c *cache) insert(class int, object *object) {
	freeObject := c.freeObjects[class]
	freeObject.push(object)
	c.totalSize += class2Size(class)
	// 清除缓存
	if c.totalSize > cacheGCLimit {
		c.gc()
	}
}

func (c *cache) gc() {
	for class := 0; class < numClasses; class++ {
		freeObject := c.freeObjects[class]
		size := class2Size(class)
		for i := 0; i < freeObject.lowMark/2; i++ {
			o := freeObject.pop()
			c.totalSize -= size
			err := c.central.insertObject(class, o)
			if err != nil {
				log.Printf("gc err: %s", err)
			}
		}

		freeObject.lowMark = freeObject.length
	}
}

func (c *cache) size() int {
	totalSize := 0
	for i := 0; i < numClasses; i++ {
		freeObject := c.freeObjects[i]
		totalSize += freeObject.length * class2Size(i)
	}
	return totalSize
}
