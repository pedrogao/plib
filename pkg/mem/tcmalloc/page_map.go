package tcmalloc

import (
	"sync"
	"unsafe"
)

const (
	radixTreeNumBits  = 18
	radixTreeNodeSize = 1 << 18               // 262144
	radixTreeMask     = radixTreeNodeSize - 1 // 262143
)

type (
	radixTreeNode struct {
		children [radixTreeNodeSize]uintptr
	}

	pageMap struct {
		root *radixTreeNode
		mu   sync.Mutex
	}
)

var radixTreeNodeMemSize = int(unsafe.Sizeof(radixTreeNode{}))

func newRadixTreeNode() *radixTreeNode {
	addr, err := mmapAnonymous(radixTreeNodeMemSize)
	if err != nil {
		panic(err)
	}
	node := (*radixTreeNode)(unsafe.Pointer(addr))
	return node
}

func newPageMap() *pageMap {
	return &pageMap{root: newRadixTreeNode()}
}

// 页开始序号
// page_number_start = span->page_start / PAGE_SIZE;
// 页结束序号
// page_number_end = page_number_start + span->num_pages - 1;
func (m *pageMap) insertAt(pageNum int, addr uintptr) error {
	i1 := getIndex(pageNum, 1)
	l2 := m.root.children[i1]
	if l2 == 0 {
		node := newRadixTreeNode()
		l2 = uintptr(unsafe.Pointer(node))
		m.root.children[i1] = l2
	}

	i2 := getIndex(pageNum, 2)
	ll2 := (*radixTreeNode)(unsafe.Pointer(l2))
	l3 := ll2.children[i2]
	if l3 == 0 {
		node := newRadixTreeNode()
		l3 = uintptr(unsafe.Pointer(node))
		ll2.children[i2] = l3
	}

	i3 := getIndex(pageNum, 3)
	ll3 := (*radixTreeNode)(unsafe.Pointer(l3))
	ll3.children[i3] = addr

	return nil
}

func (m *pageMap) insertSpan(s *span) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	starPageNum := s.pageStart / pageSize
	endPageNum := starPageNum + uintptr(s.numPages) - 1
	p := uintptr(unsafe.Pointer(s))
	err := m.insertAt(int(starPageNum), p)
	if err != nil {
		return err
	}
	err = m.insertAt(int(endPageNum), p)
	if err != nil {
		return err
	}
	return nil
}

func (m *pageMap) removeSpan(s *span) {
	m.mu.Lock()
	defer m.mu.Unlock()

	starPageNum := s.pageStart / pageSize
	endPageNum := starPageNum + uintptr(s.numPages) - 1
	m.removeAt(int(starPageNum))
	m.removeAt(int(endPageNum))
}

func (m *pageMap) removeAt(pageNum int) {
	i1 := getIndex(pageNum, 1)
	l2 := m.root.children[i1]
	if l2 == 0 {
		return
	}

	i2 := getIndex(pageNum, 2)
	ll2 := (*radixTreeNode)(unsafe.Pointer(l2))
	l3 := ll2.children[i2]
	if l3 == 0 {
		return
	}

	i3 := getIndex(pageNum, 3)
	ll3 := (*radixTreeNode)(unsafe.Pointer(l3))
	ll3.children[i3] = 0
}

func (m *pageMap) get(pageNum int) uintptr {
	m.mu.Lock()
	defer m.mu.Unlock()

	i1 := getIndex(pageNum, 1)
	l2 := m.root.children[i1]
	if l2 == 0 {
		return 0
	}

	i2 := getIndex(pageNum, 2)
	ll2 := (*radixTreeNode)(unsafe.Pointer(l2))
	l3 := ll2.children[i2]
	if l3 == 0 {
		return 0
	}

	i3 := getIndex(pageNum, 3)
	ll3 := (*radixTreeNode)(unsafe.Pointer(l3))
	return ll3.children[i3]
}

func getIndex(pageNum, level int) int {
	return (pageNum >> (radixTreeNumBits * (3 - level))) & radixTreeMask
}
