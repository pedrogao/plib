package tricolor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVM_GC1(t *testing.T) {
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
