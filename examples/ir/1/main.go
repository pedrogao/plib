package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

/*
 $ go run ./ > add.ll
 $ clang -o add ./add.ll
 $ ./add

 refer:
   - https://github.com/llir/llvm/issues/22
   - https://stackoverflow.com/questions/31092531/llvm-ir-printing-a-number
*/
func main() {
	// int g = 2;
	//
	// int add(int x, int y) {
	//   return x + y;
	// }
	//
	// int main() {
	//	 int ret = add(1, g);
	//   printf("%d\n", ret);
	//   return 0;
	// }
	m := ir.NewModule()

	zero := constant.NewInt(types.I32, 0)

	globalG := m.NewGlobalDef("g", constant.NewInt(types.I32, 2))

	printf := m.NewFunc(
		"printf",
		types.I32,
		ir.NewParam("", types.NewPointer(types.I8)),
	)
	printf.Sig.Variadic = true

	funcAdd := m.NewFunc("add", types.I32,
		ir.NewParam("x", types.I32),
		ir.NewParam("y", types.I32))

	ab := funcAdd.NewBlock("")
	ab.NewRet(ab.NewAdd(funcAdd.Params[0], funcAdd.Params[1]))

	funcMain := m.NewFunc("main", types.I32)
	mb := funcMain.NewBlock("")

	// add result
	result := mb.NewCall(funcAdd, constant.NewInt(types.I32, 1), mb.NewLoad(types.I32, globalG))

	formatStr := m.NewGlobalDef("formatStr", constant.NewCharArrayFromString("%d\n"))
	format := constant.NewGetElementPtr(types.NewArray(3, types.I8), formatStr, zero, zero)
	// https://llvm.org/docs/GetElementPtr.html
	mb.NewCall(printf, format, result)

	mb.NewRet(zero)

	fmt.Println(m.String())
}
