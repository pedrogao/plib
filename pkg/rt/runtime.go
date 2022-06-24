package rt

import (
	"reflect"
	"unsafe"
)

type GoType struct {
	Size       uintptr
	PtrData    uintptr
	Hash       uint32
	Flags      uint8
	Align      uint8
	FieldAlign uint8
	KindFlags  uint8
	Traits     unsafe.Pointer
	GCData     *byte
	Str        int32
	PtrToSelf  int32
}

type GoSlice struct {
	Ptr unsafe.Pointer
	Len int
	Cap int
}

type GoString struct {
	Ptr unsafe.Pointer
	Len int
}

type GoItab struct {
	it unsafe.Pointer
	Vt *GoType
	hv uint32
	_  [4]byte
	fn [1]uintptr
}

// refer: https://qcrao91.gitbook.io/go/interface/iface-he-eface-de-qu-bie-shi-shi-mo

type GoIface struct {
	Itab  *GoItab
	Value unsafe.Pointer
}

// GoEface 本质是 interface
type GoEface struct {
	Type  *GoType
	Value unsafe.Pointer
}

func UnpackType(t reflect.Type) *GoType {
	return (*GoType)((*GoEface)(unsafe.Pointer(&t)).Value)
}

//go:linkname growslice runtime.growslice
//goland:noinspection GoUnusedParameter
func growslice(et *GoType, old GoSlice, cap int) GoSlice

//go:noescape
//go:linkname memmove runtime.memmove
//goland:noinspection GoUnusedParameter
func memmove(to unsafe.Pointer, from unsafe.Pointer, n uintptr)
