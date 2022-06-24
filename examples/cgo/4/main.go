package main

/*
#include <stdio.h>
#include <unistd.h>

static void* Sbrk(int size) {
	void *r = sbrk(size);
	if(r == (void *)-1){
  		return NULL;
 	}
	return r;
}
*/
import "C"

import (
	"fmt"
)

func main() {
	mem := C.Sbrk(C.int(100))
	fmt.Println(mem)
}
