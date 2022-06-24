package collection

import (
	"fmt"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

const (
	bitnum = 32 << (^uint(0) >> 63) // automatically determine whether it is 32 or 64 according to the platform
)

// BitArray refer https://developpaper.com/implementation-of-bit-array-in-golang/
type BitArray struct {
	data []uint
}

// NewBitArray init
// @param len, length of bits
func NewBitArray(len int) *BitArray {
	arr := &BitArray{
		data: make([]uint, len/bitnum),
	}
	return arr
}

// Has check x is set
func (a *BitArray) Has(x int) bool {
	first, second := x/bitnum, x%bitnum
	return first < len(a.data) && a.data[first]&(1<<second) != 0
}

// Add set x
func (a *BitArray) Add(x int) {
	first, second := x/bitnum, x%bitnum
	for first >= len(a.data) {
		a.data = append(a.data, 0)
	}
	a.data[first] |= 1 << second
}

// AddAll add many
func (a *BitArray) AddAll(arr ...int) {
	for _, v := range arr {
		a.Add(v)
	}
}

// UnionWith merge with other
func (a *BitArray) UnionWith(other *BitArray) {
	for i, v := range other.data {
		if i < len(a.data) {
			a.data[i] |= v
		} else {
			a.data = append(a.data, v)
		}
	}
}

// Len return length of bits
func (a *BitArray) Len() int {
	length := 0
	for _, v := range a.data {
		for j := 0; j < bitnum; j++ {
			if v&(1<<j) != 0 {
				length++
			}
		}
	}
	return length
}

func (a *BitArray) Remove(x int) {
	if !a.Has(x) {
		return
	}
	first, second := x/bitnum, x%bitnum
	a.data[first] ^= 1 << second
}

func (a *BitArray) Clear() {
	a.data = nil
}

// String returns the set as a string of the form "{1 2 3}".
func (a *BitArray) String() string {
	var buf strings.Builder
	buf.WriteByte('{')
	for i, word := range a.data {
		if word == 0 {
			continue
		}
		for j := 0; j < bitnum; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				buf.WriteString(strconv.Itoa(bitnum*i + j))
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func (a *BitArray) Pack() []byte {
	bytes, err := jsoniter.Marshal(a.data)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (a *BitArray) UnPack(data []byte) error {
	err := jsoniter.Unmarshal(data, &a.data)
	if err != nil {
		return fmt.Errorf("unmarshal bit array err: %s", err)
	}

	return nil
}
