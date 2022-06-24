package main

import (
	"encoding/hex"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/pedrogao/plib/pkg/py"
)

func main() {
	opcodes := []*py.Opcode{
		{
			Code: "LOAD_FAST",
			Arg:  0,
		},
		{
			Code: "LOAD_FAST",
			Arg:  0,
		},
		{
			Code: "BINARY_MULTIPLY",
			Arg:  nil,
		},
		{
			Code: "LOAD_FAST",
			Arg:  1,
		},
		{
			Code: "LOAD_FAST",
			Arg:  1,
		},
		{
			Code: "BINARY_MULTIPLY",
			Arg:  nil,
		},
		{
			Code: "BINARY_SUBTRACT",
			Arg:  nil,
		},
		{
			Code: "RETURN_VALUE",
			Arg:  nil,
		},
	}
	c := py.NewCompiler(opcodes, []interface{}{nil})
	ir, err := c.Compile()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", ir)
	ir = py.Optimize(ir)
	ir = py.Optimize(ir)
	fmt.Printf("%v\n", ir)

	assembler := py.NewAssembler(4096)
	code := assembler.Assembly(ir)
	fmt.Printf("%s\n", hex.EncodeToString(code))

	area, err := syscall.Mmap(
		-1,
		0,
		4096,
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC,
		syscall.MAP_PRIVATE|syscall.MAP_ANON) // linux for MAP_ANONYMOUS
	if err != nil {
		fmt.Printf("mmap err: %v", err)
	}

	for i, by := range code {
		area[i] = by
	}

	fmt.Printf("%p\n", &area)

	type fooFunc func(x, y int64) int64
	// https://xargin.com/go1-17-new-calling-convention/
	// args: rax, rbx, rcx, rdi, rsi, r8, r9, r10, r11
	// ret: rax, rbx, rcx, rdi, rsi, r8, r9, r10, r11
	// dst <- src
	// 0xc00000c030 48 89 c6       mov rsi, rax
	// 0xc00000c033 48 89 c7       mov rdi, rax
	// 0xc00000c036 48 0f af fe    imul rdi, rsi
	// 0xc00000c03a 57             push rdi
	// 0xc00000c03b 48 89 de       mov rsi, rbx
	// 0xc00000c03e 48 89 df       mov rdi, rbx
	// 0xc00000c041 48 0f af fe    imul rdi, rsi
	// 0xc00000c045 48 89 fe       mov rsi, rdi
	// 0xc00000c048 5f             pop rdi
	// 0xc00000c049 48 29 f7       sub rdi, rsi
	// 0xc00000c04c 48 89 f8       mov rax, rdi
	// 0xc00000c04f c3             ret

	unsafePrintFunc := (uintptr)(unsafe.Pointer(&area))
	foo := *(*fooFunc)(unsafe.Pointer(&unsafePrintFunc))
	// x * x - y * y
	ret := foo(3, 2)
	// -5
	fmt.Println(ret)
}
