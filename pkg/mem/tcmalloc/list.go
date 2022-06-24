package tcmalloc

type (
	spanList struct {
		head   *span
		tail   *span
		length int
	}
)

func newList() *spanList {
	l := &spanList{}
	return l
}

func (l *spanList) begin() *span {
	return l.head
}

func (l *spanList) end() *span {
	return l.tail
}

func (l *spanList) push(n *span) {
	if n == nil {
		return
	}
	l.length++
	if l.head == nil {
		l.head = n
		l.tail = n
		n.prev = nil
		n.next = nil
		return
	}
	// 加到链表尾
	n.prev = l.tail
	n.next = nil
	l.tail.next = n
	l.tail = n
}

func (l *spanList) remove(n *span) {
	if n == nil {
		return
	}
	if n.next != nil {
		n.next.prev = n.prev
	}
	if n.prev != nil {
		n.prev.next = n.next
	}
	n.next = nil
	n.prev = nil
	l.length--
	if l.length == 0 {
		l.head = nil
		l.tail = nil
	}
}

func (l *spanList) pop() *span {
	end := l.end()
	l.remove(end)
	return end
}

func (l *spanList) len() int {
	return l.length
}

func (l *spanList) isEmpty() bool {
	return l.length <= 0
}
