package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitArray(t *testing.T) {
	assert := assert.New(t)
	// assert.Equal(bitnum, 64)

	bit := NewBitArray(100)

	assert.Equal(bit.Len(), 0)
	assert.Equal(len(bit.data), 100/bitnum)

	bit.Add(0)
	bit.Add(1)
	bit.Add(5)
	assert.Equal(bit.Len(), 3)
	assert.True(bit.Has(5))
	bit.Remove(5)
	assert.Equal(bit.Len(), 2)
	assert.False(bit.Has(5))

	bit.Clear()
	assert.Equal(bit.Len(), 0)
	bit.AddAll(10, 7, 2, 4)
	assert.Equal(bit.Len(), 4)

	o := NewBitArray(10)
	o.AddAll(10, 3, 9)
	assert.Equal(o.Len(), 3)

	bit.UnionWith(o)

	assert.Equal(bit.Len(), 6)

	t.Logf("%s\n", bit.String())
}
