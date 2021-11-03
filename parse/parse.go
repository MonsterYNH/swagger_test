package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
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

			ast.Print(&fileSet, astTree)

			if fileAST != nil {
				build.AddFile(baseName, fileAST)
			}
		}
	}

	return build, nil
}

func ParseFileAST(name string, tree *ast.File, fileSet token.FileSet, filterStr string) (*File, error) {
	file := NewFile(name, tree)

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

			functionsDesc := FunctionDesc{
				Name:    decValue.Name.Name,
				Params:  parseFuncItem(decValue.Type.Params, fileSet),
				Results: parseFuncItem(decValue.Type.Results, fileSet),
				Vars:    make(map[string]FuncItem),
				Exprs:   make(map[string]ExprItem),
				source:  decValue,
			}

			ast.Inspect(decValue.Body, func(n ast.Node) bool {
				switch nodeEntry := n.(type) {
				case *ast.CallExpr:
					name := ExprString(nodeEntry.Fun)
					args := []ExprArgItem{}
					for _, argEntry := range nodeEntry.Args {
						args = append(args, ExprArgItem{
							Name: ExprString(argEntry),
							Pos:  fileSet.Position(argEntry.Pos()).String(),
						})
					}

					functionsDesc.Exprs[name] = ExprItem{
						CallName: name,
						Args:     args,
					}
				case *ast.ValueSpec:
					if len(nodeEntry.Names) != len(nodeEntry.Values) {
						for _, nameEntry := range nodeEntry.Names {
							functionsDesc.Vars[nameEntry.Name] = FuncItem{
								Name: nameEntry.Name,
								Type: ExprString(nodeEntry.Type),
								Pos:  fileSet.Position(nameEntry.Pos()).String(),
							}
						}
					} else {
						for index := range nodeEntry.Names {
							functionsDesc.Vars[nodeEntry.Names[index].Name] = FuncItem{
								Name: nodeEntry.Names[index].Name,
								Type: ExprString(nodeEntry.Type),
								Pos:  fileSet.Position(nodeEntry.Names[index].Pos()).String(),
							}
						}
					}

				case *ast.AssignStmt:
					if nodeEntry.Tok.String() != ":=" {
						return true
					}
					if len(nodeEntry.Lhs) != len(nodeEntry.Rhs) {
						log.Println("[ERROR] var parse failed")
					}

					for index := 0; index < len(nodeEntry.Lhs); index++ {
						name := ExprString(nodeEntry.Lhs[index])
						functionsDesc.Vars[name] = FuncItem{
							Name: name,
							Type: ExprString(nodeEntry.Rhs[index]),
							Pos:  fileSet.Position(nodeEntry.Pos()).String(),
						}
					}
				}

				return true
			})

			file.Functions = append(file.Functions, functionsDesc)
		default:
			// fmt.Printf("(AST: %T) Skiping\n", decValue)
		}
	}

	return file, nil
}
