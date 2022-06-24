package collection

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRBTree_FindNode(t *testing.T) {
	assert := assert.New(t)

	rbTree := NewRBTree()

	rbTree.Insert("2", 2)
	node2 := rbTree.root
	rbTree.Insert("1", 1)
	node1 := rbTree.root.left
	rbTree.Insert("4", 4)
	node4 := rbTree.root.right
	rbTree.Insert("5", 5)
	node5 := node4.right
	rbTree.Insert("9", 9)
	node9 := node5.right
	rbTree.Insert("3", 3)
	node3 := node4.left
	rbTree.Insert("6", 6)
	node6 := node9.left
	rbTree.Insert("7", 7)
	node7 := node5.right
	rbTree.Insert("95", 15)
	node15 := node9.right
	/*
	       ___5B___
	    __2R__      7R
	  1B     4B    6B 9B
	       3R         95R
	*/
	assert.Equal(rbTree.findNode("5"), node5)
	assert.Equal(rbTree.findNode("2"), node2)
	assert.Equal(rbTree.findNode("1"), node1)
	assert.Equal(rbTree.findNode("4"), node4)
	assert.Equal(rbTree.findNode("3"), node3)
	assert.Equal(rbTree.findNode("7"), node7)
	assert.Equal(rbTree.findNode("6"), node6)
	assert.Equal(rbTree.findNode("9"), node9)
	assert.Equal(rbTree.findNode("95"), node15)
	// invalid cases
	assert.Nil(rbTree.findNode("-1"))
	assert.Nil(rbTree.findNode("52454225"))
	assert.Nil(rbTree.findNode("0"))
	assert.Nil(rbTree.findNode("401"))
	assert.Nil(rbTree.findNode("3.00001"))
}

func TestRBTree_RemoveSingleDeepChild(t *testing.T) {

}

func TestRBTree_AddRemove(t *testing.T) {
	assert := assert.New(t)

	rbTree := NewRBTree()
	iters := 1000
	keys := make([]string, 0, iters)
	for i := 0; i < iters; i++ {
		key := fmt.Sprintf("%f", rand.ExpFloat64())
		val := fmt.Sprintf("%d", i)
		keys = append(keys, key)
		rbTree.Insert(key, val)
	}
	assert.Equal(rbTree.count, iters)

	for i := 0; i < iters; i++ {
		assert.True(rbTree.Contains(keys[i]))
	}

	for i := 0; i < iters; i++ {
		rbTree.Remove(keys[i])
		assert.False(rbTree.Contains(keys[i]))
	}
}
