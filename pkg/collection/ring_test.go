package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRingWithElements(t *testing.T) {
	assert := assert.New(t)

	ring := NewRingWithElements([]int{5, 4, 3, 2, 1})
	eles := ring.Take()
	assert.Equal(len(eles), 5)
	assert.Equal(eles[0], 5)
}

func TestNewRing(t *testing.T) {
	assert := assert.New(t)

	ring := NewRing[int](4)
	ring.Add(1)
	eles := ring.Take()
	assert.Equal(len(eles), 1)
	assert.Equal(eles[0], 1)

	ring.Add(2)
	ring.Add(3)
	ring.Add(4)
	ring.Add(5)
	eles = ring.Take()
	assert.Equal(len(eles), 4)
	assert.Equal(eles[3], 5)
}
