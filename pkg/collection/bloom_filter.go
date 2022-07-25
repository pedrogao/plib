package collection

import (
	"fmt"
	"math"

	jsoniter "github.com/json-iterator/go"
	"github.com/pedrogao/plib/pkg/hash"
)

// For an explanation of the math, visit
// https://en.wikipedia.org/wiki/Bloom_filter#Probability_of_false_positives.
// A condensed explanation can be found here:
// https://stackoverflow.com/questions/658439/how-many-hash-functions-does-my-bloom-filter-need

// BloomFilter 布隆过滤器
type BloomFilter struct {
	falsePositivePob        float64
	bitArraySize, hashCount int
	bit                     *BitArray
}

// NewBloomFilter 新建布隆过滤器
func NewBloomFilter(numItems int, falsePositivePob float64) *BloomFilter {
	bf := &BloomFilter{
		falsePositivePob: falsePositivePob,
		bit:              NewBitArray(numItems),
	}
	bitArraySize := bf.calculateBitArraySize(numItems, falsePositivePob)
	hashCount := bf.calculateHashCount(bitArraySize, numItems)
	bf.hashCount = hashCount
	bf.bitArraySize = bitArraySize
	return bf
}

func (f *BloomFilter) Add(item string) {
	data := []byte(item)
	for i := 0; i < f.hashCount; i++ {
		digest := int(hash.Murmur332(data, uint32(i))) % f.bitArraySize
		f.bit.Add(digest)
	}
}

func (f *BloomFilter) Check(item string) bool {
	data := []byte(item)
	for i := 0; i < f.hashCount; i++ {
		digest := int(hash.Murmur332(data, uint32(i))) % f.bitArraySize
		if !f.bit.Has(digest) {
			return false
		}
	}
	return true
}

func (f *BloomFilter) calculateBitArraySize(numItems int,
	probability float64) int {
	// m = -(n * lg(p)) / (lg(2)^2)
	m := -(float64(numItems) * math.Log(probability) / (math.Pow(math.Log(2), 2)))
	return int(m)
}

func (f *BloomFilter) calculateHashCount(bitArraySize, numItems int) int {
	// k = (m/n) * lg(2)
	k := (float64(bitArraySize) / float64(numItems)) * math.Log(2)
	return int(k)
}

func (f *BloomFilter) Pack() []byte {
	bytes, err := jsoniter.Marshal(map[string]any{
		"falsePositivePob": f.falsePositivePob,
		"bitArraySize":     f.bitArraySize,
		"hashCount":        f.hashCount,
		"bit":              f.bit.Pack(),
	})
	if err != nil {
		panic(err)
	}

	return bytes
}

func (f *BloomFilter) UnPack(data []byte) error {
	m := map[string]any{}
	var err error
	err = jsoniter.Unmarshal(data, &m)
	if err != nil {
		return fmt.Errorf("unmarshal bloom filter err: %s", err)
	}
	bit := &BitArray{}
	f.falsePositivePob = m["falsePositivePob"].(float64)
	f.bitArraySize = m["bitArraySize"].(int)
	f.hashCount = m["hashCount"].(int)

	err = bit.UnPack(m["bit"].([]byte))
	if err != nil {
		return err
	}
	f.bit = bit

	return nil
}
