package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if err := ParseDir("."); err != nil {
		panic(err)
	}
}

func ParseDir(dir string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	fileDescs := []FileDesc{}
	for _, pkg := range pkgs {
		for fileName, file := range pkg.Files {
			ast.Print(fset, file)
			fileDesc, err := ParseFile(fileName, file)
			if err != nil {
				return err
			}
			fileDescs = append(fileDescs, *fileDesc)
		}
	}

	bytes, _ := json.Marshal(fileDescs)
	fmt.Println(string(bytes))

	return nil
}

type FileDesc struct {
	FilePath    string
	PackageName string
	Imports     []string
	Funcs       []FuncDesc
}

type FuncDesc struct {
	Name      string
	Params    []FuncItem
	Results   []FuncItem
	CallExprs []FuncCallExpr
	TypeDescs []TypeDesc
}

type FuncItem struct {
	SelectorX   string
	SelectorSel string
	Name        string
}

type FuncCallExpr struct {
	Selector FuncSelectorDesc
	Args     []FuncArgDesc
}

type FuncSelectorDesc struct {
	X   string
	Sel string
}

type FuncArgDesc struct {
	Kind        string
	Value       string
	IsObject    bool
	Op          string
	Name        string
	SelectorX   string
	SelectorSel string
}

type FuncDescFilter func(FuncDesc) bool

func genSwaggerComment(desc FuncDesc) []string {
	if !(len(desc.Params) == 1 && desc.Params[0].SelectorX == "gin" && desc.Params[0].SelectorSel == "Context") {
		return nil
	}

	name := desc.Params[0].Name

	comments := []string{}
	for _, callExpr := range desc.CallExprs {
		if callExpr.Selector.X != name {
			continue
		}

		switch callExpr.Selector.Sel {
		case "BindJSON":

		}
	}

}

func ParseFile(fileName string, file *ast.File, filters ...FuncDescFilter) (*FileDesc, error) {
	dir, _ := os.Getwd()
	filePath, err := filepath.Abs(fmt.Sprintf("%s/%s", dir, fileName))
	if err != nil {
		return nil, err
	}

	imports := []string{}
	for _, importDesc := range file.Imports {
		if importDesc.Path != nil {
			imports = append(imports, importDesc.Path.Value)
		}
	}

	fileDesc := FileDesc{
		FilePath:    filePath,
		PackageName: file.Name.Name,
		Imports:     imports,
	}

	for _, desc := range file.Decls {
		node, ok := desc.(*ast.FuncDecl)
		if !ok {
			continue
		}

		funcDesc := FuncDesc{
			Name:      node.Name.Name,
			Params:    []FuncItem{},
			Results:   []FuncItem{},
			CallExprs: []FuncCallExpr{},
			TypeDescs: parseTypeSpec(desc),
		}
		if node.Type.Params != nil && node.Type.Params.List != nil {
			funcDesc.Params = parseFuncItem(node.Type.Params.List)
		}
		if node.Type.Results != nil && node.Type.Results.List != nil {
			funcDesc.Results = parseFuncItem(node.Type.Results.List)
		}

		ast.Inspect(desc, func(n ast.Node) bool {
			switch nodeEntry := n.(type) {
			case *ast.CallExpr:
				funcEntry, ok := nodeEntry.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				funcCallExpr := FuncCallExpr{
					Selector: FuncSelectorDesc{},
					Args:     []FuncArgDesc{},
				}
				if funcEntry.X != nil {
					xEntry, ok := funcEntry.X.(*ast.Ident)
					if ok {
						funcCallExpr.Selector.X = xEntry.Name
					}
				}
				funcCallExpr.Selector.Sel = funcEntry.Sel.Name

				for _, argEntry := range nodeEntry.Args {
					argDesc := FuncArgDesc{}
					switch arg := argEntry.(type) {
					case *ast.UnaryExpr:
						if arg.X != nil {
							argDesc.Name = arg.X.(*ast.Ident).Name
							argDesc.IsObject = arg.X.(*ast.Ident).Obj != nil
						}
						argDesc.Op = arg.Op.String()
					case *ast.BasicLit:
						argDesc.Kind = arg.Kind.String()
						argDesc.Value = arg.Value
					case *ast.Ident:
						argDesc.Name = arg.Name
					case *ast.SelectorExpr:
						if arg.X != nil {
							argDesc.SelectorX = arg.X.(*ast.Ident).Name
						}
						argDesc.SelectorSel = arg.Sel.Name
					default:
						// log.Println("[ERROR] argEntry %v is unknown", argEntry)
					}
					funcCallExpr.Args = append(funcCallExpr.Args, argDesc)
				}

				funcDesc.CallExprs = append(funcDesc.CallExprs, funcCallExpr)
			}

			return true
		})

		if len(filters) == 0 {
			fileDesc.Funcs = append(fileDesc.Funcs, funcDesc)
		} else {
			for _, filter := range filters {
				if filter(funcDesc) {
					fileDesc.Funcs = append(fileDesc.Funcs, funcDesc)
				}
			}
		}

	}

	return &fileDesc, nil
}

func parseFuncItem(fields []*ast.Field) []FuncItem {
	items := []FuncItem{}

	for _, field := range fields {
		for _, nameEntry := range field.Names {
			funcItem := FuncItem{
				Name: nameEntry.Name,
			}

			switch typeEntry := field.Type.(type) {
			case *ast.Ident:
				funcItem.SelectorSel = typeEntry.Name
			case *ast.StarExpr:
				switch entry := typeEntry.X.(type) {
				case *ast.SelectorExpr:
					funcItem.SelectorSel = entry.Sel.Name
					if entry.X != nil {
						funcItem.SelectorX = entry.X.(*ast.Ident).Name
					}
				}
			}

			items = append(items, funcItem)
		}
	}

	return items
}

type TypeDesc struct {
	Name string
	Type string
}

func parseTypeSpec(node ast.Node) []TypeDesc {
	typeDescs := []TypeDesc{}
	ast.Inspect(node, func(n ast.Node) bool {
		typeDesc := TypeDesc{}
		switch typeEntry := n.(type) {
		case *ast.ValueSpec:
			for _, nameEntry := range typeEntry.Names {
				typeDesc.Name = nameEntry.Name

				switch entry := typeEntry.Type.(type) {
				case *ast.Ident:
					typeDesc.Type = entry.Name
				case *ast.SelectorExpr:
					if entry.X != nil {
						typeDesc.Type = entry.X.(*ast.Ident).Name + "." + entry.Sel.Name
					} else {
						typeDesc.Type = entry.Sel.Name
					}
				default:
					log.Println("parse valueSpec type unknown")
				}
				typeDescs = append(typeDescs, typeDesc)
			}
		case *ast.AssignStmt:
			for _, entry := range typeEntry.Lhs {
				switch nameEntry := entry.(type) {
				case *ast.Ident:
					typeDesc.Name = nameEntry.Name
				default:
					log.Println("parse AssignStmt type unknown")
				}

				for _, entry := range typeEntry.Rhs {
					switch valueEntry := entry.(type) {
					case *ast.CompositeLit:
						switch valueTypeEntry := valueEntry.Type.(type) {
						case *ast.SelectorExpr:
							if valueTypeEntry.X != nil {
								typeDesc.Type = valueTypeEntry.X.(*ast.Ident).Name + "." + valueTypeEntry.Sel.Name
							} else {
								typeDesc.Type = valueTypeEntry.Sel.Name
							}
						default:
							log.Println("parse AssignStmt value unknown")
						}
					default:
						log.Println("parse AssignStmt rhs type unknown")
					}
				}
				typeDescs = append(typeDescs, typeDesc)
			}

		default:
			log.Println("parse AssignStmt type unknows")
		}

		return true
	})

	return typeDescs
}
