package main

import (
	"fmt"

	"github.com/pedrogao/trie"
)

func main() {
	// not thread-safe

	t := trie.NewFuzzyTrie()

	t.Put("/root", "mac")
	t.Put("/root/usr", "pedro")
	t.Put("/root/env", "mike")

	val := t.Get("/root")
	fmt.Println(val)

	err := t.WalkPath("/root/*",
		func(key string, value any) error {
			fmt.Printf("key: %s, val: %s\n", key, value)
			return nil
		})
	if err != nil {
		panic(err)
	}

}
