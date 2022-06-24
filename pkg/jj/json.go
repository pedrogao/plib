package jj

import (
	"fmt"
	"reflect"
	"strconv"
)

// Unmarshal source -> struct
func Unmarshal(source string, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("input value invalid")
	}

	ast, err := UnmarshalAst(source)
	if err != nil {
		return err
	}

	switch ast.Type() {
	case ObjectType:
	}

	ref := reflect.TypeOf(v)
	numField := ref.NumField()
	for i := 0; i < numField; i++ {
		field := ref.Field(i)
		j := field.Tag.Get("json")
		if j == "" {
			continue
		}
		// rv.Index(i).Set()
	}
	return nil
}

// UnmarshalAst string -> element
func UnmarshalAst(source string) (Elem, error) {
	scanner := NewScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	elem, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return elem, err
}

// Marshal element -> string
func Marshal(elem Elem, whitespace string) string {
	switch v := elem.(type) {
	case *StringElem:
		return "\"" + v.val + "\""
	case *BoolElem:
		if v.val {
			return "true"
		} else {
			return "false"
		}
	case *NumberElem:
		return strconv.FormatFloat(v.val, 'f', -1, 64)
	case *NullElem:
		return "null"
	case *ArrayElem:
		s := "[\n"
		a := v.val
		for i := 0; i < len(a); i++ {
			value := a[i]
			s += whitespace + "  " + Marshal(value, whitespace+"  ")
			if i < len(a)-1 {
				s += ","
			}
			s += "\n"
		}

		return s + whitespace + "]"

	case *ObjectElem:
		s := "{\n"
		values := v.val
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

// PreTouch 预热 TODO
func PreTouch(v interface{}) error {
	return nil
}
