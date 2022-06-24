package json

import (
	"errors"
	"fmt"
)

func formatParseErr(base string, token *Token) error {
	s := "unexpected token '%s', type '%s', index \n%s"
	prefix := fmt.Sprintf(s, token.value, token.typ, base)
	return formatErr(prefix, token.fullSource, token.location)
}

func formatErr(base string, source []rune, index int) error {
	counter := 0
	line := 1
	column := 0
	lastLine := ""
	whitespace := ""
	for _, c := range source {
		if counter == index {
			break
		}
		if c == '\n' {
			line++
			column = 0
			lastLine = ""
			whitespace = ""
		} else if c == '\t' {
			column++
			lastLine += "  "
			whitespace += "  "
		} else {
			column++
			lastLine += string(c)
			whitespace += " "
		}

		counter++
	}

	// Continue accumulating the lastLine for debugging
	for counter < len(source) {
		c := source[counter]
		if c == '\n' {
			break
		}
		lastLine += string(c)
		counter++
	}
	s := fmt.Sprintf("%s at line %d, column %d\n%s\n%s^", base,
		line, column, lastLine, whitespace)
	return errors.New(s)
}
