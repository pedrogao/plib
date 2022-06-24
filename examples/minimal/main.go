package main

import (
	"fmt"
)

//go:nosplit
//go:noescape
func entry()

func main() {
	fmt.Println("Hello Runtime!")
	entry()
}
