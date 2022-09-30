package syncx

// concurrent linked set
// see: https://docs.google.com/presentation/d/1NMg08N1LUNDPuMxNZ-UMbdH13p8LXgMM3esbWRMowhU/edit#slide=id.ga50d357ce9_0_209
// @author pedrogao

import (
	"sync"
)

type (
	node struct {
		data Ordered
		// Each node has a lock
		// When accessing the “next” pointer, lock should be held
		// A node’s lock protects the node’s next pointer
		// Acquire the next node’s lock before releasing the current’s (“hand-over-hand”)
		mu   sync.Mutex
		next *node
	}

	// OrderedListSet concurrent ordered linked set
	OrderedListSet struct {
		head *node // node不能为空，data可以为空
	}
)

func NewOrderedListSet() *OrderedListSet {
	return &OrderedListSet{head: newNode(nil, nil)}
}

func newNode(data Ordered, next *node) *node {
	return &node{
		data: data,
		mu:   sync.Mutex{},
		next: next,
	}
}

func (l *OrderedListSet) Find(key Ordered) (Ordered, bool) {
	cur := l.head
	if cur.data == nil {
		return nil, false
	}
	for !cur.data.Equal(key) {
		// if cur > key，证明后面的更大，因此不可能找到
		if !cur.data.Less(key) {
			return nil, false
		}
		// 当访问 next 时，需要加锁
		cur.mu.Lock()
		next := cur.next
		cur.mu.Unlock()
		// 解锁完后再替换
		cur = next
	}
	return cur.data, true
}

func (l *OrderedListSet) Contains(key Ordered) bool {
	_, ok := l.Find(key)
	return ok
}

// Insert B between A and C
// 1. Acquire A’s lock
// 2. Read A.next (=C)
// 3. Allocate B with B.next=C
// 4. Write A.next = B
func (l *OrderedListSet) Insert(key Ordered) bool {
	var (
		pre *node
		cur = l.head
	)

	if cur.data == nil {
		cur.mu.Lock()
		cur.data = key
		cur.mu.Unlock()
		return true
	}

	for !cur.data.Equal(key) {
		// 找到 cur.key < key < next.key 的地方插入
		// 如果 cur.key > key，那么将 key 插入到 cur 的前面
		// 注意：是前面
		if !cur.data.Less(key) {
			pre.mu.Lock()
			next := pre.next
			n := newNode(key, next)
			pre.next = n
			pre.mu.Unlock()
			return true
		}
		// 当访问 next 时，需要加锁
		cur.mu.Lock()
		next := cur.next
		cur.mu.Unlock()
		// 解锁完后再替换
		pre = cur
		cur = next
	}
	// 相等返回false
	return false
}

// Remove B between A and C
// 1. Acquire A’s lock
// 2. Read A.next (=B)
// 3. Acquire B’s lock
// 4. Read B.next (=C)
// 5. Write A.next=C
// 6. Release B’s lock and free B
func (l *OrderedListSet) Remove(key Ordered) Ordered {
	var (
		pre *node
		cur = l.head
	)
	// cur.key > key 时，后面的只会越来越大，因此无法再找到了
	for !cur.data.Less(key) {
		// 相等时，remove
		if cur.data.Equal(key) {
			pre.mu.Lock()
			next := pre.next
			next.mu.Lock()
			nnext := next.next
			pre.next = nnext
			pre.mu.Unlock()
			next.mu.Unlock()
			next = nil // for gc
			return key
		}
		// 当访问 next 时，需要加锁
		cur.mu.Lock()
		next := cur.next
		cur.mu.Unlock()
		// 解锁完后再替换
		pre = cur
		cur = next
	}

	return nil
}
