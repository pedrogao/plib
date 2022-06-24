package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:generate go run asm.go -out add.s -stubs stub.go

func TestAdd(t *testing.T) {
	assert := assert.New(t)
	ret := Add(122, 1)
	assert.Equal(ret, uint64(123))
}
