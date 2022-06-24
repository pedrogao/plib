package py

import (
	"fmt"
)

type (
	Opcode struct {
		Code string
		Arg  interface{}
	}

	SSA struct {
		Action     string
		Arg1, Arg2 interface{}
	}

	Compiler struct {
		Opcodes   []*Opcode
		Constants []interface{}
		index     int
	}
)

func (s *SSA) String() string {
	return fmt.Sprintf("%s-%v-%v", s.Action, s.Arg1, s.Arg2)
}

func NewCompiler(opcodes []*Opcode, constants []interface{}) *Compiler {
	return &Compiler{
		Opcodes:   opcodes,
		Constants: constants,
		index:     0,
	}
}

func (c *Compiler) fetch() *Opcode {
	op := c.Opcodes[c.index]
	c.index++
	return op
}

func (c *Compiler) variable(n int) string {
	order := []string{"rax", "rbx", "rcx", "rdi"}
	return order[n]
}

func (c *Compiler) Compile() ([]*SSA, error) {
	ir := make([]*SSA, 0)

	pushSSA := func(a string, b, c interface{}) {
		ssa := &SSA{
			Action: a,
			Arg1:   b,
			Arg2:   c,
		}
		ir = append(ir, ssa)
	}

	for c.index < len(c.Opcodes) {
		op := c.fetch()
		// rdi, rsi
		switch op.Code {
		case "LOAD_FAST":
			pushSSA("push", c.variable(op.Arg.(int)), nil)
		case "STORE_FAST":
			pushSSA("pop", "rdi", nil)
			pushSSA("move", c.variable(op.Arg.(int)), "rdi")
		case "LOAD_CONST":
			pushSSA("immediate", "rdi", c.Constants[op.Arg.(int)])
			pushSSA("push", "rdi", nil)
		case "BINARY_MULTIPLY":
			pushSSA("pop", "rdi", nil)
			pushSSA("pop", "rsi", nil)
			pushSSA("imul", "rdi", "rsi")
			pushSSA("push", "rdi", nil)
		case "BINARY_ADD", "INPLACE_ADD":
			pushSSA("pop", "rdi", nil)
			pushSSA("pop", "rsi", nil)
			pushSSA("add", "rdi", "rsi")
			pushSSA("push", "rdi", nil)
		case "BINARY_SUBTRACT", "INPLACE_SUBTRACT":
			pushSSA("pop", "rsi", nil)
			pushSSA("pop", "rdi", nil)
			pushSSA("sub", "rdi", "rsi")
			pushSSA("push", "rdi", nil)
		case "UNARY_NEGATIVE":
			pushSSA("pop", "rdi", nil)
			pushSSA("neg", "rdi", nil)
			pushSSA("push", "rdi", nil)
		case "RETURN_VALUE":
			pushSSA("pop", "rax", nil)
			pushSSA("ret", nil, nil)
		default:
			return nil, fmt.Errorf("%s not support", op.Code)
		}
	}
	return ir, nil
}
