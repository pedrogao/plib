package jj

import (
	"fmt"
	"strconv"
)

// TokenType token(词)类型
type TokenType int

const (
	LeftBrace   TokenType = iota + 1 // {
	RightBrace                       // }
	LeftSquare                       // [
	RightSquare                      // ]
	Comma                            // ,
	Dot                              // .
	Minus                            // -
	Plus                             // +
	Semicolon                        // ;
	Slash                            // /
	Star                             // *
	Colon                            // :

	Identifier
	String
	Number

	False // false
	Null  // null
	True  // true

	Eof
)

// eof 结束字符
const eof = ' '

// Token 词
type Token struct {
	Type  TokenType
	Raw   string
	Value interface{}
	Line  int
}

// 关键字哈希表
var keywords = map[string]TokenType{
	"null":  Null,  // null
	"true":  True,  // true
	"false": False, // false
}

// Scanner 原始文本扫描，词法解析
type Scanner struct {
	source  string
	runes   []rune // utf8
	tokens  []*Token
	start   int
	current int
	line    int
}

// NewScanner 新建 scanner
func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		runes:   []rune(source),
		tokens:  make([]*Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() ([]*Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}
	// eof token
	s.tokens = append(s.tokens, &Token{
		Type:  Eof,
		Raw:   "",
		Value: nil,
		Line:  s.line,
	})
	return s.tokens, nil
}

func (s *Scanner) scanToken() error {
	c := s.advance()
	switch c {
	case '{':
		s.addRawToken(LeftBrace)
	case '}':
		s.addRawToken(RightBrace)
	case '[':
		s.addRawToken(LeftSquare)
	case ']':
		s.addRawToken(RightSquare)
	case ',':
		s.addRawToken(Comma)
	case '.':
		s.addRawToken(Dot)
	case '-':
		s.addRawToken(Minus)
	case '+':
		s.addRawToken(Plus)
	case ';':
		s.addRawToken(Semicolon)
	case '*':
		s.addRawToken(Star)
	case ':':
		s.addRawToken(Colon)
	case '/': // comment
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addRawToken(Slash)
		}
	case ' ', '\r', '\t': // ignore whitespace
		break
	case '\n':
		s.line += 1
	case '"':
		err := s.string()
		if err != nil {
			return err
		}
	default:
		if s.isDigital(c) {
			err := s.number()
			if err != nil {
				return err
			}
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			return fmt.Errorf("unexpected character: %v", string(c))
		}
	}
	return nil
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		return fmt.Errorf("unterminated string at %d", s.line)
	}

	s.advance() // for "
	value := string(s.runes[s.start+1 : s.current-1])
	s.addToken(String, value)
	return nil
}

func (s *Scanner) number() error {
	for s.isDigital(s.peek()) {
		s.advance()
	}
	// 小数点
	if s.peek() == '.' && s.isDigital(s.peekNext()) {
		s.advance() // for .
		for s.isDigital(s.peek()) {
			s.advance()
		}
	}
	value := string(s.runes[s.start:s.current])
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("parse %s to float err: %s", value, err)
	}
	s.addToken(Number, f)
	return nil
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	value := string(s.runes[s.start:s.current])
	var (
		typ TokenType
		ok  bool
	)
	typ, ok = keywords[value]
	if !ok {
		typ = Identifier
		s.addRawToken(typ)
	} else {
		switch typ {
		case Null:
			s.addToken(Null, nil)
		case True:
			s.addToken(True, true)
		case False:
			s.addToken(False, false)
		}
	}
}

func (s *Scanner) addToken(typ TokenType, value interface{}) {
	sub := string(s.runes[s.start:s.current])
	s.tokens = append(s.tokens, &Token{
		Type:  typ,
		Raw:   sub,
		Value: value,
		Line:  s.line,
	})
}

func (s *Scanner) addRawToken(typ TokenType) {
	s.addToken(typ, nil)
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.runes[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return eof
	}
	return s.runes[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.runes) {
		return eof
	}
	return s.runes[s.current+1]
}

func (s *Scanner) advance() rune {
	c := s.runes[s.current]
	s.current += 1
	return c
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigital(c)
}

func (s *Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *Scanner) isDigital(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.runes)
}
