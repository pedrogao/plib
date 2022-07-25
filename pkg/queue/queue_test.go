package queue

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_ = os.RemoveAll("./test")

	queue, err := New("./test")
	assert.Nil(t, err)
	assert.Truef(t, queue.Empty(), "queue must be empty")

	err = queue.Push("helloworld")
	assert.Nil(t, err)
	err = queue.Push("pedro")
	assert.Nil(t, err)

	val, err := queue.Pop()
	assert.Nil(t, err)
	assert.Equalf(t, val, "helloworld", "val must be helloworld")

	val, err = queue.Pop()
	assert.Nil(t, err)
	assert.Equalf(t, val, "pedro", "val must be pedro")
}
