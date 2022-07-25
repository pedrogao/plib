package lsm

import (
	"encoding/binary"
)

// SizedMap map with size
type SizedMap struct {
	inner     map[string]any
	totalSize int
}

// NewSizedMap new sized map
func NewSizedMap() *SizedMap {
	return &SizedMap{
		inner:     map[string]any{},
		totalSize: 0,
	}
}

// Get k
func (m *SizedMap) Get(key string) any {
	return m.inner[key]
}

// Set k->v
func (m *SizedMap) Set(key string, v any) {
	old := m.Get(key)
	if old != nil {
		m.totalSize -= binary.Size(old)
	}
	size := len(key) + binary.Size(v)
	m.totalSize += size
	m.inner[key] = v
}

// GetTotalSize total size of kvs
func (m *SizedMap) GetTotalSize() int {
	return m.totalSize
}
