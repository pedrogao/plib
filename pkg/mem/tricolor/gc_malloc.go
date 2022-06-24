package tricolor

import (
	"fmt"
	"unsafe"

	"github.com/pedrogao/plib/pkg/mem/tcmalloc"
)

const (
	stackMax           = 256
	grayMax            = 256
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
		typ    ObjectType
		next   *Object
		marked bool // black or not
		inner  interface{}
	}

	// Pair 元组对象
	Pair struct {
		head *Object
		tail *Object
	}

	// VM 虚拟机
	VM struct {
		numObjects int
		maxObjects int
		link       *Object
		stack      [stackMax]*Object // 简化方案，栈中的对象都为 GC roots
		stackSize  int

		grayStack [grayMax]*Object
		grayCount int
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
		stack:      [stackMax]*Object{},
		stackSize:  0,
		grayStack:  [grayMax]*Object{},
		grayCount:  0,
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
	v.markRoots()
	// traceReference 跟踪灰色对象
	v.traceReference()
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

func (v *VM) markRoots() {
	for i := 0; i < v.stackSize; i++ {
		v.markObject(v.stack[i])
	}
}

func (v *VM) traceReference() {
	for v.grayCount > 0 {
		v.grayCount -= 1
		object := v.grayStack[v.grayCount]
		// 出栈，然后搜索连接对象
		v.blacken(object)
	}
}

func (v *VM) sweep() {
	obj := v.link
	for obj != nil {
		if !obj.marked { // 未被标记的，就是白色对象
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

	addr := tcmalloc.Malloc(objSiz)
	obj := (*Object)(unsafe.Pointer(addr))
	obj.typ = typ
	obj.next = v.link
	obj.marked = false

	v.link = obj
	v.numObjects++

	return obj
}

func (v *VM) markObject(o *Object) {
	// 对象标记为黑色
	o.mark()
	v.grayStack[v.grayCount] = o
	// 将其加入
	v.grayCount++
}

func (o *Object) mark() {
	if o.marked {
		return
	}
	o.marked = true
}

func (v *VM) blacken(o *Object) {
	// 注意： 对象连接不是 next 指针，而是它连接的对象
	// int 对象没有连接，直接返回
	if o.typ == ObjInt {
		return
	}
	// 元组对象连接了内部的两个对象
	if o.typ == ObjPair {
		p := o.inner.(*Pair)
		v.markObject(p.head)
		v.markObject(p.tail)
	}
}
