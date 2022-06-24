package main

import (
	"reflect"

	"github.com/bytedance/sonic"
)

// Student test
type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	var s Student
	err := sonic.Pretouch(reflect.TypeOf(s))
	if err != nil {
		panic(err)
	}

}
