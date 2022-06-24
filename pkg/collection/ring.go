package collection

import "sync"

// Ring 环形切片
type Ring[T any] struct {
	elements []T
	index    int
	lock     sync.Mutex
}

// NewRing create a ring
func NewRing[T any](n int) *Ring[T] {
	if n < 1 {
		panic("n should be greater than 0")
	}
	return &Ring[T]{
		elements: make([]T, n),
	}
}

// NewRingWithElements create ring with elements
func NewRingWithElements[T any](t []T) *Ring[T] {
	return &Ring[T]{
		elements: t,
		index:    len(t),
	}
}

// Add 添加元素
func (r *Ring[T]) Add(v T) {
	r.lock.Lock()
	defer r.lock.Unlock()
	// 将元素写入切片指定位置
	// 这里的取余实现了循环写效果
	r.elements[r.index%len(r.elements)] = v
	// 更新下次写入位置
	r.index++
}

// Take 获取全部元素
// 读取顺序保持与写入顺序一致
func (r *Ring[T]) Take() []T {
	r.lock.Lock()
	defer r.lock.Unlock()

	var size int
	var start int
	// 当出现循环写的情况时
	// 开始读取位置需要通过去余实现,因为我们希望读取出来的顺序与写入顺序一致
	if r.index > len(r.elements) {
		size = len(r.elements)
		// 因为出现循环写情况,当前写入位置index开始为最旧的数据
		start = r.index % len(r.elements)
	} else {
		size = r.index
	}
	elements := make([]T, size)
	for i := 0; i < size; i++ {
		// 取余实现环形读取,读取顺序保持与写入顺序一致
		elements[i] = r.elements[(start+i)%len(r.elements)]
	}

	return elements
}
