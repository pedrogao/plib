package common

import (
	"reflect"
	"unsafe"
)

// String2Bytes string => bytes, no memory copy
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	// string to slice
	bh := reflect.SliceHeader{
		Data: sh.Data, // 直接赋值，没有拷贝
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	// 注意 []byte 是切片，不是数组
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// Bytes2String bytes to string, no memory copy
func Bytes2String(b []byte) string {
	p := unsafe.Pointer(&b) // []byte是切片
	header := (*reflect.SliceHeader)(p)
	return *(*string)(unsafe.Pointer(&header.Data))
}
