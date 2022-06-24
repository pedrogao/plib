package main

import (
	"fmt"

	"github.com/looplab/fsm"
)

func main() {
	fsm := fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},  // closed => open
			{Name: "close", Src: []string{"open"}, Dst: "closed"}, // open => closed
		},
		fsm.Callbacks{
			"before_open": func(event *fsm.Event) {
				fmt.Printf("before callback: %s\n", event.Event)
			},
			"after_open": func(event *fsm.Event) {
				fmt.Printf("after callback: %s\n", event.Event)
			},
		},
	)

	fmt.Println(fsm.Current())

	err := fsm.Event("open")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fsm.Current())

	err = fsm.Event("close")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fsm.Current())
}
