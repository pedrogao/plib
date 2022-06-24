package tcmalloc

import (
	"fmt"
)

const (
	pageHeapArraySize = 256
)

// pageHeap 管理 span list
// `pageHeap.freeSpans[i]` contains spans with `numPages == i + 1`.
// Except the last element `pageHeap.freeSpans[255]`, which contains spans
// with `numPages >= 256`.
type pageHeap struct {
	// spanList 按照 numPages(页个数) 来组织
	freeSpans [pageHeapArraySize]*spanList
	// 页表
	pm *pageMap
}

func newPageHeap(pm *pageMap) *pageHeap {
	ph := &pageHeap{
		freeSpans: [pageHeapArraySize]*spanList{},
		pm:        pm,
	}
	for i := 0; i < pageHeapArraySize; i++ {
		ph.freeSpans[i] = newList()
	}
	return ph
}

func (p *pageHeap) fill(numPages int) error {
	if numPages < 1 || numPages > maxNumPages {
		return fmt.Errorf("page numbers: %d invalid", numPages)
	}
	// 找到 >= numPage 的 span 然后分裂 span，得到合适的 span
	for i := numPages; i < pageHeapArraySize; i++ {
		freeList := p.freeSpans[i]
		if !freeList.isEmpty() {
			s := freeList.pop()
			ns := s.split(numPages, p)
			if ns == nil {
				return fmt.Errorf("split span err")
			}
			err := p.insertSpan(s)
			if err != nil {
				return err
			}
			return p.insertSpan(ns)
		}
	}
	// find larger span and split
	freeSpan := p.freeSpans[pageHeapArraySize-1]
	if !freeSpan.isEmpty() {
		for iter := freeSpan.begin(); iter != freeSpan.end(); iter = iter.next {
			if iter.numPages > numPages {
				freeSpan.remove(iter)
				ns := iter.split(numPages, p)
				if ns == nil {
					return fmt.Errorf("split span err")
				}
				err := p.insertSpan(iter)
				if err != nil {
					return err
				}
				return p.insertSpan(ns)
			}
		}
	}
	// not found, then new
	// insertAt span
	s, err := newSpan(numPages)
	if err != nil {
		return fmt.Errorf("split span err: %s", err)
	}

	return p.insertSpan(s)
}

func (p *pageHeap) getSpan(numPages int) (*span, error) {
	if numPages < pageHeapArraySize {
		freeSpan := p.freeSpans[numPages-1]
		if freeSpan.isEmpty() {
			err := p.fill(numPages) // 为空，则走 mmap 分配
			if err != nil {
				return nil, err
			}
		}
		return freeSpan.pop(), nil
	}
	// numPages >= 256
	freeSpan := p.freeSpans[pageHeapArraySize-1]
	if freeSpan.isEmpty() {
		err := p.fill(numPages)
		if err != nil {
			return nil, err
		}
	}
	for iter := freeSpan.begin(); iter != freeSpan.end(); iter = iter.next {
		if iter.numPages == numPages {
			freeSpan.remove(iter)
			return iter, nil
		}
	}

	// 没有找到 numPages 相关的 span，fill & get
	if err := p.fill(numPages); err != nil {
		return nil, err
	}
	return p.getSpan(numPages)
}

func (p *pageHeap) insertSpan(span *span) error {
	err := p.pm.insertSpan(span)
	if err != nil {
		return err
	}
	i := span.numPages - 1
	if i >= pageHeapArraySize {
		i = pageHeapArraySize - 1
	}
	p.freeSpans[i].push(span)
	span.state = free
	span.objectClass = 0
	span.freeObjects = newObjectList()
	return nil
}

func (p *pageHeap) getSpanList(span *span) *spanList {
	i := span.numPages - 1
	if i >= pageHeapArraySize {
		i = pageHeapArraySize - 1
	}
	return p.freeSpans[i]
}

func (p *pageHeap) getPageMap() *pageMap {
	return p.pm
}
