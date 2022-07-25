package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func minify(orig []byte) ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	output := make(map[string]any)
	if err := yaml.Unmarshal(orig, output); err != nil {
		return nil, err
	}
	transformMap(output, &b)
	b.WriteByte('}')
	return b.Bytes(), nil
}

func transformMap(in map[string]any, b *bytes.Buffer) {
	index := len(in)
	for k, v := range in {
		b.WriteString(k + ":")
		switch i := v.(type) {
		case string:
			if strings.TrimSpace(i) == "" {
				b.WriteString("''")
			} else {
				b.WriteString(i)
			}
		case map[string]any:
			b.WriteString("{")
			transformMap(i, b)
			b.WriteString("}")
		case []any:
			b.WriteRune('[')
			transformArr(i, b)
			b.WriteRune(']')
		}
		if index > 1 {
			b.WriteString(",")
			index--
		}
	}
}

func transformArr(arr []any, b *bytes.Buffer) {
	index := len(arr)
	for _, v := range arr {
		switch i := v.(type) {
		case string:
			if strings.TrimSpace(i) == "" {
				b.WriteString("''")
			} else {
				b.WriteString(i)
			}
		case map[string]any:
			b.WriteString("{")
			transformMap(i, b)
			b.WriteString("}")
		case []any:
			b.WriteRune('[')
			transformArr(i, b)
			b.WriteRune(']')
		}
		if index > 1 {
			b.WriteString(",")
			index--
		}
	}
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("please specify the input file path [FILE]")
	}
	data, err := ioutil.ReadFile(args[1])
	if err != nil {
		panic(err)
	}
	line, err := minify(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(line))
}
