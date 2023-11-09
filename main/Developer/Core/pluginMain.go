package Core

import (
	"GoPoc/main/User"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

// LoadPlugin 有代码用代码,无代码用json. 模式1为json格式,模式2为代码格式
func LoadPlugin(pocName string) int {
	// 如果 SendPoc 为空,则无代码. 使用json进行测试
	astNode, err := parser.ParseFile(token.NewFileSet(), pocName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	isCodeEmpty := false
	funcNum := 0
	currentFuncNum := 0
	// 计算所有的匿名函数数量
	ast.Inspect(astNode, func(n ast.Node) bool {
		_, isFuncLit := n.(*ast.FuncLit)
		if isFuncLit {
			funcNum++
			return false
		}
		return true
	})
	// 如果倒数第二个匿名函数为空则 代码为空
	ast.Inspect(astNode, func(n ast.Node) bool {
		fn, isFuncLit := n.(*ast.FuncLit)
		if isFuncLit {
			currentFuncNum++
			if len(fn.Body.List) == 0 && currentFuncNum == (funcNum-1) {
				isCodeEmpty = true
			}
			return false
		}
		return true
	})

	if !isCodeEmpty && funcNum != 0 {
		return 2
	} else if User.Json != "" {
		return 1
	} else {
		fmt.Println("[-] poc文件存在问题,请检查")
		os.Exit(1)
	}
	return 0
}
