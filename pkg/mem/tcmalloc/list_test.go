package tcmalloc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	assert := assert.New(t)

	l := newList()
	l.push(&span{
		pageStart: 5,
		numPages:  5,
	})
	l.push(&span{
		pageStart: 10,
		numPages:  10,
	})
	count := l.len()
	assert.Equal(count, 2)
	l.pop()
	count = l.len()
	assert.Equal(count, 1)
	l.pop()
	assert.True(l.isEmpty())
}
