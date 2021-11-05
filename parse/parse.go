package parse

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func ParseDir(path string, filterStr string, routeInfos map[string]gin.RouteInfo) (*Build, error) {
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

			fileAST, err := ParseFileAST(baseName, astTree, fileSet, filterStr, routeInfos)
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

func ParseFileAST(name string, tree *ast.File, fileSet token.FileSet, filterStr string, routeInfos map[string]gin.RouteInfo) (*File, error) {
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
		// 选择器,只能针对类型/对象.字段/method的选择，package.API这种不会记录在这里
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
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

				Name:        decValue.Name.Name,
				Comments:    make([]string, 0),
				PackageName: fmt.Sprintf("%s.%s", tree.Name.Name, decValue.Name.Name),
				Params:      parseFuncItemInfo(decValue.Type.Params, info),
				Results:     parseFuncItemInfo(decValue.Type.Results, info),
				Vars:        make(map[string]FuncItem),
				Exprs:       make(map[string]ExprItem),
			}

			if decValue.Doc != nil && decValue.Doc.List != nil {
				for _, comment := range decValue.Doc.List {
					functionDesc.Comments = append(functionDesc.Comments, comment.Text)
				}
			}

			if decValue.Recv != nil && decValue.Recv.List != nil {
				recv := decValue.Recv.List[0]
				functionDesc.PackageName = fmt.Sprintf("%s.%s.%s", tree.Name.Name, strings.TrimPrefix(ExprString(recv.Type), "*"), decValue.Name.Name)
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

					args := make([]ExprArgItem, 0)

					for _, argEntry := range node.Args {
						argType, exist := info.Types[argEntry]
						if !exist {
							continue
						}

						args = append(args, ExprArgItem{
							Type: argType.Type.String(),
							Name: ExprString(argEntry),
						})
					}

					functionDesc.Exprs[fmt.Sprintf("%s.%s", selectorType.Recv().String(), selectorType.Obj().Name())] = ExprItem{
						Receiver: selectorType.Recv().String(),
						Name:     selectorType.Obj().Name(),
						Args:     args,
					}

				}

				return true
			})

			functionDescs = append(functionDescs, functionDesc)

			s := &GinSwagger{FunctionDesc: &functionDesc}

			decValue.Doc = s.genSwaggerComment(routeInfos)

			printer.Fprint(os.Stdout, &fileSet, decValue)

		default:
			// fmt.Printf("(AST: %T) Skiping\n", decValue)
		}

	}

	file.Functions = functionDescs

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

func GetGinRouteInfos(app *gin.Engine) map[string]gin.RouteInfo {
	routes := make(map[string]gin.RouteInfo)
	for _, info := range app.Routes() {
		routes[info.Handler] = info
	}
	return routes
}
