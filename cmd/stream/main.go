package main

import (
	"fmt"

	"github.com/pedrogao/plib/pkg/stream"
)

func main() {
	ok := stream.Just(1, 2, 3, 4, 5).
		Concat(stream.Just(6, 9, 10)).
		Map(func(item int) int {
			return item + 100
		}).
		Filter(func(item int) bool {
			return item > 105
		}).
		Sort(func(a, b int) bool {
			return a > b
		}).
		AllMatch(func(item int) bool {
			return item > 105
		})
	fmt.Println(ok)
}
