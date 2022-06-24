package main

/**

参考资料：https://c-bata.medium.com/optimizing-go-by-avx2-using-auto-vectorization-in-llvm-118f7b366969
go 如何使用 SIMD 等指令加速计算

*/

import (
	"fmt"
	"unsafe"
)

func SumFloat64Avx2(a []float64) float64 {
	var (
		p1  = unsafe.Pointer(&a[0])
		p2  = unsafe.Pointer(uintptr(len(a)))
		res float64
	)
	__sum_float64(p1, p2, unsafe.Pointer(&res))
	return res
}

func main() {
	ret := SumFloat64Avx2([]float64{1, 2, 7, 5, 1.2})
	fmt.Println(ret)
}
