package main

import (
	"fmt"

	treemap "github.com/liyue201/gostl/ds/map"
)

func main() {
	m := treemap.New()

	m.Insert("a", "pedro")
	m.Insert("d", "mike")
	m.Insert("e", "jerry")

	v := m.LowerBound("d").Key()
	fmt.Println(v)

	v = m.LowerBound("b").Key()
	fmt.Println(v)
}
