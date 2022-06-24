package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func main() {
	m := ir.NewModule()

	zero := constant.NewInt(types.I32, 0)
	one := constant.NewInt(types.I32, 1)

	captureStruct := m.NewTypeDef("id_capture", types.NewStruct(
		types.I32,
	))
	captureTyp := types.NewPointer(captureStruct)
	idFn := m.NewFunc("id", types.I32, ir.NewParam("capture", captureTyp))
	idB := idFn.NewBlock("")
	v := idB.NewGetElementPtr(captureStruct, idFn.Params[0], zero, zero)
	idB.NewRet(idB.NewLoad(types.I32, v))
	idClosureTyp := m.NewTypeDef("id_closure", types.NewStruct(
		captureTyp,
		idFn.Type(),
	))

	mainFn := m.NewFunc("main", types.I32)
	b := mainFn.NewBlock("")
	// define a local variable `i`
	i := b.NewAlloca(types.I32)
	b.NewStore(constant.NewInt(types.I32, 10), i)
	// use alloca at here to simplify code, in real case should be `malloc` or `gc_malloc`
	captureInstance := b.NewAlloca(captureStruct)
	ptrToCapture := b.NewGetElementPtr(captureStruct, captureInstance, zero, zero)
	// capture variable
	b.NewStore(b.NewLoad(types.I32, i), ptrToCapture)
	// prepare closure
	idClosure := b.NewAlloca(idClosureTyp)
	ptrToCapturePtr := b.NewGetElementPtr(idClosureTyp, idClosure, zero, zero)
	b.NewStore(captureInstance, ptrToCapturePtr)
	ptrToFuncPtr := b.NewGetElementPtr(idClosureTyp, idClosure, zero, one)
	b.NewStore(idFn, ptrToFuncPtr)
	// assuming we transfer closure into another context
	accessCapture := b.NewGetElementPtr(idClosureTyp, idClosure, zero, zero)
	accessFunc := b.NewGetElementPtr(idClosureTyp, idClosure, zero, one)
	result := b.NewCall(b.NewLoad(idFn.Type(), accessFunc), b.NewLoad(captureTyp, accessCapture))

	printf := m.NewFunc(
		"printf",
		types.I32,
		ir.NewParam("", types.NewPointer(types.I8)),
	)
	printf.Sig.Variadic = true

	printIntegerFormat := m.NewGlobalDef("tmp", constant.NewCharArrayFromString("Pedro\n"))
	pointerToString := b.NewGetElementPtr(types.NewArray(6, types.I8), printIntegerFormat, zero, zero)
	// ignore printf
	b.NewCall(printf, pointerToString, result)

	b.NewRet(constant.NewInt(types.I32, 0))

	fmt.Println(m.String())
}
