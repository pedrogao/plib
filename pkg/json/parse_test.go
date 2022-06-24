package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	assert := assert.New(t)

	type fields struct {
		source string
	}
	tests := []struct {
		name    string
		fields  fields
		wantLen int
	}{
		{
			name: "simple json str",
			fields: fields{
				source: "{\"name\": \"pedro\", \"age\": 25 } ",
			},
			wantLen: 2,
		},
		{
			name: "complex json str",
			fields: fields{
				source: "[1,2,3,5]",
			},
			wantLen: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.fields.source)
			tokens, err := l.Lex()
			assert.NoError(err)

			p := NewParser(tokens)
			node, err := p.Parse()
			assert.NoError(err)

			if node.typ == ObjectNode {
				assert.Equalf(tt.wantLen, len(node.object), "Parse()")
			} else if node.typ == ArrayNode {
				assert.Equalf(tt.wantLen, len(node.array), "Parse()")
			}
		})
	}
}
