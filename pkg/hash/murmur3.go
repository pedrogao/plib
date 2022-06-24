package hash

import (
	"encoding/binary"
)

// MurmurHash
// refer:
//   - https://en.wikipedia.org/wiki/MurmurHash
//   - https://segmentfault.com/a/1190000016743220

func murmur32Scramble(k uint32) uint32 {
	k *= 0xcc9e2d51
	k = (k << 15) | (k >> 17)
	k *= 0x1b873593
	return k
}

// Murmur332 murmur hash for uint32
func Murmur332(key []byte, seed uint32) uint32 {
	length := len(key)
	h := seed
	var k uint32
	offset := 0
	/* Read in groups of 4. */
	for i := length >> 2; i > 0; i-- {
		k = binary.LittleEndian.Uint32(key[offset:])
		offset += 4 // 32 = 8 * 4
		h ^= murmur32Scramble(k)
		h = (h << 13) | (h >> 19)
		h = h*5 + 0xe6546b64
	}
	/* Read the rest. */
	k = 0
	for i := length & 3; i > 0; i-- {
		k <<= 8
		k |= uint32(key[i-1])
	}
	// A swap is *not* necessary here because the preceding loop already
	// places the low bytes in the low places according to whatever endianness
	// we use. Swaps only apply when the memory is copied in a chunk.
	h ^= murmur32Scramble(k)
	/* Finalize. */
	h ^= uint32(length)
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}
