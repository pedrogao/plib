package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBloomFilterOps(t *testing.T) {
	assert := assert.New(t)

	bf := NewBloomFilter(1, 0.05)
	bf.Add("pedro")

	assert.True(bf.Check("pedro"))

	bf = NewBloomFilter(4, 0.05)
	bf.Add("pedro")
	bf.Add("sara")
	bf.Add("mike")
	bf.Add("leo")

	assert.True(bf.Check("pedro"))
	assert.True(bf.Check("mike"))
	assert.True(bf.Check("sara"))
	assert.True(bf.Check("leo"))

	assert.False(bf.Check("berry"))
	assert.False(bf.Check("lucy"))
}

func TestBloomFilterDifferentProbability(t *testing.T) {
	assert := assert.New(t)

	bf := NewBloomFilter(1, 0.15)
	bf.Add("pedro")
	assert.True(bf.Check("pedro"))

	bf = NewBloomFilter(1, 0.5)
	bf.Add("pedro")
	assert.True(bf.Check("pedro"))

	bf = NewBloomFilter(1, 0.9)
	bf.Add("pedro")
	assert.True(bf.Check("pedro"))
}

func TestBloomFilterHashCount(t *testing.T) {
	assert := assert.New(t)

	bf := NewBloomFilter(20, 0.05)
	assert.Equal(bf.hashCount, 4)

	bf = NewBloomFilter(1000, 0.25)
	assert.Equal(bf.hashCount, 1)

	bf = NewBloomFilter(10000, 0.02)
	assert.Equal(bf.hashCount, 5)
}

func TestBloomFilterBitArraySize(t *testing.T) {
	assert := assert.New(t)

	bf := NewBloomFilter(20, 0.05)
	assert.Equal(bf.bitArraySize, 124)

	bf = NewBloomFilter(1000, 0.25)
	assert.Equal(bf.bitArraySize, 2885)

	bf = NewBloomFilter(10000, 0.02)
	assert.Equal(bf.bitArraySize, 81423)
}
