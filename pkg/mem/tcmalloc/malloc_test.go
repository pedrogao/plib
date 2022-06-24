package tcmalloc

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestMalloc1(t *testing.T) {
	assert := assert.New(t)

	p1 := Malloc(100)
	assert.NotEqual(int(p1), 0)
	t.Logf("%d\n", p1)

	type person struct {
		age  int      // 8
		name [20]byte // 20
		// align
	}

	t.Logf("size: %d\n", unsafe.Sizeof(person{}))

	p := (*person)(unsafe.Pointer(p1))
	// pedro
	p.name = [20]byte{'p', 'e', 'd', 'r', 'o'}
	p.age = 26

	assert.Equal(p.age, 26)
	assert.Equal(string(p.name[:5]), "pedro")
}

func TestMalloc2(t *testing.T) {
	assert := assert.New(t)

	p1 := Malloc(33768)
	assert.NotEqual(int(p1), 0)
	t.Logf("%d\n", p1)

	type person struct {
		age  int      // 8
		name [20]byte // 20
		// align
	}

	t.Logf("size: %d\n", unsafe.Sizeof(person{}))

	p := (*person)(unsafe.Pointer(p1))
	// pedro
	p.name = [20]byte{'p', 'e', 'd', 'r', 'o'}
	p.age = 26

	assert.Equal(p.age, 26)
	assert.Equal(string(p.name[:5]), "pedro")
}

func TestFree1(t *testing.T) {
	assert := assert.New(t)

	p1 := Malloc(100)
	assert.NotEqual(int(p1), 0)
	t.Logf("%d\n", p1)

	Free(p1)
}

func TestFree2(t *testing.T) {
	assert := assert.New(t)

	p1 := Malloc(33768)
	assert.NotEqual(int(p1), 0)
	t.Logf("%d\n", p1)

	Free(p1)
}

func TestFree3(t *testing.T) {
	assert := assert.New(t)

	p1 := Malloc(100)
	assert.NotEqual(int(p1), 0)
	t.Logf("%d\n", p1)
	Free(p1)

	p2 := Malloc(100)
	assert.NotEqual(int(p2), 0)
	t.Logf("%d\n", p2)
	Free(p2)

	p3 := Malloc(200)
	assert.NotEqual(int(p3), 0)
	t.Logf("%d\n", p3)
	Free(p3)

	p4 := Malloc(500)
	assert.NotEqual(int(p4), 0)
	t.Logf("%d\n", p4)
	Free(p4)
}

func TestFree4(t *testing.T) {
	assert := assert.New(t)

	pointers := make([]uintptr, 0, 1000)
	for i := 0; i < 1000; i++ {
		p1 := Malloc(100)
		assert.NotEqual(int(p1), 0)
		t.Logf("%d\n", p1)
		pointers = append(pointers, p1)
	}

	for i := 0; i < 1000; i++ {
		Free(pointers[i])
	}
}
