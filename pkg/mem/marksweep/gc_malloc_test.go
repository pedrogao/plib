package marksweep

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestVM_GC1(t *testing.T) {
	assert := assert.New(t)

	vm := NewVM()
	vm.pushInt(1)
	vm.pushInt(2)

	vm.GC()
	assert.Equal(vm.numObjects, 2)
	vm.Free()
}

func TestVM_GC2(t *testing.T) {
	assert := assert.New(t)

	vm := NewVM()
	vm.pushInt(1)
	vm.pushInt(2)
	vm.pop()
	vm.pop()

	vm.GC()
	assert.Equal(vm.numObjects, 0)
	vm.Free()
}

func TestVM_GC3(t *testing.T) {
	assert := assert.New(t)

	vm := NewVM()
	vm.pushInt(1)
	vm.pushInt(2)
	vm.pushPair()
	vm.pushInt(3)
	vm.pushInt(4)
	vm.pushPair()
	vm.pushPair()

	vm.GC()
	assert.Equal(vm.numObjects, 7)
	vm.Free()
}

func TestVM_GC4(t *testing.T) {
	assert := assert.New(t)

	vm := NewVM()
	vm.pushInt(1)
	vm.pushInt(2)
	a := vm.pushPair()
	vm.pushInt(3)
	vm.pushInt(4)
	b := vm.pushPair()

	a.inner.(*Pair).tail = b
	b.inner.(*Pair).tail = a
	assert.Equal(vm.numObjects, 6)
	vm.GC()
	assert.Equal(vm.numObjects, 4)
	vm.Free()
}

func TestVM_ObjSize(t *testing.T) {
	assert := assert.New(t)

	type t1 struct {
		typ  ObjectType // 8
		next *Object    // 8
	}

	sz := int(unsafe.Sizeof(t1{}))

	assert.Equal(sz, 16)

	type t2 struct {
		typ    ObjectType // 8
		next   *Object    // 8
		marked bool
	}

	sz = int(unsafe.Sizeof(t2{}))

	assert.Equal(sz, 24)

	type t3 struct {
		typ    ObjectType // 8
		next   *Object    // 8
		marked bool
		inner  interface{}
	}

	sz = int(unsafe.Sizeof(t3{}))

	assert.Equal(sz, 40)

	type t4 struct {
		inner interface{}
	}

	sz = int(unsafe.Sizeof(t4{}))

	assert.Equal(sz, 16)

	sz = int(unsafe.Sizeof(Object{}))

	assert.Equal(sz, 40)

	sz = int(unsafe.Sizeof(Pair{}))

	assert.Equal(sz, 16)
}

func TestVM_NewInt(t *testing.T) {
	assert := assert.New(t)

	vm := NewVM()

	i1 := vm.NewInt(2000)
	assert.Equal(i1, 2000)
	i2 := vm.NewInt(3000)
	assert.Equal(i2, 3000)

	vm.NewInt(4000)
	vm.NewInt(5000)
	vm.NewInt(6000)

	assert.Equal(vm.numObjects, 5)
	// 理论上来说，从根对象(全局变量、局部变量等)开始，进行扫描
	// 但我们没有实现另一门语言，因此使用栈来模拟根对象
	// 即：在栈上的即为根对象，不在栈上的，就不是根对象
	vm.GC()
	assert.Equal(vm.numObjects, 5)
	vm.Free()
}
