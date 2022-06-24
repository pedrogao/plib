package json

import (
	"fmt"
)

type (
	NodeType int

	TokenType int

	// Node string|number|boolean|array|object
	Node struct {
		string  *string
		number  *float64
		boolean *bool
		array   []*Node
		object  map[string]*Node
		typ     NodeType // value type of current json node
	}

	// Token json token resp
	Token struct {
		value      string
		typ        TokenType
		location   int
		fullSource []rune
	}
)

func (j TokenType) String() string {
	switch j {
	case StringToken:
		return "String"
	case NumberToken:
		return "Number"
	case SyntaxToken:
		return "Syntax"
	case BooleanToken:
		return "Boolean"
	case NullToken:
		return "Null"
	default:
		return ""
	}
}

const (
	StringNode NodeType = iota + 1
	NumberNode
	ObjectNode
	ArrayNode
	BooleanNode
	NullNode
)

const (
	StringToken TokenType = iota + 1
	NumberToken
	SyntaxToken
	BooleanToken
	NullToken
)

// UnMarshal string -> json node
func UnMarshal(source string) (*Node, error) {
	lexer := NewLexer(source)
	tokens, err := lexer.Lex()
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return ast, err
}

// Marshal TODO 支持 jit 预测后直接生成
func Marshal(v *Node, whitespace string) string {
	switch v.typ {
	case StringNode:
		return "\"" + *v.string + "\""
	case BooleanNode:
		if *v.boolean {
			return "true"
		} else {
			return "false"
		}
	case NumberNode:
		return fmt.Sprintf("%f", *v.number)
	case NullNode:
		return "null"

	case ArrayNode:
		s := "[\n"
		a := v.array
		for i := 0; i < len(a); i++ {
			value := a[i]
			s += whitespace + "  " + Marshal(value, whitespace+"  ")
			if i < len(a)-1 {
				s += ","
			}
			s += "\n"
		}

		return s + whitespace + "]"

	case ObjectNode:
		s := "{\n"
		values := v.object
		i := 0
		for key, value := range values {
			s += whitespace + "  " + "\"" + key +
				"\": " + Marshal(value, whitespace+"  ")

			if i < len(values)-1 {
				s += ","
			}

			s += "\n"
			i++
		}

		return s + whitespace + "}"
	}
	return ""
}
