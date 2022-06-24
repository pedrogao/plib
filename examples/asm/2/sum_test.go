package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:generate go run asm.go -out sum.s -stubs stub.go

func TestSum(t *testing.T) {
	assert := assert.New(t)

	sum := Sum([]uint64{1, 2, 3})
	assert.Equal(sum, uint64(6))
}
