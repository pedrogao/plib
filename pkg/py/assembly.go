package py

import (
	"encoding/binary"
)

type Assembler struct {
	index int
	code  []byte
}

func NewAssembler(size int) *Assembler {
	return &Assembler{
		index: 0,
		code:  make([]byte, size),
	}
}

func (as *Assembler) Assembly(ir []*SSA) []byte {
	for _, ssa := range ir {
		name, a, b := ssa.Action, ssa.Arg1, ssa.Arg2
		switch name {
		case "ret":
			as.ret()
		case "push":
			as.push(a.(string))
		case "pop":
			as.pop(a.(string))
		case "imul":
			as.imul(a.(string), b.(string))
		case "add":
			as.add(a.(string), b.(string))
		case "sub":
			as.sub(a.(string), b.(string))
		case "neg":
			as.neg(a.(string))
		case "mov":
			as.mov(a.(string), b.(string))
		case "immediate":
			as.immediate(a.(string), b.(int))
		}
	}
	return as.code[:as.index]
}

func (as *Assembler) immediate(a string, number int) {
	nb := make([]byte, 8)
	binary.LittleEndian.PutUint64(nb, uint64(number))
	// 小端
	as.emit(0x48, 0xb8|as.register(a, ""))
	as.emit(nb...)
}

func (as *Assembler) mov(a, b string) {
	as.emit(0x48, 0x89, 0xc0|as.register(b, a))
}

func (as *Assembler) neg(a string) {
	as.emit(0x48, 0xf7, 0xd8|as.register(a, ""))
}

func (as *Assembler) sub(a, b string) {
	as.emit(0x48, 0x29, 0xc0|as.register(b, a))
}

func (as *Assembler) add(a, b string) {
	as.emit(0x48, 0x01, 0xc0|as.register(b, a))
}

func (as *Assembler) imul(a, b string) {
	as.emit(0x48, 0x0f, 0xaf, 0xc0|as.register(a, b))
}

func (as *Assembler) pop(a string) {
	as.emit(0x58 | as.register(a, ""))
}

func (as *Assembler) push(a string) {
	as.emit(0x50 | as.register(a, ""))
}

func (as *Assembler) ret() {
	as.emit(0xc3)
}

func (as *Assembler) register(a, b string) byte {
	order := map[string]byte{
		"rax": 0,
		"rcx": 1,
		"rdx": 2,
		"rbx": 3,
		"rsp": 4,
		"rbp": 5,
		"rsi": 6,
		"rdi": 7,
	}
	enc := order[a]
	if b != "" {
		enc = enc<<3 | order[b]
	}
	return enc
}

func (as *Assembler) emit(codes ...byte) {
	//fmt.Printf("emit: %v\n", codes)
	for _, code := range codes {
		as.code[as.index] = code
		as.index++
	}
}
