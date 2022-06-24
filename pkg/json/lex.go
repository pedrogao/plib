package json

import "fmt"

/*
 json string tokenizer(lexing)
*/

type (
	// Lexer 分词器，词法解析
	Lexer struct {
		source string // source string
		runes  []rune // source string resp by rune
		index  int    // Lex position
	}

	lexerFunc func() (*Token, error)
)

// NewLexer init lexer
// @param source string to Lex
func NewLexer(source string) *Lexer {
	return &Lexer{
		source: source,
		runes:  []rune(source),
		index:  0, // default is 0
	}
}

// Lex 词法解析
// TODO 优化 first-second set算法
func (l *Lexer) Lex() ([]*Token, error) {
	tokens := make([]*Token, 0)
	genericLexers := []lexerFunc{l.lexSyntax, l.lexString, l.lexNumber,
		l.lexNull, l.lexTrue, l.lexFalse}
	for i := 0; i < len(l.runes); i++ {
		l.lexWhitespace()
		newIndex := l.curIndex()
		if newIndex != i {
			i = newIndex - 1 // -1 的原因在于，结束当前的循环后，会自动+1，如果这样就会丢掉关键信息
			continue
		}
		found := false
		for _, lexer := range genericLexers {
			token, err := lexer()
			if err != nil {
				return nil, err
			}
			newIndex := l.curIndex()
			if newIndex != i { // index 变了，表示找到了token
				token.fullSource = l.runes
				tokens = append(tokens, token)
				i = newIndex - 1
				found = true
				break
			}
		}
		if found { // 下一个token
			continue
		}

		return nil, fmt.Errorf("unable to Lex %s at %d", string(l.runes), i)
	}

	return tokens, nil
}

func (l *Lexer) lexWhitespace() {
	index := l.index
	for index < len(l.runes) {
		if l.runes[index] == ' ' ||
			l.runes[index] == '\t' ||
			l.runes[index] == '\n' {
			index++
		} else {
			break
		}
	}
	l.index = index
}

func (l *Lexer) lexSyntax() (*Token, error) {
	index := l.index
	token := &Token{
		value:    "",
		typ:      SyntaxToken,
		location: index,
	}
	c := l.runes[index]
	if c == '[' || c == ']' || c == '{' || c == '}' || c == ':' || c == ',' {
		token.value += string(c)
		index++
	}
	l.index = index
	return token, nil
}

func (l *Lexer) lexString() (*Token, error) {
	index := l.index
	token := &Token{
		value:    "",
		typ:      StringToken,
		location: index,
	}
	c := l.runes[index]
	if c != '"' {
		return token, nil
	}
	index++ // skip '"'
	for {
		if index == len(l.runes) {
			return token, fmt.Errorf("unexpected EOF while lexing string %s %d",
				string(l.runes), index)
		}
		c = l.runes[index]
		if c == '"' {
			break
		}
		token.value += string(c)
		index++
	}
	index++ // skip '"'
	l.index = index
	return token, nil
}

func (l *Lexer) lexNumber() (*Token, error) {
	index := l.index
	token := &Token{
		value:    "",
		typ:      NumberToken,
		location: index,
	}
	for {
		if index == len(l.runes) {
			break
		}

		c := l.runes[index]
		if c < '0' || c > '9' { // TODO 支持浮点数、科学计数法
			break
		}
		token.value += string(c)
		index++
	}
	l.index = index
	return token, nil
}

func (l *Lexer) lexKeyword(keyword string, typ TokenType) (*Token, error) {
	index := l.index
	token := &Token{
		value:    "",
		typ:      typ,
		location: index,
	}
	krunes := []rune(keyword)
	for index < len(l.runes) {
		if krunes[index-l.index] != l.runes[index] {
			break
		}
		index++
	}
	if index-l.index == len(krunes) {
		token.value = keyword
	}
	l.index = index
	return token, nil
}

func (l *Lexer) lexNull() (*Token, error) {
	return l.lexKeyword("null", NullToken)
}

func (l *Lexer) lexTrue() (*Token, error) {
	return l.lexKeyword("true", BooleanToken)
}

func (l *Lexer) lexFalse() (*Token, error) {
	return l.lexKeyword("false", BooleanToken)
}

func (l *Lexer) curIndex() int {
	return l.index
}
