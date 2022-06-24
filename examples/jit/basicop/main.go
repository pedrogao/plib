package main

import (
	"fmt"

	"github.com/goccy/go-jit"
)

// func f(x, y, z int) int {
//   temp1 := x * y
//   temp2 := temp1 + z
//   return temp2
// }

func main() {
	ctx := jit.NewContext()
	defer ctx.Close()

	f1, err := ctx.Build(func(ctx *jit.Context) (*jit.Function, error) {
		// create an empty function
		f := ctx.CreateFunction([]*jit.Type{jit.TypeInt, jit.TypeInt, jit.TypeInt}, jit.TypeInt)
		// add params
		x := f.Param(0)
		y := f.Param(1)
		z := f.Param(2)
		// function body
		temp1 := f.Mul(x, y)
		temp2 := f.Add(temp1, z)
		// return
		f.Return(temp2)
		// compile
		f.Compile()
		return f, nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("result = ", f1.Run(4, 5, 6))
}
