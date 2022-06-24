package main

import (
	"fmt"

	"github.com/pedrogao/plib/pkg/regex"
)

func main() {
	reg := regex.Compile("a|b")
	ok := reg.Match("a")
	fmt.Println(ok)
	ok = reg.Match("b")
	fmt.Println(ok)

	reg = regex.Compile("(abc)|d")
	ok = reg.Match("d")
	fmt.Println(ok)
	ok = reg.Match("abc")
	fmt.Println(ok)
}
