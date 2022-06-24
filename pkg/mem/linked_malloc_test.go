package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRawMemory_HeaderSize(t *testing.T) {
	assert := assert.New(t)
	t.Logf("%d\n", headerSize)
	assert.Equal(headerSize, 24)
}

func TestRawMemory_Alloc(t *testing.T) {
	assert := assert.New(t)

	mem, err := Init(1024 * 1024 * 10)
	assert.Nil(err)
	t.Logf("%s\n", mem.Stat())

	b1, err := mem.Alloc(10)
	t.Logf("%s\n", mem.Stat())
	assert.Nil(err)
	id := (*int)(b1)
	*id = 10

	b2, err := mem.Alloc(10)
	t.Logf("%s\n", mem.Stat())
	assert.Nil(err)
	age := (*int)(b2)
	*age = 1000

	err = mem.Free(b1)
	t.Logf("%s\n", mem.Stat())
	assert.Nil(err)
	err = mem.Free(b2)
	t.Logf("%s\n", mem.Stat())
	assert.Nil(err)

	t.Logf("%s\n", mem.Stat())
}

func TestRawMemory_AllocFull(t *testing.T) {
	assert := assert.New(t)
	// 100 = 24 + 76
	mem, err := Init(100 * 100)
	assert.Nil(err)
	t.Logf("%s\n", mem.Stat())

	for i := 0; i < 100; i++ {
		b1, err := mem.Alloc(76)
		t.Logf("%s\n", mem.Stat())
		assert.Nil(err)
		id := (*int)(b1)
		*id = 10
	}

	_, err = mem.Alloc(10)
	t.Logf("%s\n", mem.Stat())
	assert.NotNil(err)
}
