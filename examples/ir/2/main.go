package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os/exec"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func main() {
	fset := token.NewFileSet() // 相对于fset的position
	src := `package main

import "fmt"

func add(x int, y int) int {
	return x + y
}

func main() {
	fmt.Println(add(10, 2))
}`

	f, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		fmt.Println(err)
		return
	}

	ast.Print(fset, f)

	m := ir.NewModule()
	// builtin
	printf := m.NewFunc(
		"printf",
		types.I32,
		ir.NewParam("", types.NewPointer(types.I8)),
	)
	printf.Sig.Variadic = true

	funcMap := map[string]*ir.Func{
		"printf": printf,
	}

	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name.Name == "add" {
				translateAdd(m, fn, funcMap)
			} else if fn.Name.Name == "main" {
				translateMain(m, fn, funcMap)
			} else {
				continue
			}
		}
	}

	fmt.Println(m.String())

	err = ioutil.WriteFile("./add.ll", []byte(m.String()), 0666)
	if err != nil {
		panic(err)
	}
	err = exec.Command("clang", "./add.ll").Run()
	if err != nil {
		panic(err)
	}
}

func translateAdd(m *ir.Module, decl *ast.FuncDecl, funcMap map[string]*ir.Func) *ir.Func {
	params := make([]*ir.Param, 0)
	for _, field := range decl.Type.Params.List {
		paramName := field.Names[0].Name
		paramType := field.Type.(*ast.Ident).Name
		if paramType != "int" { // 暂不支持
			continue
		}
		params = append(params, ir.NewParam(paramName, types.I32))
	}
	returnType := decl.Type.Results.List[0].Type.(*ast.Ident).Name
	if returnType != "int" { // 暂不支持
		panic(returnType + " return type is not support")
	}

	funcDefine := m.NewFunc(decl.Name.Name, types.I32, params...)
	ab := funcDefine.NewBlock("")
	ab.NewRet(ab.NewAdd(funcDefine.Params[0], funcDefine.Params[1]))

	funcMap[decl.Name.Name] = funcDefine

	return funcDefine
}

func translateMain(m *ir.Module, decl *ast.FuncDecl, funcMap map[string]*ir.Func) *ir.Func {
	zero := constant.NewInt(types.I32, 0)

	stmt := decl.Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr).Args[0].(*ast.CallExpr)
	args := make([]value.Value, 0)
	for _, item := range stmt.Args {
		val := item.(*ast.BasicLit).Value
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			continue
		}
		args = append(args, constant.NewInt(types.I32, i))
	}

	funcMain := m.NewFunc("main", types.I32)
	mb := funcMain.NewBlock("")
	result := mb.NewCall(funcMap["add"], args...)
	formatStr := m.NewGlobalDef("formatStr", constant.NewCharArrayFromString("%d\n"))
	format := constant.NewGetElementPtr(types.NewArray(3, types.I8), formatStr, zero, zero)
	mb.NewCall(funcMap["printf"], format, result)
	mb.NewRet(zero)
	return funcMain
}
