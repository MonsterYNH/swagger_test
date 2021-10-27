package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type FileParser struct {
	FileInterfaces []FileInterfaceDesc
}

type FileInterfaceDesc struct {
	PackageName string
	Interfaces  []InterfaceDesc
}

type InterfaceDesc struct {
	Name  string
	Funcs []InterfaceFunc
}

type InterfaceFunc struct {
	Name      string
	Comments  []string
	Params    []InterfaceItem
	Results   []InterfaceItem
	Interface *ast.InterfaceType `json:"-"`
}

type InterfaceItem struct {
	Names []string
	Type  string
}

func (fp *FileParser) ParseFile(file string) error {
	astFile, err := parser.ParseFile(token.NewFileSet(), "./api.go", nil, parser.ParseComments)
	if err != nil {
		return err
	}

	fileInterfaceDesc := FileInterfaceDesc{
		PackageName: astFile.Name.Name,
		Interfaces:  []InterfaceDesc{},
	}

	for _, descl := range astFile.Decls {
		desc, ok := descl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range desc.Specs {
			specEntry, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			inter, ok := specEntry.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			interfaceDesc := InterfaceDesc{
				Name:  specEntry.Name.Name,
				Funcs: []InterfaceFunc{},
			}
			for _, method := range inter.Methods.List {
				funcType := method.Type.(*ast.FuncType)

				interfaceFunc := InterfaceFunc{
					Name:      method.Names[0].Name,
					Params:    parseFuncItem(astFile.Name.Name, funcType.Params.List),
					Results:   parseFuncItem(astFile.Name.Name, funcType.Results.List),
					Interface: inter,
				}
				for _, comment := range method.Doc.List {
					interfaceFunc.Comments = append(interfaceFunc.Comments, comment.Text)
				}

				interfaceDesc.Funcs = append(interfaceDesc.Funcs, interfaceFunc)

			}
			fileInterfaceDesc.Interfaces = append(fileInterfaceDesc.Interfaces, interfaceDesc)
		}
	}

	fp.FileInterfaces = append(fp.FileInterfaces, fileInterfaceDesc)

	return err
}

func parseFuncItem(packageName string, fields []*ast.Field) []InterfaceItem {
	items := []InterfaceItem{}
	for _, field := range fields {
		names := []string{}
		for _, nameEntry := range field.Names {
			names = append(names, nameEntry.Name)
		}

		var typeName string
		switch typeEntry := field.Type.(type) {
		case *ast.Ident:
			typeName = typeEntry.Name
		case *ast.StarExpr:
			switch xEntry := typeEntry.X.(type) {
			case *ast.SelectorExpr:
				var x, sel string
				if xEntry.X != nil {
					x = xEntry.X.(*ast.Ident).Name
				}
				if xEntry.Sel != nil {
					sel = xEntry.Sel.Name
				}
				if len(x) != 0 {
					typeName = x + "." + sel
				} else {
					typeName = sel
				}
			case *ast.Ident:
				if xEntry.Obj != nil {
					typeName = packageName + "." + xEntry.Name
				} else {
					typeName = xEntry.Name
				}

			}

		}

		items = append(items, InterfaceItem{
			Names: names,
			Type:  typeName,
		})
	}

	return items
}

func main() {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "./api.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Print(fset, astFile)

	fileParser := FileParser{
		FileInterfaces: []FileInterfaceDesc{},
	}

	if err := fileParser.ParseFile("./api.go"); err != nil {
		panic(err)
	}

	bytes, _ := json.Marshal(fileParser)
	fmt.Println(string(bytes), "============")

}
