package main

import "fmt"

type (
	// Item item of stream
	Item interface {
		comparable
	}
	// FilterFunc 筛选函数
	FilterFunc[T Item] func(T) bool
	// MapFunc 对象转换函数
	MapFunc[T Item, R Item] func(item T) R
)

// Map
func Map[T Item, R Item](a []T, m MapFunc[T, R]) []R {
	var n []R
	for _, e := range a {
		v := m(e)
		n = append(n, v)
	}
	return n
}

// Filter
func Filter[T Item](a []T, f FilterFunc[T]) []T {
	var n []T
	for _, e := range a {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func main() {
	vi := []int{1, 2, 3, 4, 5, 6}

	vi = Filter(vi, func(v int) bool {
		return v < 4
	})

	bi := Map(vi, func(v int) float32 {
		return float32(v) + 100.01
	})
	fmt.Println(bi)
}
