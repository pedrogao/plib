package main

import (
	"encoding/json"
	"fmt"

	"github.com/goccy/go-jit"

	"github.com/pedrogao/plib/pkg/common"
)

type person struct {
	Age int `json:"age"`
}

func generateJsonStr(age int) string {
	return fmt.Sprintf(`{"age":%d}`, age)
}

func main() {
	// {"name":"pedro","age":25}
	bytes, err := json.Marshal(person{
		Age: 25,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(common.Bytes2String(bytes))
	fmt.Println(generateJsonStr(25))

	ctx := jit.NewContext()
	defer ctx.Close()

	f1, err := ctx.Build(func(ctx *jit.Context) (*jit.Function, error) {
		// create an empty function
		f := ctx.CreateFunction([]*jit.Type{jit.TypeInt}, jit.TypeInt)
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

	f1.Run("pedro", 25)
}
