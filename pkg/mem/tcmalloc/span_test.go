package tcmalloc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newSpan(t *testing.T) {
	assert := assert.New(t)
	sp, err := newSpan(1)
	assert.Nil(err)
	assert.Equal(sp.numPages, 1)
	assert.Equal(sp.state, free)
	class := size2Class(32)
	sp.makeObjects(class)
	assert.False(sp.freeObjects.isEmpty())
	assert.Equal(sp.freeObjects.len(), 128)
	assert.Equal(sp.isConsistent(), true)
}
