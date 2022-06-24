package jj

import (
	"fmt"
)

// Parser 语法解析
type Parser struct {
	tokens  []*Token
	current int
}

// NewParser 新建 parser
func NewParser(tokens []*Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// Parse
func (p *Parser) Parse() (Elem, error) {
	return p.Elem()
}

// Elem
func (p *Parser) Elem() (Elem, error) {
	if p.Match(LeftBrace) {
		return p.Object()
	}
	if p.Match(LeftSquare) {
		return p.Array()
	}
	return p.Literal()
}

// Array
func (p *Parser) Array() (Elem, error) {
	var (
		elem Elem
		err  error
	)

	work := func() {
		elem, err = p.Elem()
		if err != nil {
			return
		}
	}

	val := make([]Elem, 0)
	for ok := true; ok; ok = p.Match(Comma) {
		work()

		val = append(val, elem)
	}

	_, err = p.Consume(RightSquare, "expected ']' after array")
	if err != nil {
		return nil, err
	}
	return NewArrayElem(val), nil
}

// Object
func (p *Parser) Object() (Elem, error) {
	var (
		key  *Token
		elem Elem
		err  error
	)

	work := func() {
		key, err = p.Consume(String, "expect a field name")
		if err != nil {
			return
		}
		_, err = p.Consume(Colon, "expect ':' after a field name")
		if err != nil {
			return
		}
		elem, err = p.Elem()
		if err != nil {
			return
		}
	}

	m := map[string]Elem{}

	for ok := true; ok; ok = p.Match(Comma) {
		work()

		m[key.Value.(string)] = elem
	}

	_, err = p.Consume(RightBrace, "expected '}' after object")
	if err != nil {
		return nil, err
	}

	return NewObjectElem(m), nil
}

// Literal

func (p *Parser) Literal() (Elem, error) {
	if p.Match(String) || p.Match(Number) ||
		p.Match(True) || p.Match(False) ||
		p.Match(Null) {
		previous := p.Previous()
		return NewLiteralElem(previous.Value), nil
	}
	return nil, fmt.Errorf("no match element")
}

// Consume 消费 token
func (p *Parser) Consume(typ TokenType,
	errMsg string) (*Token, error) {
	if p.Check(typ) {
		return p.Advance(), nil
	}
	return nil, fmt.Errorf(errMsg)
}

// Match 匹配 token
func (p *Parser) Match(types ...TokenType) bool {
	for _, typ := range types {
		if p.Check(typ) {
			p.Advance()
			return true
		}
	}
	return false
}

// Check 检查 token
func (p *Parser) Check(typ TokenType) bool {
	if p.IsAtEnd() {
		return false
	}
	return p.Peek().Type == typ
}

// Advance 向前
func (p *Parser) Advance() *Token {
	if !p.IsAtEnd() {
		p.current++
	}
	return p.Previous()
}

// IsAtEnd 结束
func (p *Parser) IsAtEnd() bool {
	return p.Peek().Type == Eof
}

// Previous 之前的 token
func (p *Parser) Previous() *Token {
	return p.tokens[p.current-1]
}

// Peek 当前 token
func (p *Parser) Peek() *Token {
	return p.tokens[p.current]
}

// GetCurrent 当前序号
func (p *Parser) GetCurrent() int {
	return p.current
}

// GetTokens tokens
func (p *Parser) GetTokens() []*Token {
	return p.tokens
}
