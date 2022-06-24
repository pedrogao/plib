package main

// #cgo CFLAGS: -g -Wall
// #include "string.h"
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	s1 := C.string_create_by_str(C.CString("Hello "))
	s2 := C.string_create_by_str(C.CString("Gopher"))
	s3 := C.string_concat(s1, s2)

	defer C.free(unsafe.Pointer(s1))
	defer C.free(unsafe.Pointer(s2))
	defer C.free(unsafe.Pointer(s3))

	fmt.Printf("%v\n", s1.length)
	fmt.Printf("%v\n", s2.length)
	fmt.Printf("%v\n", s3.length)
}
