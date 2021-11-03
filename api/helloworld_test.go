package api

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestHelloworld(t *testing.T) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "helloworld.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	for _, decl := range file.Decls {
		node, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		fmt.Println(strings.ReplaceAll(ExprString(node.Type), " ", ""))
	}
}

func Test1(t *testing.T) {
	data, err := parser.ParseExpr("func( http.ResponseWriter,  *http.Request) (int, string)")
	if err != nil {
		panic(err)
	}
	fmt.Println(ExprString(data))
	fmt.Println(data, err)
}
