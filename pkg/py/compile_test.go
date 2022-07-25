package py

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompilerOp(t *testing.T) {
	assert := assert.New(t)

	opcodes := []*Opcode{
		{
			Code: "LOAD_FAST",
			Arg:  0,
		},
		{
			Code: "UNARY_NEGATIVE",
			Arg:  nil,
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
	c := NewCompiler(opcodes, []any{nil})
	ir, err := c.Compile()
	assert.Nil(err)
	t.Logf("%v\n", ir)

	ir = Optimize(ir)
	// t.Logf("%v\n", ir)
	ir = Optimize(ir)
	// t.Logf("%v\n", ir)
	ir = Optimize(ir)
	t.Logf("%v\n", ir)

	assembler := NewAssembler(4096)
	code := assembler.Assembly(ir)
	t.Logf("size: %d\n", assembler.index)
	t.Logf("%s\n", hex.EncodeToString(code))
}
