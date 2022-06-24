package main

import (
	"fmt"
)

// Addable add interface
type Addable interface {
	int | int64 | int32 | float64 | float32 | string
}

func add[T Addable](a, b T) T {
	return a + b
}

func main() {
	fmt.Println(add(1, 2))
	fmt.Println(add(1.1, 2.1))
	fmt.Println(add("pedro", "gao"))
}
