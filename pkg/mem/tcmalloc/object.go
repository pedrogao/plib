package tcmalloc

type (
	object struct {
		next *object // 下一个，待分配后，此处指针就为nil，可以使用，不浪费内存
	}

	// objectList 对象元数据
	objectList struct {
		head    *object
		length  int
		lowMark int
	}
)

func newObjectList() *objectList {
	return &objectList{
		head:   nil,
		length: 0,
	}
}

func (l *objectList) isEmpty() bool {
	return l.head == nil
}

func (l *objectList) push(obj *object) {
	// set obj.next = head
	obj.next = l.head
	// reset head to obj
	l.head = obj

	l.length++
}

func (l *objectList) pop() *object {
	if l.head == nil {
		return nil
	}
	head := l.head
	// remove head link
	l.head = head.next
	head.next = nil

	l.length--

	// 每次 cache gc 的时候会回收一半
	// 因此 lowMark 不能大于 length
	if l.length < l.lowMark {
		l.lowMark = l.length
	}

	// return head
	return head
}

func (l *objectList) len() int {
	return l.length
}
