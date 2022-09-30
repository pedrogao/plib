package syncx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type item int

func (i item) Equal(other Ordered) bool {
	n, ok := other.(item)
	if !ok {
		return false
	}
	return n == i
}

func (i item) Less(other Ordered) bool {
	n, ok := other.(item)
	if !ok {
		return false
	}
	return i < n
}

func TestOrderedListSet_Insert(t *testing.T) {
	set := NewOrderedListSet()
	set.Insert(item(1))
	assert.True(t, set.Contains(item(1)))
	set.Insert(item(2))
	assert.True(t, set.Contains(item(2)))
}
