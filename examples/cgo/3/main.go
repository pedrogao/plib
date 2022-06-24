package main

/*
#include <stdio.h>

static int Add(int i, int j)
{
  int res = 0;
  __asm__ ("add %1, %2"
    : "=r" (res)
    : "r" (i), "0" (j)
  );
  return res;
}
*/
import "C"
import (
	"fmt"
)

func main() {
	r := C.Add(C.int(2022), C.int(18))
	fmt.Println(r)
}
