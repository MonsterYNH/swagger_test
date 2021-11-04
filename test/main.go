package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
	"sort"
)

func main() {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "../example/gin/main.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	config := types.Config{
		Importer: importer.ForCompiler(fset, "source", nil),
	}

	info := types.Info{
		// 表达式对应的类型
		Types: make(map[ast.Expr]types.TypeAndValue),
		// 被定义的标示符
		Defs: make(map[*ast.Ident]types.Object),
		// 被使用的标示符
		Uses: make(map[*ast.Ident]types.Object),
		// 隐藏节点，匿名import包，type-specific时的case对应的当前类型，声明函数的匿名参数如var func(int)
		Implicits: make(map[ast.Node]types.Object),
		// 选择器,只能针对类型/对象.字段/method的选择，package.API这种不会记录在这里
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		// scope 记录当前库scope下的所有域，*ast.File/*ast.FuncType/... 都属于scope，详情看Scopes说明
		// scope关系: 最外层Universe scope,之后Package scope，其他子scope
		Scopes: make(map[ast.Node]*types.Scope),
		// 记录所有package级的初始化值
		InitOrder: make([]*types.Initializer, 0, 0),
	}

	_, err = config.Check("", fset, []*ast.File{file}, &info)
	if err != nil {
		panic(err)
	}

	fmt.Println("------------ def -----------")
	for _, typed := range info.Types {
		fmt.Println(typed.Type.String())
	}

	// for _, node := range getSortedKeys(info.Types) {
	// 	expr := node.(ast.Expr)
	// 	typeValue := info.Types[expr]
	// 	fmt.Printf("%s - %s %T it's value: %v type: %s\n",
	// 		fset.Position(expr.Pos()),
	// 		fset.Position(expr.End()),
	// 		expr,
	// 		typeValue.Value,
	// 		typeValue.Type.String(),
	// 	)
	// 	if typeValue.Assignable() {
	// 		fmt.Print("assignable ")
	// 	}

	// 	if typeValue.Addressable() {
	// 		fmt.Print("addressable ")
	// 	}
	// 	if typeValue.IsNil() {
	// 		fmt.Print("nil ")
	// 	}
	// 	if typeValue.HasOk() {
	// 		fmt.Print("has ok ")
	// 	}
	// 	if typeValue.IsBuiltin() {
	// 		fmt.Print("builtin ")
	// 	}
	// 	if typeValue.IsType() {
	// 		fmt.Print("is type ")
	// 	}
	// 	if typeValue.IsValue() {
	// 		fmt.Print("is value ")
	// 	}
	// 	if typeValue.IsVoid() {
	// 		fmt.Print("void ")
	// 	}
	// 	fmt.Println()
	// }

	// // 打印defs
	// fmt.Println("------------ def -----------")
	// for _, node := range getSortedKeys(info.Defs) {
	// 	ident := node.(*ast.Ident)
	// 	fmt.Println(parse.ExprString(ident), "---------------------")
	// 	object := info.Defs[ident]
	// 	fmt.Printf("%s - %s %T",
	// 		fset.Position(ident.Pos()),
	// 		fset.Position(ident.End()),
	// 		object,
	// 	)
	// 	if object != nil {
	// 		fmt.Printf(" it's object: %s type: %s",
	// 			object.Name(),
	// 			object.Type().Underlying().String(),
	// 		)

	// 	}
	// 	fmt.Println()
	// }

	// // 打印Uses
	// fmt.Println("------------ uses -----------")
	// for _, node := range getSortedKeys(info.Uses) {
	// 	ident := node.(*ast.Ident)
	// 	object := info.Uses[ident]
	// 	fmt.Printf("%s - %s %T",
	// 		fset.Position(ident.Pos()),
	// 		fset.Position(ident.End()),
	// 		object,
	// 	)
	// 	if object != nil {
	// 		fmt.Printf(" it's object: %s type: %s",
	// 			object,
	// 			object.Type().String(),
	// 		)

	// 	}
	// 	fmt.Println()
	// }

	// // 打印Implicits
	// fmt.Println("------------ implicits -----------")
	// for _, node := range getSortedKeys(info.Implicits) {
	// 	object := info.Implicits[node]
	// 	fmt.Printf("%s - %s %T it's object: %s\n",
	// 		fset.Position(node.Pos()),
	// 		fset.Position(node.End()),
	// 		node,
	// 		object.Name(),
	// 	)
	// }

	// // 打印Selections
	// fmt.Println("------------ selections -----------")
	// for _, node := range getSortedKeys(info.Selections) {
	// 	sel := node.(*ast.SelectorExpr)
	// 	typeSel := info.Selections[sel]
	// 	fmt.Printf("%s - %s it's selection: %s\n",
	// 		fset.Position(sel.Pos()),
	// 		fset.Position(sel.End()),
	// 		typeSel.Obj().Name(),
	// 	)
	// 	fmt.Printf("receive: %s index: %v obj: %s\n", typeSel.Recv(), typeSel.Index(), typeSel.Obj().Name())
	// }

	// // 打印Scopes
	// fmt.Println("------------ scopes -----------")
	// //打印package scope
	// fmt.Printf("package level scope %s\n",
	// 	pkg.Scope().String(),
	// )

	// // 打印宇宙级scope
	// fmt.Printf("universe level scope %s\n",
	// 	pkg.Scope().Parent().String(),
	// )
	// for _, node := range getSortedKeys(info.Scopes) {
	// 	scope := info.Scopes[node]
	// 	fmt.Printf("%s - %s %T it's scope %s\n",
	// 		fset.Position(node.Pos()),
	// 		fset.Position(node.End()),
	// 		node,
	// 		scope.Names(),
	// 	)
	// }
}

// 排序规则order by Pos(), End()
func sortNodes(nodes []ast.Node) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Pos() == nodes[j].Pos() {
			return nodes[i].End() < nodes[j].End()
		}
		return nodes[i].Pos() < nodes[j].Pos()
	})
}

// map中的元素是无序的，对key排序打印更好查看
func getSortedKeys(m interface{}) []ast.Node {
	mValue := reflect.ValueOf(m)
	nodes := make([]ast.Node, mValue.Len())
	keys := mValue.MapKeys()
	for i := range keys {
		nodes[i] = keys[i].Interface().(ast.Node)
	}
	sortNodes(nodes)
	return nodes
}
