package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer_lexWhitespace(t *testing.T) {
	assert := assert.New(t)

	type fields struct {
		source   string
		curIndex int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "simple spaces",
			fields: fields{
				source:   "     ",
				curIndex: 5,
			},
		},
		{
			name: "spaces with name prefix",
			fields: fields{
				source:   "pedro     ",
				curIndex: 0,
			},
		},
		{
			name: "spaces with name suffix",
			fields: fields{
				source:   "     pedro",
				curIndex: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.fields.source)
			l.lexWhitespace()
			curIndex := l.curIndex()
			assert.Equal(curIndex, tt.fields.curIndex)
		})
	}
}

func TestLexer_lexSyntax(t *testing.T) {
	assert := assert.New(t)

	type fields struct {
		source string
		index  int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "[",
			fields: fields{
				source: "[",
				index:  1,
			},
		},
		{
			name: "]",
			fields: fields{
				source: "]",
				index:  1,
			},
		},
		{
			name: "{",
			fields: fields{
				source: "{",
				index:  1,
			},
		},
		{
			name: "}",
			fields: fields{
				source: "}",
				index:  1,
			},
		},
		{
			name: ":",
			fields: fields{
				source: ":",
				index:  1,
			},
		},
		{
			name: ",",
			fields: fields{
				source: ",",
				index:  1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.fields.source)

			got, err := l.lexSyntax()
			assert.NoError(err)
			assert.Equal(got.value, tt.name)
			assert.Equal(got.typ, SyntaxToken)
			assert.Equal(got.location, 0)
			assert.Equal(l.curIndex(), tt.fields.index)
		})
	}
}

func TestLexer_lexString(t *testing.T) {
	assert := assert.New(t)

	type fields struct {
		source string
		index  int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "pedro",
			fields: fields{
				source: "\"pedro\"",
				index:  7,
			},
			want: "pedro",
		},
		{
			name: "pedro mike",
			fields: fields{
				source: "\" pedro mike \"",
				index:  14,
			},
			want: " pedro mike ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			l := NewLexer(tt.fields.source)
			got, err := l.lexString()
			assert.NoError(err)
			assert.Equal(got.value, tt.want)
			assert.Equal(got.typ, StringToken)
			assert.Equal(got.location, 0)
			assert.Equal(l.curIndex(), tt.fields.index)
		})
	}
}

func TestLexer_lexNumber(t *testing.T) {
	assert := assert.New(t)

	type fields struct {
		source string
		index  int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "91",
			fields: fields{
				source: "91",
				index:  2,
			},
			want: "91",
		},
		{
			name: "91 whitespace",
			fields: fields{
				source: "91 ",
				index:  2,
			},
			want: "91",
		},
		{
			name: "9 whitespace 1",
			fields: fields{
				source: "9 1",
				index:  1,
			},
			want: "9",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.fields.source)
			got, err := l.lexNumber()

			assert.NoError(err)
			assert.Equal(got.value, tt.want)
			assert.Equal(got.typ, NumberToken)
			assert.Equal(got.location, 0)
			assert.Equal(l.curIndex(), tt.fields.index)
		})
	}
}

func TestLexer_lexNull(t *testing.T) {
	assert := assert.New(t)

	type fields struct {
		source string
		index  int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "null",
			fields: fields{
				source: "null",
				index:  4,
			},
			want: "null",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.fields.source)
			got, err := l.lexNull()

			assert.NoError(err)
			assert.Equal(got.value, tt.want)
			assert.Equal(got.typ, NullToken)
			assert.Equal(got.location, 0)
			assert.Equal(l.curIndex(), tt.fields.index)
		})
	}
}

func TestLexer_lexTrue(t *testing.T) {
	assert := assert.New(t)

	type fields struct {
		source string
		index  int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "true",
			fields: fields{
				source: "true",
				index:  4,
			},
			want: "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.fields.source)
			got, err := l.lexTrue()

			assert.NoError(err)
			assert.Equal(got.value, tt.want)
			assert.Equal(got.typ, BooleanToken)
			assert.Equal(got.location, 0)
			assert.Equal(l.curIndex(), tt.fields.index)
		})
	}
}

func TestLexer_lexFalse(t *testing.T) {
	assert := assert.New(t)

	type fields struct {
		source string
		index  int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "false",
			fields: fields{
				source: "false",
				index:  5,
			},
			want: "false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.fields.source)
			got, err := l.lexFalse()

			assert.NoError(err)
			assert.Equal(got.value, tt.want)
			assert.Equal(got.typ, BooleanToken)
			assert.Equal(got.location, 0)
			assert.Equal(l.curIndex(), tt.fields.index)
		})
	}
}

var testJson = `{
	"glossary": {
		"GlossDiv": {
		"GlossList": {
		"GlossEntry": {
		"Abbrev": "ISO 8879:1986",
		"Acronym": "SGML",
		"GlossDef": {
			"GlossSeeAlso": [
				"GML",
				"XML"
			],
		"para": "A meta-markup language, used to create markup languages such as DocBook."
		},
		"GlossSee": "markup",
		"GlossTerm": "Standard Generalized Markup Language",
		"ID": "SGML",
		"SortAs": "SGML"
	}
	},
	"title": "S"
	},
	"title": "example glossary"
	}
}`

func TestLexer_Lex(t *testing.T) {
	assert := assert.New(t)

	type fields struct {
		source string
	}
	tests := []struct {
		name      string
		fields    fields
		wantLen   int
		wantIndex int
	}{
		{
			name: "simple json str",
			fields: fields{
				source: "{\"name\": \"pedro\", \"age\": 25 } ",
			},
			wantLen:   9,
			wantIndex: 30,
		},
		{
			name: "simple json str with space",
			fields: fields{
				source: "{\n\t\n\n\t   \"name\":    \"pedro\",   \"age\": 25 } ",
			},
			wantLen:   9,
			wantIndex: 43,
		},
		{
			name: "complex json str",
			fields: fields{
				source: testJson,
			},
			wantLen:   65,
			wantIndex: 444,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.fields.source)

			got, err := l.Lex()
			assert.NoError(err)
			assert.Equal(tt.wantLen, len(got), "Lex()")
			assert.Equal(tt.wantIndex, l.curIndex(), "Lex()")
		})
	}
}
