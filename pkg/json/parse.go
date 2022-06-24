package json

import (
	"fmt"
	"strconv"
)

/*
 json string Parse
*/

type (
	// Parser 语法解析
	Parser struct {
		tokens []*Token
		index  int
	}
)

func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens: tokens,
		index:  0,
	}
}

func (p *Parser) Parse() (*Node, error) {
	token := p.tokens[p.index]
	switch token.typ {
	case NumberToken:
		f, err := strconv.ParseFloat(token.value, 64)
		if err != nil {
			return nil, fmt.Errorf("parse '%s' to number err: %s", token.value, err)
		}
		p.index += 1
		return &Node{
			number: &f,
			typ:    NumberNode,
		}, nil
	case NullToken:
		p.index += 1
		return &Node{
			typ: NullNode,
		}, nil
	case BooleanToken:
		p.index += 1
		v := token.value == "true"
		return &Node{
			boolean: &v,
			typ:     BooleanNode,
		}, nil
	case StringToken:
		p.index += 1
		return &Node{
			string: &token.value,
			typ:    StringNode,
		}, nil
	case SyntaxToken:
		p.index += 1
		if token.value == "[" {
			array, err := p.parseArray()
			if err != nil {
				return nil, err
			}
			return &Node{
				array: array,
				typ:   ArrayNode,
			}, nil
		}
		if token.value == "{" {
			obj, err := p.parseObject()
			if err != nil {
				return nil, err
			}
			return &Node{
				object: obj,
				typ:    ObjectNode,
			}, nil
		}
		return nil, fmt.Errorf("unexpected token: %s, type: %s", token.value, token.typ)
	}
	return nil, fmt.Errorf("unexpected token: %s, type: %s", token.value, token.typ)
}

func (p *Parser) parseArray() ([]*Node, error) {
	// TODO 支持嵌套
	children := make([]*Node, 0)
	for p.index < len(p.tokens) {
		t := p.tokens[p.index]
		if t.typ == SyntaxToken {
			if t.value == "]" {
				// 数组结束
				return children, nil
			}
			if t.value == "," {
				// 数组继续
				p.index++
				t = p.tokens[p.index]
			} else if len(children) > 0 {
				return nil, fmt.Errorf("expected comma after element in array: %s", t.value)
			}
		}
		child, err := p.Parse()
		if err != nil {
			return nil, err
		}
		children = append(children, child)
	}
	return nil, fmt.Errorf("unexpected EOF while parsing array: %s", p.tokens[p.index].value)
}

func (p *Parser) parseObject() (map[string]*Node, error) {
	// TODO 支持嵌套
	values := make(map[string]*Node)
	for p.index < len(p.tokens) {
		t := p.tokens[p.index]
		if t.typ == SyntaxToken {
			if t.value == "}" {
				return values, nil
			}

			if t.value == "," {
				p.index++
				t = p.tokens[p.index]
			} else if len(values) > 0 {
				return nil,
					fmt.Errorf("expected comma after element in object: %s", t.value)
			} else {
				return nil,
					fmt.Errorf("expected key-value pair or closing brace in object: %s", t.value)
			}
		}
		key, err := p.Parse() // key
		if err != nil {
			return nil, err
		}

		if key.typ != StringNode {
			return nil, fmt.Errorf("expected string key in object: %s", t.value)
		}
		t = p.tokens[p.index]

		if !(t.typ == SyntaxToken && t.value == ":") { // 如果不是 : 符号，则报错
			return nil, fmt.Errorf("expected colon after key in object: %s", t.value)
		}
		p.index++
		t = p.tokens[p.index]

		value, err := p.Parse() // value
		if err != nil {
			return nil, err
		}

		values[*key.string] = value
	}

	return values, nil
}

func (p *Parser) curIndex() int {
	return p.index
}
