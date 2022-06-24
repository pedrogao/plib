package main

import (
	"unsafe"
)

//go:noescape
func __sum_float64(buf, len, res unsafe.Pointer)
