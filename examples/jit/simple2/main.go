package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

type fooFunc func(x, y int64) int64

func main() {
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
	fooFunction, err := hex.DecodeString("4889c64889c7480faffe574889de4889df480faffe4889fe5f4829f74889f8c3")
	if err != nil {
		log.Fatalln(err)
	}

	executableFooFunc, err := syscall.Mmap(
		-1,
		0,
		128,
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC,
		syscall.MAP_PRIVATE|syscall.MAP_ANON) // linux for MAP_ANONYMOUS
	if err != nil {
		log.Fatalf("mmap err: %v", err)
	}

	for i := range fooFunction {
		executableFooFunc[i] = fooFunction[i]
	}

	unsafePrintFunc := (uintptr)(unsafe.Pointer(&executableFooFunc))
	printer := *(*fooFunc)(unsafe.Pointer(&unsafePrintFunc))
	ret := printer(3, 2)
	fmt.Println(ret)
}
