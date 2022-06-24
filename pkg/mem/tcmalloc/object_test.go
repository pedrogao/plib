package tcmalloc

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestObject(t *testing.T) {
	assert := assert.New(t)

	addr, err := mmapAnonymous(1024)
	assert.Nil(err)

	a := object{}
	(*(*object)(unsafe.Pointer(addr))).next = &a

	assert.True((*object)(unsafe.Pointer(addr)).next == &a)
}

func TestObjectList(t *testing.T) {
	assert := assert.New(t)

	list := newObjectList()
	list.push(&object{})
	assert.Equal(list.len(), 1)
	obj := list.pop()
	assert.NotNil(obj)
	assert.Equal(list.len(), 0)
}
