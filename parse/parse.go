package parse

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"strings"
)

func ParseDir(path string, filterStr string) (*Build, error) {
	var fileSet token.FileSet

	packages, err := parser.ParseDir(&fileSet, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// create new build for the file set
	build := NewBuild()

	// iterate over all packages in the directory
	for _, pkg := range packages {
		// iterate over all files within the package
		for name, astTree := range pkg.Files {
			baseName := filepath.Base(name)

			fileAST, err := ParseFileAST(baseName, astTree, fileSet, filterStr)
			if err != nil {
				return nil, err
			}

			// ast.Print(&fileSet, astTree)

			if fileAST != nil {
				build.AddFile(baseName, fileAST)
			}
		}
	}

	return build, nil
}

func ParseFileAST(name string, tree *ast.File, fileSet token.FileSet, filterStr string) (*File, error) {
	file := NewFile(name, tree)

	config := types.Config{
		Importer: importer.ForCompiler(&fileSet, "source", nil),
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

	if _, err := config.Check("", &fileSet, []*ast.File{tree}, &info); err != nil {
		return nil, err
	}

	functionDescs := []FunctionDesc{}

	for _, declaration := range tree.Decls {
		switch decValue := declaration.(type) {
		case *ast.FuncDecl:
			expr, err := parser.ParseExpr(filterStr)
			if err != nil {
				return nil, err
			}

			except := strings.ReplaceAll(ExprString(expr), " ", "")
			act := strings.ReplaceAll(ExprString(decValue.Type), " ", "")
			if except != act {
				continue
			}

			log.Printf("[INFO] match %s named %s at line %s", except, decValue.Name.Name, fileSet.Position(decValue.Pos()))

			functionDesc := FunctionDesc{
				source: decValue,

				Name:    decValue.Name.Name,
				Params:  parseFuncItemInfo(decValue.Type.Params, info),
				Results: parseFuncItemInfo(decValue.Type.Results, info),
				Vars:    make(map[string]FuncItem),
				Exprs:   make(map[string]ExprItem),
			}

			ast.Inspect(decValue.Body, func(n ast.Node) bool {
				switch node := n.(type) {
				// 获取函数体变量
				case *ast.Ident:
					if info.Defs[node] != nil {
						functionDesc.Vars[node.Name] = FuncItem{
							Name: node.Name,
							Type: info.Defs[node].Type().String(),
						}
					}
				// 获取函数内函数调用
				case *ast.CallExpr:
					selector, ok := node.Fun.(*ast.SelectorExpr)
					if !ok {
						return true
					}

					selectorType, exist := info.Selections[selector]
					if !exist {
						return true
					}

					if selectorType.Kind() != types.MethodVal {
						return true
					}

					for _, argEntry := range node.Args {
						argType, exist := info.Defs[argEntry]
						if !exist {
							continue
						}

						argType.
					}

					functionDesc.Exprs[fmt.Sprintf("%s.%s", selectorType.Recv().String(), selectorType.Obj().Name())] = ExprItem{
						Receiver: selectorType.Recv().String(),
						Name: selectorType.Obj().Name(),

					}

					functionDesc.Exprs = append(functionDesc.Exprs, ExprItem{

					})

					fmt.Println(, )
				}
				return true
			})

			parseFuncItem(decValue.Type.Params, fileSet)

			// functionsDesc := FunctionDesc{
			// 	Name:    decValue.Name.Name,
			// 	Params:  parseFuncItem(decValue.Type.Params, fileSet),
			// 	Results: parseFuncItem(decValue.Type.Results, fileSet),
			// 	Vars:    make(map[string]FuncItem),
			// 	Exprs:   make(map[string]ExprItem),
			// 	source:  decValue,
			// }

			// ast.Inspect(decValue.Body, func(n ast.Node) bool {
			// 	switch nodeEntry := n.(type) {
			// 	case *ast.CallExpr:
			// 		name := ExprString(nodeEntry.Fun)
			// 		args := []ExprArgItem{}
			// 		for _, argEntry := range nodeEntry.Args {
			// 			args = append(args, ExprArgItem{
			// 				Name: ExprString(argEntry),
			// 				Pos:  fileSet.Position(argEntry.Pos()).String(),
			// 			})
			// 		}

			// 		functionsDesc.Exprs[name] = ExprItem{
			// 			CallName: name,
			// 			Args:     args,
			// 		}
			// 	case *ast.ValueSpec:
			// 		if len(nodeEntry.Names) != len(nodeEntry.Values) {
			// 			for _, nameEntry := range nodeEntry.Names {
			// 				functionsDesc.Vars[nameEntry.Name] = FuncItem{
			// 					Name: nameEntry.Name,
			// 					Type: ExprString(nodeEntry.Type),
			// 					Pos:  fileSet.Position(nameEntry.Pos()).String(),
			// 				}
			// 			}
			// 		} else {
			// 			for index := range nodeEntry.Names {
			// 				functionsDesc.Vars[nodeEntry.Names[index].Name] = FuncItem{
			// 					Name: nodeEntry.Names[index].Name,
			// 					Type: ExprString(nodeEntry.Type),
			// 					Pos:  fileSet.Position(nodeEntry.Names[index].Pos()).String(),
			// 				}
			// 			}
			// 		}

			// 	case *ast.AssignStmt:
			// 		if nodeEntry.Tok.String() != ":=" {
			// 			return true
			// 		}
			// 		if len(nodeEntry.Lhs) != len(nodeEntry.Rhs) {
			// 			log.Println("[ERROR] var parse failed")
			// 		}

			// 		for index := 0; index < len(nodeEntry.Lhs); index++ {
			// 			name := ExprString(nodeEntry.Lhs[index])
			// 			functionsDesc.Vars[name] = FuncItem{
			// 				Name: name,
			// 				Type: ExprString(nodeEntry.Rhs[index]),
			// 				Pos:  fileSet.Position(nodeEntry.Pos()).String(),
			// 			}
			// 		}
			// 	}

			// 	return true
			// })

			// file.Functions = append(file.Functions, functionsDesc)
		default:
			// fmt.Printf("(AST: %T) Skiping\n", decValue)
		}
	}

	return file, nil
}

func parseFuncItemInfo(node *ast.FieldList, info types.Info) []FuncItem {
	items := []FuncItem{}

	if node == nil || node.List == nil {
		return items
	}

	for _, field := range node.List {
		for _, nameEntry := range field.Names {
			value, exist := info.Types[field.Type]
			if !exist {
				continue
			}

			items = append(items, FuncItem{
				Name: nameEntry.Name,
				Type: value.Type.String(),
			})
		}
	}

	return items
}
