package main

import (
	"fmt"
	"sort"
)

func main() {
	s := make([]int, 12)

	fmt.Println(len(s))

	arr := []int{2, 3, 1, 5, 4}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})
	fmt.Println(arr)
}
