package marksweep

import (
	"fmt"
	"unsafe"

	"github.com/pedrogao/plib/pkg/mem/tcmalloc"
)

const (
	stackMax           = 256
	initialGCThreshold = 8
)

// ObjectType 对象类型
type ObjectType int

const (
	// ObjInt int 对象
	ObjInt ObjectType = iota + 1
	// ObjPair 元组对象
	ObjPair
)

type (
	// Object 对象
	Object struct {
		typ    ObjectType // 8
		next   *Object    // 8
		marked bool       // 8
		inner  any        // 16
	}

	// Pair 元组对象
	Pair struct {
		head *Object // 8
		tail *Object // 8
	}

	// VM 虚拟机
	VM struct {
		numObjects int
		maxObjects int
		link       *Object
		stack      [stackMax]*Object
		stackSize  int
	}
)

var (
	objSiz  = int(unsafe.Sizeof(Object{}))
	pairSiz = int(unsafe.Sizeof(Pair{}))
)

// NewVM 新建虚拟机
func NewVM() *VM {
	vm := &VM{
		numObjects: 0,
		maxObjects: initialGCThreshold,
		stack:      [256]*Object{},
		stackSize:  0,
	}
	return vm
}

// Free all objects
func (v *VM) Free() {
	v.stackSize = 0
	v.GC()
}

// GC 垃圾灰色
func (v *VM) GC() {
	numObjects := v.numObjects
	if numObjects == 0 {
		return
	}

	// mark
	v.markAll()
	// sweep
	v.sweep()

	if v.numObjects == 0 {
		v.maxObjects = initialGCThreshold
	} else {
		v.maxObjects *= 2
	}
	fmt.Printf("Collected %d objects, %d remaining.\n",
		numObjects-v.numObjects, v.numObjects)
}

func (v *VM) markAll() {
	for i := 0; i < v.stackSize; i++ {
		v.stack[i].mark()
	}
}

func (v *VM) sweep() {
	obj := v.link
	for obj != nil {
		if !obj.marked {
			unreached := obj
			obj = unreached.next
			addr := uintptr(unsafe.Pointer(unreached))
			tcmalloc.Free(addr)
			v.numObjects--
		} else {
			obj.marked = false
			obj = obj.next
		}
	}
}

func (v *VM) push(value *Object) {
	v.stack[v.stackSize] = value
	v.stackSize++
}

func (v *VM) pushInt(n int) {
	obj := v.newObject(ObjInt)
	obj.inner = n

	v.push(obj)
}

func (v *VM) pushPair() *Object {
	obj := v.newObject(ObjPair)

	//obj.inner = &Pair{
	//	head: v.pop(),
	//	tail: v.pop(),
	//}

	addr := tcmalloc.Malloc(pairSiz)
	pair := (*Pair)(unsafe.Pointer(addr))
	pair.head = v.pop()
	pair.tail = v.pop()
	obj.inner = pair

	v.push(obj)
	return obj
}

func (v *VM) pop() *Object {
	v.stackSize--
	return v.stack[v.stackSize]
}

func (v *VM) newObject(typ ObjectType) *Object {
	if v.numObjects == v.maxObjects {
		v.GC()
	}

	//obj := &Object{
	//	typ:    typ,
	//	next:   v.link,
	//	marked: false,
	//}

	addr := tcmalloc.Malloc(objSiz)
	obj := (*Object)(unsafe.Pointer(addr))
	obj.typ = typ
	obj.next = v.link
	obj.marked = false

	v.link = obj
	v.numObjects++

	return obj
}

func (v *VM) NewInt(n int) int {
	obj := v.newObject(ObjInt)
	obj.inner = n

	v.push(obj)

	return obj.inner.(int)
}

func (o *Object) mark() {
	if o.marked {
		return
	}
	o.marked = true
	if o.typ == ObjPair {
		p := o.inner.(*Pair)
		p.head.mark()
		p.tail.mark()
	}
}
