package skiplist

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkipList_Operations(t *testing.T) {
	skipList := NewSkipList[int, string]()
	for i := 1; i <= 1000; i++ {
		skipList.Insert(i, strconv.FormatInt(int64(i), 10))
	}

	for i := 1; i <= 1000; i++ {
		find, b := skipList.Find(i)
		assert.Equal(t, b, true)
		assert.Equal(t, find, strconv.FormatInt(int64(i), 10))
	}

	for i := 1; i <= 1000; i++ {
		skipList.Delete(i)
	}

	for i := 1; i <= 1000; i++ {
		_, b := skipList.Find(i)
		assert.Equal(t, b, false)
	}
}
